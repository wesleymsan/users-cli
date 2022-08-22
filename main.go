package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/gookit/color.v1"
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
			{
				Name:    "all",
				Aliases: []string{"l"},
				Usage:   "list all users",
				Action: func(c *cli.Context) error {
					users, err := getAll()
					if err != nil {
						if err == mongo.ErrNoDocuments {
							fmt.Print("Nothing to see here.\nRun `add 'user'` to add a user")
							return nil
						}

						return err
					}

					printUsers(users)
					return nil
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

func getAll() ([]*User, error) {
	filter := bson.D{{}}
	return filterUsers(filter)
}

func filterUsers(filter interface{}) ([]*User, error) {
	var users []*User

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return users, err
	}

	for cur.Next(ctx) {
		var t User
		err := cur.Decode(&t)
		if err != nil {
			return users, err
		}

		users = append(users, &t)
	}

	if err := cur.Err(); err != nil {
		return users, err
	}

	cur.Close(ctx)

	if len(users) == 0 {
		return users, mongo.ErrNoDocuments
	}

	return users, nil
}

func printUsers(users []*User) {
	for i, v := range users {
		if v.Completed {
			color.Green.Printf("%d: %s\n", i+1, v.Name)
		} else {
			color.Yellow.Printf("%d: %s\n", i+1, v.Name)
		}
	}
}
