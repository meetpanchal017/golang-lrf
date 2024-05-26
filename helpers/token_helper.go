package helper

import (
	"context"
	"golang-jwt-demo/database"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "User")

type SignedDetails struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Uid       string `json:"uid"`
	UserType  string `json:"user_type"`
	jwt.StandardClaims
}

var SECRET_KEY []byte

func init() {
	// Load secret key from environment variable
	SECRET_KEY = []byte(os.Getenv("SECRET_KEY"))
}

func GenerateAllTokens(email, firstName, lastName, userType, uid string) (signedToken, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		UserType:  userType,
		Uid:       uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 168).Unix(), // Expires in 7 days
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(SECRET_KEY)
	if err != nil {
		log.Panic(err)
		return "", "", err
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(SECRET_KEY)
	if err != nil {
		log.Panic(err)
		return "", "", err
	}
	return token, refreshToken, err
}

func UpdateAllTokens(signedToken, signedRefreshToken, userID string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	updateObj := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "refresh_token", Value: signedRefreshToken},
			{Key: "updated_at", Value: time.Now()},
		}},
	}

	upsert := true
	filter := bson.M{"user_id": userID}
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
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		msg = err.Error()
		return nil, msg
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok || !token.Valid {
		msg = "token is invalid"
		return nil, msg
	}

	if claims.ExpiresAt < time.Now().Unix() {
		msg = "token expired"
		return nil, msg
	}

	return claims, ""
}

type SignedDetailsNew struct {
	PhoneNumber string `json:"phone_number"`
	Uid         string `json:"uid"`
	jwt.StandardClaims
}

func GenerateTokensNew(phoneNumber, uid string) (signedToken, signedRefreshToken string, err error) {
	// Create access token claims
	accessTokenClaims := SignedDetailsNew{
		PhoneNumber: phoneNumber,
		Uid:         uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours
		},
	}

	// Create refresh token claims
	refreshTokenClaims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 168).Unix(), // Expires in 7 days
	}

	// Generate access token
	accessToken, err := generateToken(accessTokenClaims)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshToken, err := generateToken(refreshTokenClaims)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func generateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		log.Panic(err)
		return "", err
	}
	return signedToken, nil
}
