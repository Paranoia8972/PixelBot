package utils

import (
	"context"

	"github.com/Paranoia8972/PixelBot/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddAutoRole(guildID string, roleID string) error {
	collection := db.GetCollection(cfg.DBName, "auto_roles")
	filter := bson.M{"guild_id": guildID}
	update := bson.M{"$addToSet": bson.M{"role_ids": roleID}}
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	return err
}

func GetAutoRoles(guildID string) ([]string, error) {
	collection := db.GetCollection(cfg.DBName, "auto_roles")
	filter := bson.M{"guild_id": guildID}
	var autoRoles struct {
		RoleIDs []string `bson:"role_ids"`
	}
	err := collection.FindOne(context.Background(), filter).Decode(&autoRoles)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []string{}, nil
		}
		return nil, err
	}
	return autoRoles.RoleIDs, nil
}

func RemoveAutoRole(guildID string, roleID string) error {
	collection := db.GetCollection(cfg.DBName, "auto_roles")
	filter := bson.M{"guild_id": guildID}
	update := bson.M{"$pull": bson.M{"role_ids": roleID}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}
