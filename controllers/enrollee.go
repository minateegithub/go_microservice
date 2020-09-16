package controllers

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	guuid "github.com/google/uuid"
	"github.com/minateegithub/go_microservice/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// DATABASE INSTANCE
var collection *mongo.Collection

//EnrolleeCollection - Initialize the collection
func EnrolleeCollection(c *mongo.Database) {
	collectionName := os.Getenv("DB_ENROLLE_COLLECTION_NAME")
	collection = c.Collection(collectionName)
}

//GetAllEnrollees - Gets all enrollee records from db
func GetAllEnrollees(c *gin.Context) {
	enrollees := []models.Enrollee{}
	cursor, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		log.Printf("Error while getting all todos, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	defer cursor.Close(context.TODO())

	// Iterate through the returned cursor.
	for cursor.Next(context.TODO()) {
		var enrollee models.Enrollee
		err := cursor.Decode(&enrollee)
		if err != nil {
			log.Fatal(err)
		}
		enrollees = append(enrollees, enrollee)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All Enrollees",
		"data":    enrollees,
	})
	return
}

//CreateEnrollee - creats a new enrollee
func CreateEnrollee(c *gin.Context) {
	var enrollee models.Enrollee

	err2 := c.BindJSON(&enrollee)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	validationErrors := validateNewEnroleeData(enrollee)
	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	name := enrollee.Name
	isActive := enrollee.IsActive
	birthDate := enrollee.BirthDate
	phoneNumber := enrollee.PhoneNumber
	dependents := enrollee.Dependents
	id := guuid.New().String()

	if dependents == nil {
		dependents = []models.Dependent{}
	} else {

		for i := 0; i < len(dependents); i++ {
			dependentID := guuid.New().String()
			dependents[i].ID = &dependentID
			dependents[i].CreatedAt = time.Now()
			dependents[i].UpdatedAt = time.Now()
		}
	}

	newEnrollee := models.Enrollee{
		ID:          &id,
		Name:        name,
		IsActive:    isActive,
		BirthDate:   birthDate,
		PhoneNumber: phoneNumber,
		Dependents:  dependents,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := collection.InsertOne(context.TODO(), newEnrollee)

	if err != nil {
		log.Printf("Error while inserting new enrolle into db, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Internal error",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Enrollee created Successfully",
	})
	return
}

//GetSingleEnrollee - Gets the enrollee record for the provided id
func GetSingleEnrollee(c *gin.Context) {
	enrolleeID := c.Param("enrolleeId")

	enrollee := models.Enrollee{}
	err := collection.FindOne(context.TODO(), bson.M{"id": enrolleeID}).Decode(&enrollee)
	if err != nil {
		log.Printf("Error while getting the enrollee, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Enrollee not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Single Enrollee",
		"data":    enrollee,
	})
	return
}

//EditEnrollee - Updates enrollee data
func EditEnrollee(c *gin.Context) {
	enrolleeID := c.Param("enrolleeId")
	var enrollee models.Enrollee

	err2 := c.BindJSON(&enrollee)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	validationErrors := validateEditedEnroleeData(enrollee)
	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	currentEnrollee := models.Enrollee{}
	err := collection.FindOne(context.TODO(), bson.M{"id": enrolleeID}).Decode(&currentEnrollee)
	if err != nil {
		log.Printf("Error while getting the enrollee, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Enrollee not found",
		})
		return
	}
	if enrollee.Name != nil {
		currentEnrollee.Name = enrollee.Name
	}
	if enrollee.BirthDate != nil {
		currentEnrollee.BirthDate = enrollee.BirthDate
	}
	if enrollee.PhoneNumber != nil {
		currentEnrollee.PhoneNumber = enrollee.PhoneNumber
	}
	if enrollee.IsActive != nil {
		currentEnrollee.IsActive = enrollee.IsActive
	}

	newData := bson.M{
		"$set": bson.M{
			"name":         currentEnrollee.Name,
			"is_active":    currentEnrollee.IsActive,
			"birth_date":   currentEnrollee.BirthDate,
			"phone_number": currentEnrollee.PhoneNumber,
			"updated_at":   time.Now(),
		},
	}

	_, err1 := collection.UpdateOne(context.TODO(), bson.M{"id": enrolleeID}, newData)
	if err1 != nil {
		log.Printf("Error, Reason: %v\n", err1)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Enrollee Edited Successfully",
	})
	return
}

//DeleteEnrollee - Deletes enrollee data
func DeleteEnrollee(c *gin.Context) {
	enrolleeID := c.Param("enrolleeId")

	result, err := collection.DeleteOne(context.TODO(), bson.M{"id": enrolleeID})
	if err != nil {
		log.Printf("Error while deleting the enrollee, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Internal error",
		})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "No Enrollee is deleted",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Enrollee is deleted successfully",
		})
	}

	return
}
