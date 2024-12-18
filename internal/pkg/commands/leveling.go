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
	case "set_channel_requirement":
		SetChannelRequirementCommand(s, i)
	case "get_channel_requirement":
		GetChannelRequirementsCommand(s, i)
	case "delete_channel_requirement":
		DeleteChannelRequirementCommand(s, i)
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
	cmdOptions := i.ApplicationCommandData().Options[0].Options

	level := int(cmdOptions[0].IntValue())
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

func SetChannelRequirementCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmdOptions := i.ApplicationCommandData().Options[0].Options
	channel := cmdOptions[0].ChannelValue(s)
	requiredLevel := cmdOptions[1].IntValue()
	guildID := i.GuildID

	_, err := db.GetCollection(cfg.DBName, "channel_requirements").UpdateOne(
		context.TODO(),
		bson.M{
			"guild_id":   guildID,
			"channel_id": channel.ID,
		},
		bson.M{
			"$set": bson.M{
				"required_level": requiredLevel,
			},
		},
		options.Update().SetUpsert(true),
	)

	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to set channel requirement.",
				Flags:   64,
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Set level requirement for <#%s> to level %d", channel.ID, requiredLevel),
			Flags:   64,
		},
	})
}

func GetChannelRequirementsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guildID := i.GuildID

	cursor, err := db.GetCollection(cfg.DBName, "channel_requirements").Find(
		context.TODO(),
		bson.M{"guild_id": guildID},
	)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to fetch channel requirements",
				Flags:   64,
			},
		})
		return
	}
	defer cursor.Close(context.TODO())

	var requirements []struct {
		ChannelID     string `bson:"channel_id"`
		RequiredLevel int    `bson:"required_level"`
	}
	if err = cursor.All(context.TODO(), &requirements); err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to process requirements",
				Flags:   64,
			},
		})
		return
	}

	if len(requirements) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No channel requirements set",
				Flags:   64,
			},
		})
		return
	}

	content := "Channel Level Requirements:\n"
	for _, req := range requirements {
		content += fmt.Sprintf("<#%s>: Level %d\n", req.ChannelID, req.RequiredLevel)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   64,
		},
	})
}

func DeleteChannelRequirementCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channel := i.ApplicationCommandData().Options[0].Options[0].ChannelValue(s)
	guildID := i.GuildID

	result, err := db.GetCollection(cfg.DBName, "channel_requirements").DeleteOne(
		context.TODO(),
		bson.M{
			"guild_id":   guildID,
			"channel_id": channel.ID,
		},
	)

	if err != nil || result.DeletedCount == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No requirement found for this channel",
				Flags:   64,
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Removed level requirement for <#%s>", channel.ID),
			Flags:   64,
		},
	})
}
