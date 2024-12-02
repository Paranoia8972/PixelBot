package commands

import (
	"log"
	"strconv"
	"time"

	"github.com/Paranoia8972/PixelBot/internal/pkg/utils"
	"github.com/bwmarrin/discordgo"
)

func AdventCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	currentDay := time.Now().Day()
	currentMonth := time.Now().Month()
	currentYear := time.Now().Year()

	if currentYear > 2024 || (currentYear == 2024 && currentMonth == time.December && currentDay > 24) {
		currentDay = 24
	}
	userID := i.Member.User.ID

	adventClick, err := utils.GetAdventClick(userID)
	if err != nil {
		log.Printf("Failed to retrieve advent clicks: %v", err)
	}

	emojis := map[int]string{
		1:  "1313186442557526026",
		2:  "1313186433170800734",
		3:  "1313186420952662027",
		4:  "1313186411494768681",
		5:  "1313186401701068920",
		6:  "1313186389390655508",
		7:  "1313186379211210752",
		8:  "1313186370742915152",
		9:  "1313186360508551198",
		10: "1313186351885058129",
		11: "1313183932228960256",
		12: "1313183956258127933",
		13: "1313183970149797989",
		14: "1313183993109151845",
		15: "1313184006900289558",
		16: "1313184022767206441",
		17: "1313184038885916733",
		18: "1313184058611601418",
		19: "1313184073585528922",
		20: "1313184088185634867",
		21: "1313184102765170830",
		22: "1313184114882379978",
		23: "1313184129273303150",
		24: "1313184144708210769",
	}

	emojiNames := map[int]string{
		1:  "1_",
		2:  "2_",
		3:  "3_",
		4:  "4_",
		5:  "5_",
		6:  "6_",
		7:  "7_",
		8:  "8_",
		9:  "9_",
		10: "10",
		11: "11",
		12: "12",
		13: "13",
		14: "14",
		15: "15",
		16: "16",
		17: "17",
		18: "18",
		19: "19",
		20: "20",
		21: "21",
		22: "22",
		23: "23",
		24: "24",
	}

	var buttons []discordgo.MessageComponent
	for j := 1; j <= 24; j++ {
		style := discordgo.SecondaryButton
		disabled := true
		emojiID := emojis[j]
		emojiName := emojiNames[j]

		if j <= currentDay {
			style = discordgo.SuccessButton
			disabled = false
		}

		// Check if the button has been clicked by the user
		if adventClick != nil && utils.HasButtonBeenClicked(adventClick, "advent_"+strconv.Itoa(j)) {
			emojiID = ""
			emojiName = "🔓"
		}

		if j == currentDay {
			style = discordgo.DangerButton
		}

		buttons = append(buttons, discordgo.Button{
			CustomID: "advent_" + strconv.Itoa(j),
			Style:    style,
			Disabled: disabled,
			Emoji: &discordgo.ComponentEmoji{
				ID:   emojiID,
				Name: emojiName,
			},
		})
	}

	var rows []discordgo.MessageComponent
	for k := 0; k < 24; k += 5 {
		rows = append(rows, discordgo.ActionsRow{
			Components: buttons[k:min(k+5, 24)],
		})
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🎄 Advent Calendar - Day " + strconv.Itoa(currentDay),
		Description: "## Open today's surprise gift! Check back daily for more festive surprises! 🎁",
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Happy holidays from the OnThePixel.net team! 🎅",
		},
		Color: 0x248045,
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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

	levelUpDay := strconv.Itoa(day)

	err = utils.StoreAdventClick(userID, username, buttonID, levelUpDay)
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
	}

	if adventFunc, exists := adventFunctions[day]; exists {
		adventFunc(s, i)
	} else {
		log.Printf("No advent function found for day %d", day)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
