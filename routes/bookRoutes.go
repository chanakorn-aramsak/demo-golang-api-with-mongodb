package routes

import (
	"REST/models"
	"context"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"os"

	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
)

var collection *mongo.Collection

func init() {
	// Initialize MongoDB connection
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGO_URI")
	fmt.Println(mongoURI)
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("bookstore").Collection("books")
	fmt.Println("Connected to MongoDB!")
}

// GetBooks fetch all books
// @Summary Get all books
// @Description Get all book data
// @Produce json
// @Success 200 {array} models.Book
// @Router /api/books [get]
func GetBooks(c *gin.Context) {
	bookList := []models.Book{}
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching books"})
		return
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var book models.Book
		cursor.Decode(&book)
		bookList = append(bookList, book)
	}
	c.JSON(http.StatusOK, bookList)
}

// GetBook fetch a single book by ID
// @Summary Get a book by ID
// @Description Get a single book
// @Produce json
// @Param id path string true "ID of the book to get"
// @Success 200 {object} models.Book
// @Router /api/books/{id} [get]
func GetBook(c *gin.Context) {
	id := c.Param("id")
	objectID, _ := primitive.ObjectIDFromHex(id)
	var book models.Book
	err := collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&book)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

// AddBook to add a new book
// @Summary Add a new book
// @Description Add a new book
// @Accept  json
// @Produce  json
// @Param   book body models.Book true "Book info"
// @Success 201 {string} string	"ok"
// @Router /api/books [post]
func AddBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
		return
	}
	book.ID = primitive.NewObjectID()
	insertResult, err := collection.InsertOne(context.TODO(), book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting book"})
		return
	}
	c.JSON(http.StatusCreated, insertResult.InsertedID)
}

// UpdateBook to update a book
// @Summary Update an existing book
// @Description Update an existing book by ID
// @Accept  json
// @Produce  json
// @Param   id path string true "ID of the book to update"
// @Param   book body models.Book true "Book info"
// @Success 200 {object} models.Book
// @Router /api/books/{id} [put]
func UpdateBook(c *gin.Context) {
	id := c.Param("id")
	objectID, _ := primitive.ObjectIDFromHex(id)
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
		return
	}
	updateResult, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectID},
		bson.D{
			{"$set", bson.D{{"title", book.Title}, {"author", book.Author}}},
		},
	)
	if err != nil || updateResult.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found or not updated"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book updated"})
}

// DeleteBook to delete a book
// @Summary Delete a book
// @Description Delete a book by ID
// @Produce  json
// @Param id path string true "ID of the book to delete"
// @Success 200 {object} models.Book
// @Router /api/books/{id} [delete]
func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	objectID, _ := primitive.ObjectIDFromHex(id)
	deleteResult, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil || deleteResult.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
}
