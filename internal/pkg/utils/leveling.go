package utils

import (
	"context"
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

func GiveLevelRewards(s *discordgo.Session, guildID, userID string, level int) {
	switch level {
	case 5:
		roleID := "ROLE_ID_FOR_LEVEL_5"
		s.GuildMemberRoleAdd(guildID, userID, roleID)
	case 10:
		roleID := "ROLE_ID_FOR_LEVEL_10"
		s.GuildMemberRoleAdd(guildID, userID, roleID)
	}
}
