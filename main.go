package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var collection *mongo.Collection
var ctx = context.TODO()

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Text      string             `bson:"text"`
	Completed bool               `bson:"completed"`
}

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("user").Collection("users")
}

func main() {
	app := &cli.App{
		Name:  "users",
		Usage: "A simple CLI program to manage your users",
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add a user to the list",
				Action: func(c *cli.Context) error {
					name := c.Args().First()
					if name == "" {
						return errors.New("Cannot add an empty user")
					}

					user := &User{
						ID:        primitive.NewObjectID(),
						Name:      name,
						CreatedAt: time.Time{},
						UpdatedAt: time.Time{},
						Completed: false,
					}

					return createUser(user)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func createUser(user *User) error {
	_, err := collection.InsertOne(ctx, user)
	return err
}
