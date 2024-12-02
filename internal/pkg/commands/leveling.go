package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Paranoia8972/PixelBot/internal/db"
	"github.com/Paranoia8972/PixelBot/internal/pkg/utils"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func LevelCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var user *discordgo.User
	options := i.ApplicationCommandData().Options
	if len(options) > 0 && options[0].UserValue(s) != nil {
		user = options[0].UserValue(s)
	} else {
		user = i.Member.User
	}
	guildID := i.GuildID
	userID := user.ID

	currentXP, currentLevel := utils.GetUserXPLevel(guildID, userID)
	xpNeeded := utils.CalculateXPNeeded(currentLevel)

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("Level for %s", user.Username),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Level",
				Value:  strconv.Itoa(currentLevel),
				Inline: true,
			},
			{
				Name:   "XP",
				Value:  fmt.Sprintf("%d/%d", currentXP, xpNeeded),
				Inline: true,
			},
		},
		Color: 0x00FF00,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

func SetLevelChannelCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmdOptions := i.ApplicationCommandData().Options
	channelID := cmdOptions[0].ChannelValue(nil).ID
	guildID := i.GuildID

	_, err := db.GetCollection(cfg.DBName, "level_up_channels").UpdateOne(
		context.TODO(),
		bson.M{"guild_id": guildID},
		bson.M{"$set": bson.M{"channel_id": channelID}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to set the level-up channel.",
				Flags:   64,
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Level-up messages will be sent to <#%s>.", channelID),
			Flags:   64,
		},
	})
}
