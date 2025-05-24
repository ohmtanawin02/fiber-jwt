package handlers

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	DB          *mongo.Client
	DBName      string
	Database    *mongo.Database
	Collections map[string]*mongo.Collection
}

func NewHandler(db *mongo.Client, dbName string) *Handler {
	database := db.Database(dbName)
	handler := &Handler{
		DB:          db,
		DBName:      dbName,
		Database:    database,
		Collections: make(map[string]*mongo.Collection),
	}

	handler.initCollections()

	return handler
}

func (h *Handler) initCollections() {
	h.Collections["users"] = h.Database.Collection("users")
	// h.Collections["users"] = h.Database.Collection("users")
}

func (h *Handler) GetCollection(name string) *mongo.Collection {
	if collection, exists := h.Collections[name]; exists {
		return collection
	}
	collection := h.Database.Collection(name)
	h.Collections[name] = collection
	return collection
}
