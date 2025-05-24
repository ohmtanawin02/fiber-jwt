package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// note omitempty คือ field นี้ เป็น zero value ค่าจะไม่ถูกส่งไป
// มีประโยชน์กรณี patch

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Password  string             `json:"password,omitempty" bson:"password,omitempty"`
	IsDeleted bool               `json:"is_deleted,omitempty" bson:"is_deleted,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

type UserResponse struct {
	ID        primitive.ObjectID `json:"_id,omitempty"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	IsDeleted bool               `json:"is_deleted"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt *time.Time         `json:"updated_at"`
	DeletedAt *time.Time         `json:"deleted_at"`
}

type GetUserPaginateRequest struct {
	Search string `json:"search"`
	Email  string `json:"email"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

type GetUserRequest struct {
	ID string `json:"_id"`
}
