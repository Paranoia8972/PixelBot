package utils

import (
	"context"
	"fmt"
	"log"

	"github.com/Paranoia8972/PixelBot/internal/db"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

func Advent1(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "🕯️ First Advent: 1. December",
		Description: "Use the code `advent1_500` to redeem 500 coins once OnThePixel.net opens!\nEach code can only be redeemed once.",
		Color:       0x248045,
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  64,
		},
	})
}

func Advent2(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := i.Member.User.ID
	guildID := i.GuildID

	var adventClick AdventClick
	collection := db.GetCollection(cfg.DBName, "advent_clicks")
	filter := bson.M{"user_id": userID}
	err := collection.FindOne(context.Background(), filter).Decode(&adventClick)
	if err == nil {
		for _, day := range adventClick.LevelUpDays {
			if day == "2" {
				embed := &discordgo.MessageEmbed{
					Title:       "🎄 Advent Calendar: 2. December",
					Description: "You have already received the level-up for today.",
					Color:       0x248045,
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{embed},
						Flags:  64,
					},
				})
				return
			}
		}
	}

	_, currentLevel := GetUserXPLevel(guildID, userID)
	newLevel := currentLevel + 1
	newXP := 0
	SetUserXPLevel(guildID, userID, newXP, newLevel)

	err = StoreAdventClick(userID, i.Member.User.Username, "", "2")
	if err != nil {
		log.Println(err)
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🎄 Advent Calendar: 2. December",
		Description: fmt.Sprintf("You've received a Levelup! 🎉 You are now level %d", newLevel),
		Color:       0x248045,
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  64,
		},
	})
}

func Advent3(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "☃️ Advent Calendar: 3. December",
		Description: "Use the code `snowman3_2024` to redeem for the Pet `Snowman` once OnThePixel.net opens!\n Each code can only be redeemed once.",
		Color:       0x248045,
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  64,
		},
	})
}

func Advent4(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "❄️ Advent Calendar: 4. December",
		Description: "Use the code `advent4_200` to redeem 200 coins once OnThePixel.net opens!\n Each code can only be redeemed once.",
		Color:       0x248045,
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  64,
		},
	})
}

func Advent5(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := i.Member.User.ID
	guildID := i.GuildID

	var adventClick AdventClick
	collection := db.GetCollection(cfg.DBName, "advent_clicks")
	filter := bson.M{"user_id": userID}
	err := collection.FindOne(context.Background(), filter).Decode(&adventClick)
	if err == nil {
		for _, day := range adventClick.LevelUpDays {
			if day == "5" {
				embed := &discordgo.MessageEmbed{
					Title:       "🎄 Advent Calendar: 5. December",
					Description: "You have already received the level-up for today.",
					Color:       0x248045,
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{embed},
						Flags:  64,
					},
				})
				return
			}
		}
	}

	_, currentLevel := GetUserXPLevel(guildID, userID)
	newLevel := currentLevel + 1
	newXP := 0
	SetUserXPLevel(guildID, userID, newXP, newLevel)

	err = StoreAdventClick(userID, i.Member.User.Username, "", "2")
	if err != nil {
		log.Println(err)
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🎄 Advent Calendar: 5. December",
		Description: fmt.Sprintf("You've received a Levelup! 🎉 You are now level %d", newLevel),
		Color:       0x248045,
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  64,
		},
	})
}

func Advent6(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "🎄 Advent Calendar: 6. December",
		Description: "Use the code `nick6_2024` to redeem a for a santa hat once OnThePixel.net opens!",
		Color:       0x248045,
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  64,
		},
	})
}

func Advent7(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "☃️ Advent Calendar: 7. December",
		Description: "Use the code `advent7_200` to redeem 200 coins once OnThePixel.net opens!\nEach code can only be redeemed once.",
		Color:       0x248045,
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  64,
		},
	})
}

func Advent8(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := i.Member.User.ID
	guildID := i.GuildID

	var adventClick AdventClick
	collection := db.GetCollection(cfg.DBName, "advent_clicks")
	filter := bson.M{"user_id": userID}
	err := collection.FindOne(context.Background(), filter).Decode(&adventClick)
	if err == nil {
		for _, day := range adventClick.LevelUpDays {
			if day == "8" {
				embed := &discordgo.MessageEmbed{
					Title:       "🎄 Advent Calendar: 8. December",
					Description: "You have already received the level-up for today.",
					Color:       0x248045,
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{embed},
						Flags:  64,
					},
				})
				return
			}
		}
	}

	_, currentLevel := GetUserXPLevel(guildID, userID)
	newLevel := currentLevel + 1
	newXP := 0
	SetUserXPLevel(guildID, userID, newXP, newLevel)

	err = StoreAdventClick(userID, i.Member.User.Username, "", "2")
	if err != nil {
		log.Println(err)
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🎄 Advent Calendar: 8. December",
		Description: fmt.Sprintf("You've received a Levelup! 🎉 You are now level %d", newLevel),
		Color:       0x248045,
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  64,
		},
	})
}

