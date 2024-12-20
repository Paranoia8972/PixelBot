package utils

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"log"
	"net/http"

	"github.com/Paranoia8972/PixelBot/internal/db"
	"github.com/bwmarrin/discordgo"
	"github.com/fogleman/gg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUserXPLevel(guildID, userID string) (int, int) {
	var result struct {
		XP    int `bson:"xp"`
		Level int `bson:"level"`
	}
	err := db.GetCollection(cfg.DBName, "levels").FindOne(context.TODO(), bson.M{
		"guild_id": guildID,
		"user_id":  userID,
	}).Decode(&result)
	if err != nil {
		return 0, 1
	}
	return result.XP, result.Level
}

func SetUserXPLevel(guildID, userID string, xp, level int) {
	_, err := db.GetCollection(cfg.DBName, "levels").UpdateOne(
		context.TODO(),
		bson.M{
			"guild_id": guildID,
			"user_id":  userID,
		},
		bson.M{
			"$set": bson.M{
				"xp":    xp,
				"level": level,
			},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		log.Printf("Failed to update user level: %v", err)
	}
}

func CalculateXPNeeded(level int) int {
	return level * 100
}

func GetLevelUpChannel(guildID string) string {
	var result struct {
		ChannelID string `bson:"channel_id"`
	}
	err := db.GetCollection(cfg.DBName, "level_up_channels").FindOne(context.TODO(), bson.M{
		"guild_id": guildID,
	}).Decode(&result)
	if err != nil {
		return ""
	}
	return result.ChannelID
}

func GiveLevelRewards(s *discordgo.Session, guildID string, level int) {
	var eligibleUsers []struct {
		UserID string `bson:"user_id"`
		Level  int    `bson:"level"`
	}

	cursor, err := db.GetCollection(cfg.DBName, "levels").Find(context.TODO(),
		bson.M{
			"guild_id": guildID,
			"level":    bson.M{"$gte": level},
		})
	if err != nil {
		log.Printf("Failed to find eligible users: %v", err)
		return
	}
	defer cursor.Close(context.TODO())

	if err = cursor.All(context.TODO(), &eligibleUsers); err != nil {
		log.Printf("Failed to decode eligible users: %v", err)
		return
	}

	roleName := fmt.Sprintf("Level %d", level)
	var roleID string

	roles, err := s.GuildRoles(guildID)
	if err != nil {
		log.Printf("Failed to get guild roles: %v", err)
		return
	}

	for _, role := range roles {
		if role.Name == roleName {
			roleID = role.ID
			break
		}
	}

	if roleID == "" {
		color := 0x00FF00
		perms := int64(discordgo.PermissionSendMessages)
		hoist := true
		newRole, err := s.GuildRoleCreate(guildID, &discordgo.RoleParams{
			Name:        roleName,
			Color:       &color,
			Permissions: &perms,
			Hoist:       &hoist,
		})
		if err != nil {
			log.Printf("Failed to create role: %v", err)
			return
		}
		roleID = newRole.ID
	}

	for _, user := range eligibleUsers {
		err = s.GuildMemberRoleAdd(guildID, user.UserID, roleID)
		if err != nil {
			log.Printf("Failed to add role to user %s: %v", user.UserID, err)
			continue
		}
	}

	members, err := s.GuildMembers(guildID, "", 1000)
	if err != nil {
		log.Printf("Failed to get guild members: %v", err)
		return
	}

	for _, member := range members {
		hasRole := false
		for _, memberRole := range member.Roles {
			if memberRole == roleID {
				hasRole = true
				break
			}
		}

		if !hasRole {
			continue
		}

		_, userLevel := GetUserXPLevel(guildID, member.User.ID)
		if userLevel >= level {
			continue
		}
		err = s.GuildMemberRoleRemove(guildID, member.User.ID, roleID)
		if err != nil {
			log.Printf("Failed to remove role from user %s: %v", member.User.ID, err)
		}
	}
}

func GetChannelRequirement(guildID, channelID string) int {
	var result struct {
		RequiredLevel int `bson:"required_level"`
	}
	err := db.GetCollection(cfg.DBName, "channel_requirements").FindOne(
		context.TODO(),
		bson.M{
			"guild_id":   guildID,
			"channel_id": channelID,
		},
	).Decode(&result)

	if err != nil {
		return 0
	}
	return result.RequiredLevel
}

func GetUserRank(guildID, userID string, weekly bool) (int, int) {
	collection := "levels"
	if weekly {
		collection = "weekly_levels"
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"guild_id": guildID,
			},
		},
		{
			"$sort": bson.M{
				"xp": -1,
			},
		},
	}

	cursor, err := db.GetCollection(cfg.DBName, collection).Aggregate(context.TODO(), pipeline)
	if err != nil {
		return 0, 0
	}
	defer cursor.Close(context.TODO())

	rank := 1
	totalUsers := 0
	for cursor.Next(context.TODO()) {
		var result struct {
			UserID string `bson:"user_id"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		if result.UserID == userID {
			return rank, totalUsers + 1
		}
		rank++
		totalUsers++
	}
	return 0, totalUsers + 1
}

func GetWeeklyXP(guildID, userID string) int {
	var result struct {
		XP int `bson:"xp"`
	}
	err := db.GetCollection(cfg.DBName, "weekly_levels").FindOne(context.TODO(), bson.M{
		"guild_id": guildID,
		"user_id":  userID,
	}).Decode(&result)
	if err != nil {
		return 0
	}
	return result.XP
}

func GenerateLevelImage(s *discordgo.Session, user *discordgo.User, guildID string) (*bytes.Buffer, error) {
	const width = 800
	const height = 200

	dc := gg.NewContext(width, height)

	dc.SetRGBA(0, 0, 0, 0)
	dc.Clear()

	// Left block - Profile section
	dc.SetRGB(0.15, 0.15, 0.15)
	dc.DrawRoundedRectangle(20, 20, 560, 160, 10)
	dc.Fill()

	// Right top block - Level
	dc.DrawRoundedRectangle(600, 20, 180, 75, 10)
	dc.Fill()

	// Right bottom block - XP
	dc.DrawRoundedRectangle(600, 105, 180, 75, 10)
	dc.Fill()

	// Get user avatar
	resp, err := http.Get(user.AvatarURL("128"))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	avatar, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	// Avatar circle
	dc.DrawCircle(100, 100, 50)
	dc.Clip()
	dc.DrawImage(avatar, 50, 50)
	dc.ResetClip()

	xp, level := GetUserXPLevel(guildID, user.ID)
	weeklyXP := GetWeeklyXP(guildID, user.ID)
	totalRank, _ := GetUserRank(guildID, user.ID, false)
	weeklyRank, _ := GetUserRank(guildID, user.ID, true)
	xpNeeded := CalculateXPNeeded(level)

	// Text for username
	dc.SetRGB(1, 1, 1)
	dc.LoadFontFace("assets/Jersey20-Regular.ttf", 32)
	dc.DrawString(user.Username, 170, 60)

	// Stats labels
	dc.SetRGB(0.6, 0.6, 0.6)
	dc.LoadFontFace("assets/Jersey20-Regular.ttf", 24)
	dc.DrawString("Server Rank", 170, 90)
	dc.DrawString("Weekly Rank", 300, 90)
	dc.DrawString("Weekly XP", 430, 90)

	// Stats values
	dc.LoadFontFace("assets/Jersey20-Regular.ttf", 30)

	// Server rank
	dc.SetRGB(0.4, 0.8, 1.0)
	dc.DrawStringAnchored(fmt.Sprintf("%d", totalRank), 215, 120, 0.5, 0)

	// Weekly stats
	dc.SetRGB(1, 1, 1)
	dc.DrawStringAnchored(fmt.Sprintf("%d", weeklyRank), 345, 120, 0.5, 0)
	dc.DrawStringAnchored(fmt.Sprintf("%d", weeklyXP), 475, 120, 0.5, 0)

	// Level in right top block
	dc.LoadFontFace("assets/Jersey20-Regular.ttf", 30)
	dc.DrawStringAnchored(fmt.Sprintf("Level %d", level), 690, 60, 0.5, 0.5)

	// XP in right bottom block
	dc.SetRGB(1, 1, 1)
	dc.DrawStringAnchored(fmt.Sprintf("%d/%d XP", xp, xpNeeded), 690, 145, 0.5, 0.5)

	// XP Progress bar
	barWidth := 560
	barHeight := 10
	progress := float64(xp) / float64(xpNeeded)

	// Background bar
	dc.SetRGB(0.2, 0.2, 0.2)
	dc.DrawRoundedRectangle(20, 170, float64(barWidth), float64(barHeight), 5)
	dc.Fill()

	// Progress bar
	dc.SetRGB(0.4, 0.4, 1.0)
	dc.DrawRoundedRectangle(20, 170, float64(barWidth)*progress, float64(barHeight), 5)
	dc.Fill()

	buf := new(bytes.Buffer)
	err = dc.EncodePNG(buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
