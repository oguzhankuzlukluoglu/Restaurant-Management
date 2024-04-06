package helpers

import (
	"RESTAURANT-MANAGEMENT/database"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	jwt.Claims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstName string, lastName string, uid string) (signedToken string, signedRefreshToken string, err error) {

	claims := &SignedDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		Uid:        uid,
		Claims: jwt.Claims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshClaims := &SignedDetails{
		Claims: jwt.Claims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}
	return token, refreshToken, err
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	updateObj := bson.D{
		{"$set", bson.D{
			{"token", signedToken},
			{"refresh_token", signedRefreshToken},
			{"updated_at", time.Now().Format(time.RFC3339)},
		}},
	}

	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(ctx, filter, updateObj, &opt)
	if err != nil {
		log.Panic(err)
		return
	}
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		msg = fmt.Sprintf("the token is invalid: %v", err)
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		return
	}

	if claims.ExpiresAt < time.Now().Unix() {
		msg = fmt.Sprintf("token is expired")
		return
	}

	return claims, msg
}