// TODO
func Advent9(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "❄️ Advent Calendar: 9. December",
		Description: "Use the code `Advent500` to redeem 500 coins once OnThePixel.net opens!",
		Color:       0x248045,
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  64,
		},
	})
}

// TODO
func Advent10(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "🎁 Advent Calendar: 10. December",
		Description: "Use the code `Advent500` to redeem 500 coins once OnThePixel.net opens!",
		Color:       0x248045,
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  64,
		},
	})
}

// TODO
func Advent11(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "🎄 Advent Calendar: 11. December",
		Description: "Use the code `Advent500` to redeem 500 coins once OnThePixel.net opens!",
		Color:       0x248045,
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  64,
		},
	})
}

// TODO
func Advent12(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "☃️ Advent Calendar: 12. December",
		Description: "Use the code `Advent500` to redeem 500 coins once OnThePixel.net opens!",
		Color:       0x248045,
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  64,
		},
	})
}

func Advent13(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "❄️ Advent Calendar: 13. December",
		Description: "You were unlucky and got a piece of coal! 🪨",
		Color:       0x248045,
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  64,
		},
	})
}

func Advent14(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "🎁 Advent Calendar: 14. December",
		Description: "Use the code `snow14_2024` to redeem for a ❄️ snow particle trail once OnThePixel.net opens!\nEach code can only be redeemed once.",
		Color:       0x248045,
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  64,
		},
	})
}

func Advent15(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "🕯️🕯️🕯️ Third Advent: 15. December",
		Description: "Use the code `advent15_500` to redeem 500 coins once OnThePixel.net opens!\nEach code can only be redeemed once.",
		Color:       0x248045,
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  64,
		},
	})
}

func Advent16(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := i.Member.User.ID
	guildID := i.GuildID

	var adventClick AdventClick
	collection := db.GetCollection(cfg.DBName, "advent_clicks")
	filter := bson.M{"user_id": userID}
	err := collection.FindOne(context.Background(), filter).Decode(&adventClick)
	if err == nil {
		for _, day := range adventClick.LevelUpDays {
			if day == "16" {
				embed := &discordgo.MessageEmbed{
					Title:       "🎄 Advent Calendar: 16. December",
					Description: "You have already received the level-up for today.",
					Color:       0x248045,
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{embed},
						Flags:  64,
					},
				})
				return
			}
		}
	}

	_, currentLevel := GetUserXPLevel(guildID, userID)
	newLevel := currentLevel + 1
	newXP := 0
	SetUserXPLevel(guildID, userID, newXP, newLevel)

	err = StoreAdventClick(userID, i.Member.User.Username, "", "2")
	if err != nil {
		log.Println(err)
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🎄 Advent Calendar: 16. December",
		Description: fmt.Sprintf("You've received a Levelup! 🎉 You are now level %d", newLevel),
		Color:       0x248045,
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  64,
		},
	})
}

func Advent17(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "☃️ Advent Calendar: 17. December",
		Description: "Use the code `advent17_200` to redeem 200 coins once OnThePixel.net opens!\nEach code can only be redeemed once.",
		Color:       0x248045,
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  64,
		},
	})
}

// TODO
func Advent18(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "❄️ Advent Calendar: 18. December",
		Description: "Use the code `Advent500` to redeem 500 coins once OnThePixel.net opens!",
		Color:       0x248045,
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  64,
		},
	})
}

// TODO
func Advent19(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "🎁 Advent Calendar: 19. December",
		Description: "Use the code `Advent500` to redeem 500 coins once OnThePixel.net opens!",
		Color:       0x248045,
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  64,
		},
	})
}

func Advent20(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "☃️ Advent Calendar: 20. December",
		Description: "Just a few more days until Christmas!\nYou've received the role `Holiday Hero` on discord 🎅",
		Color:       0x248045,
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  64,
		},
	})
	err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, "1312868677800034376")
	if err != nil {
		log.Println(err)
		return
	}
}

// TODO
func Advent21(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "☃️ Advent Calendar: 21. December",
		Description: "Use the code `Advent500` to redeem 500 coins once OnThePixel.net opens!",
		Color:       0x248045,
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  64,
		},
	})
}

func Advent22(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "🕯️🕯️🕯️🕯️ Fourth Advent: 22. December",
		Description: "Use the code `advent22_500` to redeem 500 coins once OnThePixel.net opens!\nEach code can only be redeemed once.",
		Color:       0x248045,
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  64,
		},
	})
}

func Advent23(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "❄️ Advent Calendar: 23. December",
		Description: "Just one more day until Christmas!\nYou've received the role `Christmas Elf` on discord 🎄",
		Color:       0x248045,
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  64,
		},
	})
}

func Advent24(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "🎅🎁🎄 Christmas Day: 24. December",
		Description: "Use the code `xmas-hero` to redeem for a custom role on the Minecraft Server once OnThePixel.net opens!\nYou'll receive special perks and a custom color for your name.",
		Color:       0x248045,
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  64,
		},
	})
}
