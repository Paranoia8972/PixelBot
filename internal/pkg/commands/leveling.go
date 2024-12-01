package commands

import (
	"context"
	"log"
	"strconv"

	"github.com/Paranoia8972/PixelBot/internal/db"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func LevelCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.Member.User
	if len(i.ApplicationCommandData().Options) > 0 {
		user = i.ApplicationCommandData().Options[0].UserValue(s)
	}

	collection := db.GetCollection(cfg.DBName, "users")
	filter := bson.M{"user_id": user.ID}

	var result struct {
		XP    int64 `bson:"xp"`
		Level int64 `bson:"level"`
	}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		result.XP = 0
		result.Level = 0
	}

	message := user.Username + " is level " + strconv.FormatInt(result.Level, 10) +
		" with " + strconv.FormatInt(result.XP, 10) + " XP."

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}

func SetLevelChannelCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channelID := i.ApplicationCommandData().Options[0].ChannelValue(s).ID
	guildID := i.GuildID

	collection := db.GetCollection(cfg.DBName, "guild_settings")
	filter := bson.M{"guild_id": guildID}
	update := bson.M{"$set": bson.M{"level_channel_id": channelID}}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		log.Printf("Failed to set level-up channel: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to set level-up channel.",
				Flags:   64,
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Level-up messages will be sent to <#" + channelID + ">.",
			Flags:   64,
		},
	})
}
