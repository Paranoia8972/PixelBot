package events

import (
	"context"
	"log"
	"math"
	"strconv"

	"github.com/Paranoia8972/PixelBot/internal/db"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	collection := db.GetCollection(cfg.DBName, "users")
	filter := bson.M{"user_id": m.Author.ID}
	update := bson.M{"$inc": bson.M{"xp": 10}}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var result struct {
		UserID string `bson:"user_id"`
		XP     int64  `bson:"xp"`
		Level  int64  `bson:"level"`
	}

	err := collection.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&result)
	if err != nil {
		log.Printf("Failed to update XP: %v", err)
		return
	}

	newLevel := int64(math.Sqrt(float64(result.XP / 100)))
	if newLevel > result.Level {
		_, err := collection.UpdateOne(context.TODO(), filter, bson.M{"$set": bson.M{"level": newLevel}})
		if err != nil {
			log.Printf("Failed to update level: %v", err)
			return
		}

		// Get the level-up channel from the database
		settingsCollection := db.GetCollection(cfg.DBName, "guild_settings")
		var settings struct {
			LevelChannelID string `bson:"level_channel_id"`
		}
		err = settingsCollection.FindOne(context.TODO(), bson.M{"guild_id": m.GuildID}).Decode(&settings)
		if err != nil || settings.LevelChannelID == "" {
			settings.LevelChannelID = m.ChannelID
		}

		s.ChannelMessageSend(settings.LevelChannelID, m.Author.Mention()+" has leveled up to level "+strconv.FormatInt(newLevel, 10)+"!")
	}
}
