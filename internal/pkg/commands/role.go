package commands

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func RoleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		respondWithMessage(s, i, "No subcommand provided.")
		return
	}

	switch options[0].Name {
	case "add":
		addRole(s, i)
	case "remove":
		removeRole(s, i)
	case "addall":
		addRoleToAll(s, i)
	case "removeall":
		removeRoleFromAll(s, i)
	default:
		respondWithMessage(s, i, "Unknown subcommand.")
	}
}

func addRole(s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.ApplicationCommandData().Options[0].Options[0].UserValue(s)
	roleID := i.ApplicationCommandData().Options[0].Options[1].RoleValue(s, "").ID

	err := s.GuildMemberRoleAdd(i.GuildID, user.ID, roleID)
	if err != nil {
		respondWithMessage(s, i, "Failed to add role.")
		return
	}

	respondWithMessage(s, i, "Role added successfully!")
}

func removeRole(s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.ApplicationCommandData().Options[0].Options[0].UserValue(s)
	roleID := i.ApplicationCommandData().Options[0].Options[1].RoleValue(s, "").ID

	err := s.GuildMemberRoleRemove(i.GuildID, user.ID, roleID)
	if err != nil {
		respondWithMessage(s, i, "Failed to remove role.")
		return
	}

	respondWithMessage(s, i, "Role removed successfully!")
}

func addRoleToAll(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Failed to defer interaction response: %v", err)
		return
	}

	roleID := i.ApplicationCommandData().Options[0].Options[0].RoleValue(s, "").ID
	members, err := s.GuildMembers(i.GuildID, "", 1000)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
			Content: "Failed to retrieve members.",
		})
		return
	}

	var failedMembers []string
	for _, member := range members {
		err := s.GuildMemberRoleAdd(i.GuildID, member.User.ID, roleID)
		if err != nil {
			log.Printf("Failed to add role to member %s: %v", member.User.ID, err)
			failedMembers = append(failedMembers, member.User.ID)
		}
		time.Sleep(100 * time.Millisecond)
	}

	if len(failedMembers) > 0 {
		failedMessage := "Failed to add role to some members."
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &failedMessage,
		})
		if err != nil {
			log.Printf("Failed to edit interaction response: %v", err)
		}
	} else {
		successMessage := "Role added to all members successfully!"
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &successMessage,
		})
		if err != nil {
			log.Printf("Failed to edit interaction response: %v", err)
		}
	}
}

func removeRoleFromAll(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Printf("Failed to defer interaction response: %v", err)
		return
	}

	roleID := i.ApplicationCommandData().Options[0].Options[0].RoleValue(s, "").ID
	members, err := s.GuildMembers(i.GuildID, "", 1000)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
			Content: "Failed to retrieve members.",
		})
		return
	}

	var failedMembers []string
	for _, member := range members {
		err := s.GuildMemberRoleRemove(i.GuildID, member.User.ID, roleID)
		if err != nil {
			log.Printf("Failed to remove role from member %s: %v", member.User.ID, err)
			failedMembers = append(failedMembers, member.User.ID)
		}
		time.Sleep(100 * time.Millisecond)
	}

	if len(failedMembers) > 0 {
		failedMessage := "Failed to remove role from some members."
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &failedMessage,
		})
		if err != nil {
			log.Printf("Failed to edit interaction response: %v", err)
		}
	} else {
		successMessage := "Role removed from all members successfully!"
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &successMessage,
		})
		if err != nil {
			log.Printf("Failed to edit interaction response: %v", err)
		}
	}
}
