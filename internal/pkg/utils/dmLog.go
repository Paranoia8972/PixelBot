package utils

import (
	"context"

	"github.com/Paranoia8972/PixelBot/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DMLogChannel struct {
	GuildID   string `bson:"guild_id"`
	ChannelID string `bson:"channel_id"`
}

func SetDMLogChannel(guildID, channelID string) error {
	collection := db.GetCollection(cfg.DBName, "dm_log_channels")
	filter := bson.M{"guild_id": guildID}
	update := bson.M{"$set": bson.M{"channel_id": channelID}}
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	return err
}

func GetDMLogChannel(guildID string) (string, error) {
	collection := db.GetCollection(cfg.DBName, "dm_log_channels")
	filter := bson.M{"guild_id": guildID}
	var result DMLogChannel
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return "", err
	}
	return result.ChannelID, nil
}
