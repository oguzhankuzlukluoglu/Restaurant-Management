package controllers

import (
	"RESTAURANT-MANAGEMENT/database"
	"RESTAURANT-MANAGEMENT/models"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"strconv"
	"time"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10

		}

		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}
		startIndex := (page - 1) * recordPerPage

		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
			}},
		}

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, projectStage,
		})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
			var allUsers [0]bson.M
			if err = result.All(ctx, &allUsers); err != nil {
				log.Fatal(err)
			}
			c.JSON(http.StatusOK, allUsers[])
		}

		//either pass an error

		//ideally want to return all the users based on the various query parameters

	}
}
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		userId := c.Param("user_id")

		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})

		}
		c.JSON(http.StatusOK, user)

	}
}
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		//convert to JSON data coming from postman a something that golang understands
		//validate the data based on user struct
		//you'll check if the email has already been used by another user
		//hash password
		//you'll also check if the phone number has already been used by another user
		//get some extra details for the user object like - created_at,updated_at,ID
		//generate token and refresh token (generate all tokens function from helper)
		//if all okey, than you insert this new user into the user collection
		//return status OK and send the result back

	}
}
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		//convert the login data from postman which is JSON to golang readable format
		//find a user with that email and see if that user even exist
		//then you will verify the password
		//if all goes well, then you'll generate tokens
		//update tokens - token and refreh token
		//return statusOK

	}
}
func HashPassword(password string) string {
	
}
func VerifyPassword(password string, providePassword string) (bool, string) {

}
