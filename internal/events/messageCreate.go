package events

import (
	"fmt"

	"github.com/Paranoia8972/PixelBot/internal/pkg/utils"
	"github.com/bwmarrin/discordgo"
)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
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
