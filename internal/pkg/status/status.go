package status

import (
	"context"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Status struct {
	Message string                 `bson:"message"`
	Type    discordgo.ActivityType `bson:"type"`
}

type Manager struct {
	session    *discordgo.Session
	collection *mongo.Collection
}

func NewManager(s *discordgo.Session, collection *mongo.Collection) *Manager {
	m := &Manager{
		session:    s,
		collection: collection,
	}
	go m.startRotation()
	return m
}

func (m *Manager) AddStatus(message string, activityType discordgo.ActivityType) error {
	_, err := m.collection.InsertOne(context.Background(), Status{
		Message: message,
		Type:    activityType,
	})
	return err
}

func (m *Manager) RemoveStatus(message string) bool {
	result, err := m.collection.DeleteOne(context.Background(), bson.M{"message": message})
	if err != nil {
		return false
	}
	return result.DeletedCount > 0
}

func (m *Manager) GetStatuses() []Status {
	var statuses []Status
	cursor, err := m.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return []Status{}
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &statuses); err != nil {
		return []Status{}
	}
	return statuses
}

func (m *Manager) startRotation() {
	ticker := time.NewTicker(time.Duration(5+rand.Intn(5)) * time.Second)
	for range ticker.C {
		m.updateRandomStatus()
	}
}

func (m *Manager) updateRandomStatus() {
	statuses := m.GetStatuses()
	if len(statuses) == 0 {
		return
	}

	status := statuses[rand.Intn(len(statuses))]
	m.session.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{{
			Name: status.Message,
			Type: status.Type,
		}},
		Status: "online",
		AFK:    false,
	})
}
