package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func getTransactions(c *gin.Context) {

	var transactions []Transaction

	cursor, err := db.Find(context.Background(), bson.M{})

	if err != nil {
		log.Fatal("Error searching in database: ", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to find transactions"})
		return
	}

	defer cursor.Close(context.Background())

	if !cursor.Next(ctx) {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	if err := cursor.All(context.Background(), &transactions); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to deserialize transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transactions obtained successfully",
		"body":    transactions})
}

func getTransaction(c *gin.Context) {

	idParam := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(idParam)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var transaction Transaction

	if err := db.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&transaction); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction obtained successfully",
		"body":    transaction})
}

func createTransaction(c *gin.Context) {

	var transaction Transaction

	if err := c.BindJSON(&transaction); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	transaction.CreatedAt = time.Now()

	result, err := db.InsertOne(context.Background(), transaction)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	transaction.ID = result.InsertedID

	c.JSON(http.StatusCreated, gin.H{
		"message": "Transaction created succesfully",
		"body":    transaction,
	})

}

func updateTransaction(c *gin.Context) {

	idParam := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(idParam)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var updatedTransaction Transaction
	if err := c.BindJSON(&updatedTransaction); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	updatedTransaction.UpdatedAt = time.Now()

	update := bson.M{"$set": updatedTransaction}

	result, err := db.UpdateOne(context.Background(), bson.M{"_id": objectId}, update)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction"})
		return
	}

	if result.MatchedCount == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	updatedTransaction.ID = result.UpsertedID

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction updated succesfully",
		"body":    updatedTransaction,
	})

}

func deleteTransaction(c *gin.Context) {

	idParam := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(idParam)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	result, err := db.DeleteOne(context.Background(), bson.M{"_id": objectId})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction"})
		return
	}

	if result.DeletedCount == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
