package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	guuid "github.com/google/uuid"
	"github.com/minateegithub/go_microservice/models"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
)

//AddDependent - Adds new dependents to an enrollee
func AddDependent(c *gin.Context) {
	enrolleeID := c.Param("enrolleeId")
	var dependents []models.Dependent

	err2 := c.BindJSON(&dependents)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	newDependents := []models.Dependent{}

	for i := 0; i < len(dependents); i++ {

		validationErrors := validateNewDependentData(dependents[i])
		if len(validationErrors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}

		name := dependents[i].Name
		birthDate := dependents[i].BirthDate
		id := guuid.New().String()
		newDependent := models.Dependent{
			ID:        &id,
			Name:      name,
			BirthDate: birthDate,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		newDependents = append(newDependents, newDependent)
	}

	_, err := collection.UpdateOne(context.TODO(), bson.M{"id": enrolleeID}, bson.M{"$push": bson.M{"dependents": bson.M{"$each": newDependents}}})

	if err != nil {
		log.Printf("Error while adding new dependents into db, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Internal error",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Dependents are added successfully",
	})
	return
}

//EditDependent - Updates dependent data
func EditDependent(c *gin.Context) {
	enrolleeID := c.Param("enrolleeId")
	dependentID := c.Param("dependentId")

	var dependent models.Dependent

	err2 := c.BindJSON(&dependent)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	validationErrors := validateEditedDependentData(dependent)
	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	currentEnrollee := models.Enrollee{}
	err := collection.FindOne(context.TODO(), bson.M{"id": enrolleeID, "dependents.id": dependentID}).Decode(&currentEnrollee)
	if err != nil {
		log.Printf("Error while getting the enrollee, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Enrollee not found",
		})
		return
	}
	/* s, _ := json.MarshalIndent(currentEnrollee, "", "\t")
	fmt.Println(string(s)) */

	currentDependent := models.Dependent{}

	for i := 0; i < len(currentEnrollee.Dependents); i++ {
		if currentEnrollee.Dependents[i].ID != nil && *currentEnrollee.Dependents[i].ID == dependentID {
			currentDependent = currentEnrollee.Dependents[i]
		}
	}

	/* s1, _ := json.MarshalIndent(currentDependent, "", "\t")
	fmt.Println(string(s1)) */

	if dependent.Name != nil {
		currentDependent.Name = dependent.Name
	}
	if dependent.BirthDate != nil {
		currentDependent.BirthDate = dependent.BirthDate
	}

	newData := bson.M{
		"$set": bson.M{
			"dependents.$": bson.M{
				"id":         currentDependent.ID,
				"name":       currentDependent.Name,
				"birth_date": currentDependent.BirthDate,
				"created_at": currentDependent.CreatedAt,
				"updated_at": time.Now(),
			},
		},
	}

	_, err1 := collection.UpdateOne(context.TODO(), bson.M{"id": enrolleeID, "dependents.id": dependentID}, newData)
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
		"message": "Dependent Edited Successfully",
	})
	return
}

//DeleteDependent - Deletes dependent data
func DeleteDependent(c *gin.Context) {
	enrolleeID := c.Param("enrolleeId")
	dependentID := c.Param("dependentId")

	result, err := collection.UpdateOne(context.TODO(), bson.M{"id": enrolleeID}, bson.M{"$pull": bson.M{"dependents": bson.M{"id": dependentID}}})
	if err != nil {
		log.Printf("Error while deleting the dependent, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Internal error",
		})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "No Dependent is deleted",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Dependent is deleted successfully",
		})
	}

	return
}

//GetDependents - Gets all the dependents for an enrollee
func GetDependents(c *gin.Context) {
	enrolleeID := c.Param("enrolleeId")

	enrollee := models.Enrollee{}
	opts := options.FindOne().SetProjection(bson.M{"dependents": 1, "_id": 0})
	err := collection.FindOne(context.TODO(), bson.M{"id": enrolleeID}, opts).Decode(&enrollee)
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
		"message": "Dependents",
		"data":    enrollee.Dependents,
	})
	return
}
