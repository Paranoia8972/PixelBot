package events

import (
	"fmt"
	"time"

	"github.com/Paranoia8972/PixelBot/internal/pkg/utils"
	"github.com/bwmarrin/discordgo"
)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if m.GuildID == "" {
		logChannelID, err := utils.GetDMLogChannel(cfg.GuildID)
		if err != nil {
			return
		}

		if logChannelID != "" {
			embed := &discordgo.MessageEmbed{
				Author: &discordgo.MessageEmbedAuthor{
					Name:    m.Author.Username,
					IconURL: m.Author.AvatarURL(""),
				},
				Description: m.Content,
				Color:       0x248045,
				Timestamp:   m.Timestamp.Format(time.RFC3339),
			}

			if len(m.Attachments) > 0 {
				embed.Image = &discordgo.MessageEmbedImage{
					URL: m.Attachments[0].URL,
				}
			}

			s.ChannelMessageSendEmbed(logChannelID, embed)
		}
	}

	xpGain := 10
	userID := m.Author.ID
	guildID := m.GuildID

	currentXP, currentLevel := utils.GetUserXPLevel(guildID, userID)
	newXP := currentXP + xpGain
	xpNeeded := utils.CalculateXPNeeded(currentLevel)

	if newXP >= xpNeeded {
		newLevel := currentLevel + 1
		newXP -= xpNeeded
		utils.SetUserXPLevel(guildID, userID, newXP, newLevel)

		levelUpChannelID := utils.GetLevelUpChannel(guildID)
		if levelUpChannelID != "" {
			s.ChannelMessageSend(levelUpChannelID, fmt.Sprintf("Congratulations %s, you've reached level %d!", m.Author.Mention(), newLevel))
		}

		utils.GiveLevelRewards(s, guildID, newLevel)
	} else {
		utils.SetUserXPLevel(guildID, userID, newXP, currentLevel)
	}
}
