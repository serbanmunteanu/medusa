package users

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"medusa/src/common/mongodb"
)

type UserRepository struct {
	client     *mongodb.Client
	collection string
}

func NewUserRepository(client *mongodb.Client, collection string) *UserRepository {
	return &UserRepository{
		client:     client,
		collection: collection,
	}
}

func (us *UserRepository) Insert(user *UserDbModel) error {
	id, err := us.client.Insert(us.collection, user)

	if err != nil {
		return err
	}

	user.ID = id.(primitive.ObjectID)

	return nil
}
