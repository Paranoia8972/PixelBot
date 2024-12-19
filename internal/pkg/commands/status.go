package commands

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Paranoia8972/PixelBot/internal/db"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

type Status struct {
	Message string                 `bson:"message"`
	Type    discordgo.ActivityType `bson:"type"`
}

func StatusCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if len(i.ApplicationCommandData().Options) == 0 {
		RespondWithMessage(s, i, "Please provide a subcommand.")
		return
	}

	switch i.ApplicationCommandData().Options[0].Name {
	case "add":
		addStatus(s, i)
	case "remove":
		removeStatus(s, i)
	case "list":
		listStatuses(s, i)
	}
}

func addStatus(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options[0].Options
	message := options[0].StringValue()
	activityType := discordgo.ActivityType(options[1].IntValue())

	collection := db.GetCollection(cfg.DBName, "statuses")
	_, err := collection.InsertOne(context.TODO(), Status{
		Message: message,
		Type:    activityType,
	})

	if err != nil {
		RespondWithMessage(s, i, "Failed to add status.")
		return
	}

	RespondWithMessage(s, i, "Status added successfully!")
}

func removeStatus(s *discordgo.Session, i *discordgo.InteractionCreate) {
	message := i.ApplicationCommandData().Options[0].Options[0].StringValue()

	collection := db.GetCollection(cfg.DBName, "statuses")
	result, err := collection.DeleteOne(context.TODO(), bson.M{"message": message})

	if err != nil {
		RespondWithMessage(s, i, "Failed to remove status.")
		return
	}

	if result.DeletedCount == 0 {
		RespondWithMessage(s, i, "Status not found.")
		return
	}

	RespondWithMessage(s, i, "Status removed successfully!")
}

func listStatuses(s *discordgo.Session, i *discordgo.InteractionCreate) {
	collection := db.GetCollection(cfg.DBName, "statuses")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		RespondWithMessage(s, i, "Failed to fetch statuses.")
		return
	}
	defer cursor.Close(context.TODO())

	var statuses []Status
	if err = cursor.All(context.TODO(), &statuses); err != nil {
		RespondWithMessage(s, i, "Failed to decode statuses.")
		return
	}

	if len(statuses) == 0 {
		RespondWithMessage(s, i, "No statuses found.")
		return
	}

	message := "Current statuses:\n"
	for _, status := range statuses {
		message += fmt.Sprintf("- %s (%s)\n", status.Message, getActivityTypeName(status.Type))
	}

	RespondWithMessage(s, i, message)
}

func StartStatusRotation(s *discordgo.Session) {
	go func() {
		for {
			time.Sleep(time.Duration(5+rand.Intn(5)) * time.Second)
			rotateStatus(s)
		}
	}()
}

func rotateStatus(s *discordgo.Session) {
	collection := db.GetCollection(cfg.DBName, "statuses")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return
	}
	defer cursor.Close(context.TODO())

	var statuses []Status
	if err = cursor.All(context.TODO(), &statuses); err != nil {
		return
	}

	if len(statuses) == 0 {
		return
	}

	status := statuses[rand.Intn(len(statuses))]
	s.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{
			{
				Name: status.Message,
				Type: status.Type,
			},
		},
		Status: "online",
	})
}

func getActivityTypeName(activityType discordgo.ActivityType) string {
	switch activityType {
	case discordgo.ActivityTypeGame:
		return "Playing"
	case discordgo.ActivityTypeStreaming:
		return "Streaming"
	case discordgo.ActivityTypeListening:
		return "Listening to"
	case discordgo.ActivityTypeWatching:
		return "Watching"
	case discordgo.ActivityTypeCustom:
		return "Custom"
	case discordgo.ActivityTypeCompeting:
		return "Competing in"
	default:
		return "Unknown"
	}
}
