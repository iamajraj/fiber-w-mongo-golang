package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type User struct {
	ID   string `json:"id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name,omitempty"`
	Age  int    `json:"age" bson:"age,omitempty"`
}

func main() {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		fmt.Println("failed to connect to db")
		os.Exit(1)
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	// var user User

	// user.Name = "Damn"
	// user.Age = 41

	// data, _ := bson.Marshal(user)

	// client.Database("mydb").CreateCollection(context.TODO(), "users")
	userCollection := client.Database("mydb").Collection("users")
	// _, err = userCollection.InsertOne(context.TODO(), data)

	if err != nil {
		fmt.Println("failed to insert to db")
		os.Exit(1)
	}

	app := fiber.New()

	app.Post("/user", func(c *fiber.Ctx) error {
		var user User
		err := json.Unmarshal(c.Body(), &user)

		if err != nil {
			return fiber.NewError(400, "please provide data name and age")
		}

		data, err := bson.Marshal(user)

		if err != nil {
			return fiber.NewError(500, "internal server error: failed to convert to bson")
		}

		_, err = userCollection.InsertOne(context.TODO(), data)

		if err != nil {
			return fiber.NewError(500, "failed to insert into db")
		}

		return c.SendString("The user has been created")
	})

	app.Get("/user", func(c *fiber.Ctx) error {

		cursor, err := userCollection.Find(context.TODO(), bson.D{})

		if err != nil {
			return fiber.NewError(500, "collection find error")
		}

		var results []User

		err = cursor.All(context.TODO(), &results)

		if err != nil {
			return fiber.NewError(500, "failed to fetch all users")
		}

		if err != nil {
			return fiber.NewError(500, "failed to parse users to json")
		}

		return c.JSON(results)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		var v map[string]interface{}
		err := json.Unmarshal(c.Body(), &v)
		if err != nil {
			fmt.Println("can't parse the json")
			return fiber.NewError(500, "can't parse the json")
		}
		fmt.Println(v["token"])

		return c.JSON(v)
	})

	fmt.Println("Server started on port :8000")
	app.Listen("127.0.0.1:8000")
}
