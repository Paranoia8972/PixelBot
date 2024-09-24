package commands

import (
	"time"

	"github.com/Paranoia8972/PixelBot/internal/app/config"
	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

var cfg *config.Config

func init() {
	cfg = config.LoadConfig()
}

func RegisterCommands(s *discordgo.Session, cfg *config.Config) {

	Commands := []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Responds with the Bot's latency.",
		},
		{
			Name:                     "clear",
			Description:              "Deletes messages from a channel.",
			DefaultMemberPermissions: &[]int64{discordgo.PermissionManageMessages}[0],
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "count",
					Description: "Number of messages to delete.",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Required:    true,
				},
				{
					Name:        "user",
					Description: "User whose messages to delete.",
					Type:        discordgo.ApplicationCommandOptionUser,
					Required:    false,
				},
			},
		},
		{
			Name:        "radio",
			Description: "Controls the radio.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "start",
					Description: "Starts the radio.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "stop",
					Description: "Stops the radio.",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
		{
			Name:                     "role",
			Description:              "Role management",
			DefaultMemberPermissions: &[]int64{discordgo.PermissionManageRoles}[0],
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "all",
					Description: "Assign a role to all members",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "role",
							Description: "Role to assign",
							Type:        discordgo.ApplicationCommandOptionRole,
							Required:    true,
						},
					},
				},
				{
					Name:        "add",
					Description: "Assign a role to a specific user",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "user",
							Description: "User to assign the role to",
							Type:        discordgo.ApplicationCommandOptionUser,
							Required:    true,
						},
						{
							Name:        "role",
							Description: "Role to assign",
							Type:        discordgo.ApplicationCommandOptionRole,
							Required:    true,
						},
					},
				},
				{
					Name:        "remove",
					Description: "Remove a role from a specific user",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "user",
							Description: "User to remove the role from",
							Type:        discordgo.ApplicationCommandOptionUser,
							Required:    true,
						},
						{
							Name:        "role",
							Description: "Role to remove",
							Type:        discordgo.ApplicationCommandOptionRole,
							Required:    true,
						},
					},
				},
			},
		},
		{
			Name:                     "say",
			Description:              "Repeats a message.",
			DefaultMemberPermissions: &[]int64{discordgo.PermissionManageMessages}[0],
		},
		{
			Name:                     "edit",
			Description:              "Edits a message.",
			DefaultMemberPermissions: &[]int64{discordgo.PermissionManageMessages}[0],
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "message_id",
					Description: "ID of the message to edit.",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		{
			Name:                     "welcome",
			Description:              "Manage welcome channel",
			DefaultMemberPermissions: &[]int64{discordgo.PermissionAdministrator}[0],
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "set",
					Description: "Set the welcome channel",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "channel",
							Description: "Channel to set as welcome channel",
							Type:        discordgo.ApplicationCommandOptionChannel,
							Required:    true,
						},
						{
							Name:        "message",
							Description: "Welcome message",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
				{
					Name:        "get",
					Description: "Get the current welcome channel",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "delete",
					Description: "Delete the welcome channel entry",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
		{
			Name:                     "social",
			Description:              "Manage social updates channel and accounts",
			DefaultMemberPermissions: &[]int64{discordgo.PermissionAdministrator}[0],
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "set",
					Description: "Set the social updates channel and accounts",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "channel",
							Description: "Channel to set as social updates channel",
							Type:        discordgo.ApplicationCommandOptionChannel,
							Required:    true,
						},
						{
							Name:        "youtube",
							Description: "YouTube username",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    false,
						},
						{
							Name:        "twitch",
							Description: "Twitch username",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    false,
						},
						{
							Name:        "twitter",
							Description: "Twitter username",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    false,
						},
					},
				},
				{
					Name:        "get",
					Description: "Get the current social updates channel and accounts",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "delete",
					Description: "Delete the social updates channel entry",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
		{
			Name:                     "ticket",
			Description:              "Manage ticket system",
			DefaultMemberPermissions: &[]int64{discordgo.PermissionAdministrator}[0],
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "setup",
					Description: "Setup the ticket system",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name: "channel",

							Description: "Channel to set as ticket channel",
							Type:        discordgo.ApplicationCommandOptionChannel,
							Required:    true,
						},
						{
							Name:        "category",
							Description: "Category to set as ticket category",
							Type:        discordgo.ApplicationCommandOptionChannel,
							Required:    true,
						},
						{
							Name:        "transcript",
							Description: "Channel to send transcripts to",
							Type:        discordgo.ApplicationCommandOptionChannel,
							Required:    true,
						},
					},
				},
				{
					Name:        "addmodrole",
					Description: "Set the moderator roles",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "role",
							Description: "Role to set as moderator role",
							Type:        discordgo.ApplicationCommandOptionRole,
							Required:    true,
							Options: []*discordgo.ApplicationCommandOption{
								{
									Name:        "role",
									Description: "Role to set as moderator role",
									Type:        discordgo.ApplicationCommandOptionRole,
									Required:    true,
								},
							},
						},
					},
				},
				{
					Name:        "getmodroles",
					Description: "Get the moderator roles",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "removemodrole",
					Description: "Remove the moderator roles",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "role",
							Description: "Role to remove from moderator roles",
							Type:        discordgo.ApplicationCommandOptionRole,
							Required:    true,
							Options: []*discordgo.ApplicationCommandOption{
								{
									Name:        "role",
									Description: "Role to remove from moderator role",
									Type:        discordgo.ApplicationCommandOptionRole,
									Required:    true,
								},
							},
						},
					},
				},
			},
		},
		{
			Name:                     "giveaway",
			Description:              "Manage giveaways",
			DefaultMemberPermissions: &[]int64{discordgo.PermissionAdministrator}[0],
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "create",
					Description: "Create a giveaway",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "duration",
							Description: "Duration of the giveaway",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
						{
							Name:        "winners",
							Description: "Number of winners",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
						{
							Name:        "prize",
							Description: "Prize of the giveaway",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
				{
					Name:        "edit",
					Description: "Edit a giveaway",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "message_id",
							Description: "ID of the message to edit",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
						{
							Name:        "duration",
							Description: "Duration of the giveaway",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    false,
						},
						{
							Name:        "winners",
							Description: "Number of winners",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    false,
						},
					},
				},
				{
					Name:        "end",
					Description: "End a giveaway",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "message_id",
							Description: "ID of the message to end",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
				{
					Name:        "reroll",
					Description: "Reroll a giveaway",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "message_id",
							Description: "ID of the message to reroll",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
			},
		},
		{
			Name:                     "autorole",
			Description:              "Manage auto roles",
			DefaultMemberPermissions: &[]int64{discordgo.PermissionManageRoles}[0],
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "add",
					Description: "Add an auto role",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "role",
							Description: "Role to add",
							Type:        discordgo.ApplicationCommandOptionRole,
							Required:    true,
						},
					},
				},
				{
					Name:        "get",
					Description: "Get auto roles",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "remove",
					Description: "Remove an auto role",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "role",
							Description: "Role to remove",
							Type:        discordgo.ApplicationCommandOptionRole,
							Required:    true,
						},
					},
				},
			},
		},
	}

	commands := make([]*discordgo.ApplicationCommand, len(Commands))
	copy(commands, Commands)

	start := time.Now()
	_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, cfg.GuildID, commands)
	duration := time.Since(start)

	if err != nil {
		if rErr, ok := err.(*discordgo.RESTError); ok && rErr.Response.StatusCode == 429 {
			color.Yellow("Rate limited: %v", rErr.Message)
		} else {
			color.Red("Cannot bulk overwrite commands: %v", err)
		}
	} else {
		for _, cmd := range commands {
			color.Blue("Registered command: %s", cmd.Name)
		}
		color.Blue("Registered %d commands in %v", len(commands), duration)
	}
}
