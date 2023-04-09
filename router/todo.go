package router

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zigbee-s/go-todo-mongo-api/database"
	"github.com/zigbee-s/go-todo-mongo-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddTodoGroup(app *fiber.App) {
	todoGroup := app.Group("/todo")

	todoGroup.Get("/", getTodos)
	todoGroup.Get("/:id", getTodo)
	todoGroup.Post("/", createTodo)
	todoGroup.Put("/:id", updateTodo)
	todoGroup.Delete("/:id", deleteTodo)
}

func getTodos(c *fiber.Ctx) error {
	coll := database.GetDBCollection("todos")

	// find all todos
	dbTodos := make([]models.DbTodo, 0)
	cursor, err := coll.Find(c.Context(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// iterate over the cursor
	for cursor.Next(c.Context()) {
		dbTodo := models.DbTodo{}
		err := cursor.Decode(&dbTodo)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		dbTodos = append(dbTodos, dbTodo)
	}

	return c.Status(200).JSON(fiber.Map{
		"data": dbTodos,
	})
}

func getTodo(c *fiber.Ctx) error {
	coll := database.GetDBCollection("todos")

	// find the todo
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "id is required",
		})
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid Id",
		})
	}

	dbTodo := models.DbTodo{}

	err = coll.FindOne(c.Context(), bson.M{"_id": objectId}).Decode(&dbTodo)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": dbTodo,
	})
}

type createDTO struct {
	Task      string    `json:"task" bson:"task"`
	Done      bool      `json:"done" bson:"done"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}

func createTodo(c *fiber.Ctx) error {
	// Validate the body
	b := new(createDTO)
	if err := c.BodyParser(b); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid Body",
		})
	}

	if b.Task == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Task is required",
		})
	}

	b.CreatedAt = time.Now()

	// create the todo
	coll := database.GetDBCollection("todos")
	result, err := coll.InsertOne(c.Context(), b)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create todo",
			"message": err.Error(),
		})
	}

	// return the book
	return c.Status(201).JSON(fiber.Map{
		"result": result,
	})
}

type updateDTO struct {
	Task        string    `json:"task" bson:"task"`
	Done        bool      `json:"done" bson:"done"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	CompletedAt time.Time `json:"completedAt,omitempty" bson:"completedAt,omitempty"`
}

func updateTodo(c *fiber.Ctx) error {
	// validate the body
	b := new(updateDTO)
	if err := c.BodyParser(b); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid body",
		})
	}
	// get the id
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "id is required",
		})
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	// update the todo
	coll := database.GetDBCollection("todos")
	result, err := coll.UpdateOne(c.Context(), bson.M{"_id": objectId}, bson.M{"$set": b})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to update todo",
			"message": err.Error(),
		})
	}

	// return the todo
	return c.Status(200).JSON(fiber.Map{
		"result": result,
	})
}

func deleteTodo(c *fiber.Ctx) error {
	// get the id
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "id is required",
		})
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	// delete the book
	coll := database.GetDBCollection("todos")
	result, err := coll.DeleteOne(c.Context(), bson.M{"_id": objectId})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to delete book",
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"result": result,
	})
}
