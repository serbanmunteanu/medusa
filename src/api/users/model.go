package users

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"medusa/src/common/users"
	"time"
)

type SignUpInput struct {
	Name      string    `json:"name" bson:"name" binding:"required"`
	Email     string    `json:"email" bson:"email" binding:"required"`
	Password  string    `json:"password" bson:"password" binding:"required,min=8"`
	Role      string    `json:"role" bson:"role"`
	Verified  bool      `json:"verified" bson:"verified"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type SignInInput struct {
	Email    string `json:"email" bson:"email" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}

type UserResponse struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Role      string             `json:"role,omitempty" bson:"role,omitempty"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

func MapToUserResponse(user *users.UserDbModel) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
