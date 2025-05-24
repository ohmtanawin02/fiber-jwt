package services

import "go.mongodb.org/mongo-driver/mongo"

type Service struct {
	DB          *mongo.Client
	DBName      string
	Database    *mongo.Database
	Collections map[string]*mongo.Collection
}

func NewService(db *mongo.Client, dbName string) *Service {
	database := db.Database(dbName)
	service := &Service{
		DB:          db,
		DBName:      dbName,
		Database:    database,
		Collections: make(map[string]*mongo.Collection),
	}

	service.initCollections()

	return service
}

func (h *Service) initCollections() {
	h.Collections["users"] = h.Database.Collection("users")
	// h.Collections["users"] = h.Database.Collection("users")
}

func (h *Service) GetCollection(name string) *mongo.Collection {
	if collection, exists := h.Collections[name]; exists {
		return collection
	}
	collection := h.Database.Collection(name)
	h.Collections[name] = collection
	return collection
}
