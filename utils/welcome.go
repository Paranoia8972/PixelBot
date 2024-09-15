package utils

import (
	"context"

	"github.com/paranoia8972/PixelBot/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WelcomeChannel struct {
	GuildID        string `bson:"guild_id"`
	ChannelID      string `bson:"channel_id"`
	WelcomeMessage string `bson:"welcome_message"`
}

func SetWelcomeChannel(guildID, channelID, welcomeMessage string) error {
	collection := db.GetCollection("PixelBot", "welcome_channels")
	filter := bson.M{"guild_id": guildID}
	update := bson.M{"$set": bson.M{"channel_id": channelID, "welcome_message": welcomeMessage}}
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	return err
}

func GetWelcomeChannel(guildID string) (WelcomeChannel, error) {
	collection := db.GetCollection("PixelBot", "welcome_channels")
	filter := bson.M{"guild_id": guildID}
	var result WelcomeChannel
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	return result, err
}

func DeleteWelcomeChannel(guildID string) error {
	collection := db.GetCollection("PixelBot", "welcome_channels")
	filter := bson.M{"guild_id": guildID}
	_, err := collection.DeleteOne(context.Background(), filter)
	return err
}
