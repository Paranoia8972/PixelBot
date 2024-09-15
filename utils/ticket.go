package utils

import (
	"context"

	"github.com/paranoia8972/PixelBot/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TicketSetup struct {
	GuildID             string `bson:"guild_id"`
	ChannelID           string `bson:"channel_id"`
	CategoryID          string `bson:"category_id"`
	TranscriptChannelID string `bson:"transcript_channel_id"`
}

type Tickets struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	GuildID    string             `bson:"guild_id"`
	UserID     string             `bson:"user_id"`
	ChannelID  string             `bson:"channel_id"`
	Transcript string             `bson:"transcript"`
}

type UserTicket struct {
	GuildID   string `bson:"guild_id"`
	UserID    string `bson:"user_id"`
	TicketNum int    `bson:"ticket_num"`
}

type ModeratorRoles struct {
	GuildID string   `bson:"guild_id"`
	RoleIDs []string `bson:"role_ids"`
}

func SetTicketSetup(guildID, channelID, categoryID, transcriptChannelID string) error {
	collection := db.GetCollection("PixelBot", "ticket_setup")
	filter := bson.M{"guild_id": guildID}
	update := bson.M{"$set": bson.M{"channel_id": channelID, "category_id": categoryID, "transcript_channel_id": transcriptChannelID}}
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	return err
}

func GetTicketSetup(guildID string) (TicketSetup, error) {
	collection := db.GetCollection("PixelBot", "ticket_setup")
	filter := bson.M{"guild_id": guildID}
	var ticketSetup TicketSetup
	err := collection.FindOne(context.Background(), filter).Decode(&ticketSetup)
	return ticketSetup, err
}

func GetNextTicketNumber(guildID, userID string) (int, error) {
	collection := db.GetCollection("PixelBot", "user_tickets")
	filter := bson.M{"guild_id": guildID, "user_id": userID}
	var userTicket UserTicket
	err := collection.FindOne(context.Background(), filter).Decode(&userTicket)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 1, nil
		}
		return 0, err
	}
	return userTicket.TicketNum + 1, nil
}

func IncrementTicketNumber(guildID, userID string, ticketNum int) error {
	collection := db.GetCollection("PixelBot", "user_tickets")
	filter := bson.M{"guild_id": guildID, "user_id": userID}
	update := bson.M{"$set": bson.M{"ticket_num": ticketNum}}
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	return err
}

func AddModeratorRoles(guildID string, roleIDs []string) error {
	collection := db.GetCollection("PixelBot", "moderator_roles")
	filter := bson.M{"guild_id": guildID}
	update := bson.M{"$addToSet": bson.M{"role_ids": bson.M{"$each": roleIDs}}}
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	return err
}

func GetModeratorRoles(guildID string) ([]string, error) {
	collection := db.GetCollection("PixelBot", "moderator_roles")
	filter := bson.M{"guild_id": guildID}
	var modRoles ModeratorRoles
	err := collection.FindOne(context.Background(), filter).Decode(&modRoles)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []string{}, nil
		}
		return nil, err
	}
	return modRoles.RoleIDs, nil
}

func RemoveModeratorRoles(guildID string, roleID string) error {
	collection := db.GetCollection("PixelBot", "moderator_roles")
	filter := bson.M{"guild_id": guildID}
	update := bson.M{"$pull": bson.M{"role": roleID}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}

func StoreTranscript(guildID, userID, channelID string, transcriptJSON []byte) (primitive.ObjectID, error) {
	collection := db.GetCollection("PixelBot", "tickets")
	document := Tickets{
		GuildID:    guildID,
		UserID:     userID,
		ChannelID:  channelID,
		Transcript: string(transcriptJSON),
	}
	result, err := collection.InsertOne(context.Background(), document)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}
