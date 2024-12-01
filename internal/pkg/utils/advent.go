package utils

import (
	"context"
	"sort"
	"strconv"
	"time"

	"github.com/Paranoia8972/PixelBot/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AdventClick struct {
	UserID    string    `bson:"user_id"`
	Username  string    `bson:"username"`
	Buttons   []string  `bson:"buttons"`
	Timestamp time.Time `bson:"timestamp"`
}

func StoreAdventClick(userID, username, buttonID string) error {
	collection := db.GetCollection(cfg.DBName, "advent_clicks")
	filter := bson.M{"user_id": userID}

	// Retrieve the current document
	var adventClick AdventClick
	err := collection.FindOne(context.Background(), filter).Decode(&adventClick)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	// Add the new buttonID if it doesn't already exist
	buttonExists := false
	for _, b := range adventClick.Buttons {
		if b == buttonID {
			buttonExists = true
			break
		}
	}
	if !buttonExists {
		adventClick.Buttons = append(adventClick.Buttons, buttonID)
	}

	// Sort the buttons
	sort.Slice(adventClick.Buttons, func(i, j int) bool {
		dayI, _ := strconv.Atoi(adventClick.Buttons[i][len("advent_"):])
		dayJ, _ := strconv.Atoi(adventClick.Buttons[j][len("advent_"):])
		return dayI < dayJ
	})

	update := bson.M{
		"$set": bson.M{
			"username":  username,
			"buttons":   adventClick.Buttons,
			"timestamp": time.Now(),
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err = collection.UpdateOne(context.Background(), filter, update, opts)
	return err
}
