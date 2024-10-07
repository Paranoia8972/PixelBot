package events

import (
	"image"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Paranoia8972/PixelBot/internal/pkg/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

func Ready(s *discordgo.Session, r *discordgo.Ready) {
	color.Blue("Logged in as %s", r.User.Username)
	status := discordgo.UpdateStatusData{
		IdleSince: nil,
		Activities: []*discordgo.Activity{
			{
				Name: "on OnThePixel.net",
				Type: discordgo.ActivityTypeGame,
			},
		},
		Status: "online",
		AFK:    false,
	}

	s.UpdateStatusComplex(status)
}

func AutoRole(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	autoRoles, err := utils.GetAutoRoles(m.GuildID)
	if err != nil {
		log.Printf("Failed to get auto roles: %v", err)
		return
	}

	log.Printf("Auto roles for guild %s: %v", m.GuildID, autoRoles)

	for _, roleID := range autoRoles {
		log.Printf("Assigning role %s to user %s", roleID, m.User.ID)
		err := s.GuildMemberRoleAdd(m.GuildID, m.User.ID, roleID)
		if err != nil {
			log.Printf("Failed to add role %s to user %s: %v", roleID, m.User.ID, err)
		}
	}
}

func MemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	// Auto role logic
	autoRoles, err := utils.GetAutoRoles(m.GuildID)
	if err != nil {
		log.Printf("Failed to get auto roles: %v", err)
		return
	}
	log.Printf("Auto roles for guild %s: %v", m.GuildID, autoRoles)
	for _, roleID := range autoRoles {
		log.Printf("Assigning role %s to user %s", roleID, m.User.ID)
		err := s.GuildMemberRoleAdd(m.GuildID, m.User.ID, roleID)
		if err != nil {
			log.Printf("Failed to add role %s to user %s: %v", roleID, m.User.ID, err)
		}
	}

	// Call the function to send the welcome message
	sendWelcomeMessage(s, m)
}

func sendWelcomeMessage(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	// Welcome message logic
	welcomeChannel, err := utils.GetWelcomeChannel(m.GuildID)
	if err != nil {
		log.Printf("Failed to get welcome channel: %v", err)
		return
	}
	welcomeMessage := welcomeChannel.WelcomeMessage
	if welcomeMessage == "" {
		welcomeMessage = "Welcome to the server, {user}!"
	}
	welcomeMessage = strings.ReplaceAll(welcomeMessage, "{user}", m.User.Mention())

	// Welcome image logic
	bgFile, err := os.Open("assets/bg.png")
	if err != nil {
		log.Printf("Failed to open background image: %v", err)
		return
	}
	defer bgFile.Close()
	bgImg, err := png.Decode(bgFile)
	if err != nil {
		log.Printf("Failed to decode background image: %v", err)
		return
	}

	// Get user's profile picture
	pfpResp, err := http.Get(m.User.AvatarURL("512"))
	if err != nil {
		log.Printf("Failed to download profile picture: %v", err)
		return
	}
	defer pfpResp.Body.Close()
	pfpImg, err := png.Decode(pfpResp.Body)
	if err != nil {
		log.Printf("Failed to decode profile picture: %v", err)
		return
	}

	// Resize profile picture
	pfpImg = resize.Resize(512, 512, pfpImg, resize.Lanczos3)

	// Create new image
	outputImg := image.NewRGBA(bgImg.Bounds())
	draw.Draw(outputImg, bgImg.Bounds(), bgImg, image.Point{}, draw.Over)
	dc := gg.NewContextForRGBA(outputImg)

	// Position profile picture in the center
	pfpX := float64((bgImg.Bounds().Dx() - pfpImg.Bounds().Dx()) / 2)
	pfpY := float64((bgImg.Bounds().Dy() - pfpImg.Bounds().Dy()) / 2)
	pfpY -= 50.0 // Move up by 50 pixels

	// Circular mask for profile picture
	dc.DrawCircle(pfpX+256, pfpY+256, 256)
	dc.Clip()
	dc.DrawImage(pfpImg, int(pfpX), int(pfpY))
	dc.ResetClip()

	// White border for profile picture
	dc.SetLineWidth(10)
	dc.SetRGB(1, 1, 1)
	dc.DrawCircle(pfpX+256, pfpY+256, 256)
	dc.Stroke()

	// Username
	dc.SetRGB(1, 1, 1)
	err = dc.LoadFontFace("assets/Geist-Bold.ttf", 100)
	if err != nil {
		log.Printf("Failed to load font: %v", err)
		return
	}
	spacing := 50.0 // Spacing between profile picture and username
	dc.DrawStringAnchored(m.User.Username, float64(bgImg.Bounds().Dx()/2), float64(pfpY+float64(pfpImg.Bounds().Dy())+30+spacing), 0.5, 0.5)
	dc.Fill()

	// Save output image
	outputFile, err := os.Create("welcome.png")
	if err != nil {
		log.Printf("Failed to create output file: %v", err)
		return
	}
	defer outputFile.Close()
	png.Encode(outputFile, outputImg)

	// Send welcome message
	file, err := os.Open("welcome.png")
	if err != nil {
		log.Printf("Failed to open output image file: %v", err)
		return
	}
	defer file.Close()
	_, err = s.ChannelMessageSendComplex(welcomeChannel.ChannelID, &discordgo.MessageSend{
		Content: welcomeMessage,
		Files: []*discordgo.File{
			{
				Name:   "welcome.png",
				Reader: file,
			},
		},
	})
	if err != nil {
		log.Printf("Failed to send welcome message: %v", err)
	}
}
