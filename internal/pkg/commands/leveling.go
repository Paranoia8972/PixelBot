package commands

import (
	"context"
	"fmt"
	"log"
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

func LevelingCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if len(i.ApplicationCommandData().Options) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No subcommand provided.",
			},
		})
		return
	}
	switch i.ApplicationCommandData().Options[0].Name {
	case "setlevelchannel":
		SetLevelChannelCommand(s, i)
	case "set_reward":
		AddLevelRewardCommand(s, i)
	case "get_reward":
		GetLevelRewardsCommand(s, i)
	case "remove_reward":
		RemoveLevelRewardCommand(s, i.GuildID, i.Member.User.ID)
	default:
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Unknown subcommand.",
			},
		})
	}
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

func AddLevelRewardCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmdOptions := i.ApplicationCommandData().Options
	level := cmdOptions[0].IntValue()
	role := cmdOptions[1].RoleValue(s, i.GuildID)

	guildID := i.GuildID

	_, err := db.GetCollection(cfg.DBName, "level_rewards").UpdateOne(
		context.TODO(),
		bson.M{"guild_id": guildID, "level": level},
		bson.M{"$set": bson.M{"role_id": role.ID}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		log.Printf("Failed to add level reward: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to add the level reward.",
				Flags:   64,
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Role <@&%s> will be assigned at level %d.", role.ID, level),
			Flags:   64,
		},
	})
}

func GetLevelRewardsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guildID := i.GuildID
	var results []struct {
		Level  int    `bson:"level"`
		RoleID string `bson:"role_id"`
	}
	cursor, err := db.GetCollection(cfg.DBName, "level_rewards").Find(context.TODO(), bson.M{"guild_id": guildID})
	if err != nil {
		log.Printf("Failed to get level rewards: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to get level rewards.",
				Flags:   64,
			},
		})
		return
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var result struct {
			Level  int    `bson:"level"`
			RoleID string `bson:"role_id"`
		}
		if err := cursor.Decode(&result); err != nil {
			log.Printf("Failed to decode level reward: %v", err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Failed to get level rewards.",
					Flags:   64,
				},
			})
			return
		}
		results = append(results, result)
	}

	var content string
	for _, result := range results {
		content += fmt.Sprintf("Level %d: <@&%s>\n", result.Level, result.RoleID)
	}

	if content == "" {
		content = "No level rewards found."
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   64,
		},
	})
}

func RemoveLevelRewardCommand(s *discordgo.Session, guildID, userID string) {
	var results []struct {
		Level  int    `bson:"level"`
		RoleID string `bson:"role_id"`
	}
	cursor, err := db.GetCollection(cfg.DBName, "level_rewards").Find(context.TODO(), bson.M{"guild_id": guildID})
	if err != nil {
		log.Printf("Failed to get level rewards: %v", err)
		return
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var result struct {
			Level  int    `bson:"level"`
			RoleID string `bson:"role_id"`
		}
		if err := cursor.Decode(&result); err != nil {
			log.Printf("Failed to decode level reward: %v", err)
			return
		}
		results = append(results, result)
	}

	for _, result := range results {
		if result.RoleID == userID {
			_, err := db.GetCollection(cfg.DBName, "level_rewards").DeleteOne(
				context.TODO(),
				bson.M{"guild_id": guildID, "level": result.Level},
			)
			if err != nil {
				log.Printf("Failed to remove level reward: %v", err)
			}
		}
	}
}
