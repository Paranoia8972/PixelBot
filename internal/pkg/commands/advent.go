package commands

import (
	"log"
	"strconv"
	"time"

	"github.com/Paranoia8972/PixelBot/internal/pkg/utils"
	"github.com/bwmarrin/discordgo"
)

// var debugDay int

// func init() {
// 	SetDebugDay(24)
// }

func AdventCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	currentDay := time.Now().Day()
	//* debugDay start
	// if debugDay >= 1 && debugDay <= 24 {
	// 	currentDay = debugDay
	// }
	//* debugDay end
	var buttons []discordgo.MessageComponent
	for j := 1; j <= 24; j++ {
		style := discordgo.SecondaryButton
		disabled := true
		label := strconv.Itoa(j)
		if j <= currentDay {
			style = discordgo.SuccessButton
			disabled = false
		}
		if j == currentDay {
			label = "🎄 " + label
		}
		buttons = append(buttons, discordgo.Button{
			Label:    label,
			CustomID: "advent_" + strconv.Itoa(j),
			Style:    style,
			Disabled: disabled,
		})
	}

	var rows []discordgo.MessageComponent
	for k := 0; k < 24; k += 5 {
		rows = append(rows, discordgo.ActionsRow{
			Components: buttons[k:min(k+5, 24)],
		})
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Advent Calendar",
		Description: "# Advent, Advent, a little light is burning! 🎄🎁",
		Color:       0x248045,
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: rows,
			Flags:      64,
		},
	})
	if err != nil {
		log.Printf("Failed to respond to interaction: %v", err)
	}
}

func HandleAdventButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := i.Member.User.ID
	username := i.Member.User.Username
	buttonID := i.MessageComponentData().CustomID
	day, err := strconv.Atoi(buttonID[len("advent_"):])
	if err != nil {
		log.Printf("Invalid button ID: %v", err)
		return
	}
	err = utils.StoreAdventClick(userID, username, buttonID)
	if err != nil {
		log.Printf("Failed to store advent click: %v", err)
	}

	adventFunctions := map[int]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		1:  utils.Advent1,
		2:  utils.Advent2,
		3:  utils.Advent3,
		4:  utils.Advent4,
		5:  utils.Advent5,
		6:  utils.Advent6,
		7:  utils.Advent7,
		8:  utils.Advent8,
		9:  utils.Advent9,
		10: utils.Advent10,
		11: utils.Advent11,
		12: utils.Advent12,
		13: utils.Advent13,
		14: utils.Advent14,
		15: utils.Advent15,
		16: utils.Advent16,
		17: utils.Advent17,
		18: utils.Advent18,
		19: utils.Advent19,
		20: utils.Advent20,
		21: utils.Advent21,
		22: utils.Advent22,
		23: utils.Advent23,
		24: utils.Advent24,
	}

	if adventFunc, ok := adventFunctions[day]; ok {
		adventFunc(s, i)
	}
}

// func SetDebugDay(day int) {
// 	if day >= 1 && day <= 24 {
// 		debugDay = day
// 	} else {
// 		debugDay = 0
// 	}
// }

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
