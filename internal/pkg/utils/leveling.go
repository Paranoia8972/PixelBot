package utils

import (
	"context"
	"fmt"
	"log"

	"github.com/Paranoia8972/PixelBot/internal/db"
	"github.com/bwmarrin/discordgo"
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
