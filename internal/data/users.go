package data

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	CreatedAt time.Time          `bson:"createdAt" json:"created_at"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Password  password           `bson:"password" json:"-"`
	IsAuthor  bool               `bson:"isAuthor" json:"activated"`
}

type password struct {
	plaintext string           `bson:"-"`
	Hash      primitive.Binary `bson:"hash"`
}

type UserModel struct {
	Collection *mongo.Collection
}

func (p *password) Get() []byte {
	return p.Hash.Data
}

func (p *password) Set(plainTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = plainTextPassword
	p.Hash.Data = hash
	return nil
}

func (b UserModel) Insert(user *User) (*mongo.InsertOneResult, error) {

	userBson, err := bson.Marshal(user)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := b.Collection.InsertOne(ctx, userBson)
	if err != nil {
		return nil, err
	}
	return result, nil

}

func (b UserModel) GetByEmail(email string) (*User, error) {

	filter := bson.D{{"email", email}}
	fmt.Println(email)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := b.Collection.FindOne(ctx, filter)

	var user User
	err := result.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (b UserModel) Update(user *User) (*mongo.UpdateResult, error) {
	filter := bson.D{{"_id", user.ID}}
	update := bson.D{{"$set", bson.D{
		{"name", user.Name},
		{"email", user.Email},
		{"password.hash", user.Password.Hash},
		{"isAuthor", user.IsAuthor}}}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := b.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b UserModel) Delete(user *User) (*mongo.DeleteResult, error) {
	filter := bson.D{{"_id", user.ID}}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := b.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}
