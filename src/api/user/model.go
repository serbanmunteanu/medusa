package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"medusa/src/common/user"
	"time"
)

type UserResponse struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Role      string             `json:"role,omitempty" bson:"role,omitempty"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

func MapToUserResponse(user *user.UserDbModel) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func MapToUsersResponse(users []user.UserDbModel) []UserResponse {
	response := make([]UserResponse, 0)
	for _, user := range users {
		response = append(response, MapToUserResponse(&user))
	}
	return response
}
