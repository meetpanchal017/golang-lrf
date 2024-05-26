package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang-jwt-demo/database"
	helper "golang-jwt-demo/helpers"
	"golang-jwt-demo/models"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	// "github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection

// var validate = validator.New()

func init() {
	userCollection = database.OpenCollection(database.Client, "User")
}

var appTimeOut = time.Second * 100

func HashPassword(password string) (string, error) {
	hashValue, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error generating password hash:", err)
		return "", err
	}
	return string(hashValue), nil
}

func VerifyPassword(userPassword string, providedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	if err != nil {
		return false, errors.New("email or password is incorrect")
	}
	return true, nil
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := helper.CheckUserType(c, "ADMIN")
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Meta: models.Meta{Message: err.Error()}})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), appTimeOut)
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.DefaultQuery("recordPerPage", "10"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage

		matchStage := bson.D{
			{Key: "$match", Value: bson.D{}},
		}

		groupStage := bson.D{
			primitive.E{Key: "$group", Value: bson.D{
				primitive.E{Key: "_id", Value: nil},
				primitive.E{Key: "total_count", Value: bson.D{
					primitive.E{Key: "$sum", Value: 1},
				}},
				primitive.E{Key: "data", Value: bson.D{
					primitive.E{Key: "$push", Value: "$$ROOT"},
				}},
			}},
		}

		projectStage := bson.D{
			primitive.E{Key: "$project", Value: bson.D{
				primitive.E{Key: "_id", Value: 0},
				primitive.E{Key: "total_count", Value: 1},
				primitive.E{Key: "user_items", Value: bson.D{
					primitive.E{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}},
				}},
			}},
		}

		cursor, err := userCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Meta: models.Meta{Message: "error occurred while listing users"}})
			return
		}
		defer cursor.Close(ctx)

		var allUsers []bson.M
		if err := cursor.All(ctx, &allUsers); err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{Meta: models.Meta{Message: err.Error()}})
			return
		}

		c.JSON(http.StatusOK, allUsers)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.Param("uid")
		println("MeetDebug check user id here: ", uid)

		ctx, cancel := context.WithTimeout(context.Background(), appTimeOut)
		defer cancel()

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"uid": uid}).Decode(&user)
		if err != nil {
			println("MeetDebug check - user exist ")
			c.JSON(http.StatusInternalServerError, models.Response{Meta: models.Meta{Message: err.Error()}})
			return
		}
		println("MeetDebug check - out side function ")
		c.JSON(http.StatusOK, models.Response{Data: user})
	}
}

func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body == nil {
			c.JSON(http.StatusOK, models.Response{Meta: models.Meta{Message: "Please pass the fields which you want to update"}})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), appTimeOut)
		defer cancel()

		var userProfile map[string]interface{}
		if err := json.NewDecoder(c.Request.Body).Decode(&userProfile); err != nil {
			c.JSON(http.StatusOK, models.Response{Meta: models.Meta{Message: "Something went wrong"}})
			return
		}

		token := c.Request.Header.Get("Authorization")

		// Define a whitelist of allowed fields
		allowedFields := map[string]bool{
			"first_name":   true,
			"last_name":    true,
			"phone_number": true,
			// Add more allowed fields as needed
		}

		update := bson.M{}
		for key, value := range userProfile {
			if allowedFields[key] {
				update[key] = value
			}
		}

		if len(update) == 0 {
			c.JSON(http.StatusOK, models.Response{Meta: models.Meta{Message: "No valid fields to update"}})
			return
		}

		err := updateUserFields(ctx, token, update)

		if err != nil {
			c.JSON(http.StatusOK, models.Response{Meta: models.Meta{Message: err.Error()}})
			return
		}

		c.JSON(http.StatusOK, models.Response{Meta: models.Meta{Message: "User updated successfully"}})
	}
}

func UpdateProfileImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		ctx, cancel := context.WithTimeout(context.Background(), appTimeOut)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file from request"})
			return
		}
		defer file.Close()

		bucket, err := gridfs.NewBucket(
			database.Client.Database("image-server"),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create GridFS bucket"})
			return
		}

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read file data"})
			return
		}

		filename := time.Now().Format(time.RFC3339) + "_" + header.Filename
		uploadStream, err := bucket.OpenUploadStream(filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload stream"})
			return
		}
		defer uploadStream.Close()

		_, err = uploadStream.Write(buf.Bytes())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file data to GridFS"})
			return
		}

		fileID := uploadStream.FileID.(primitive.ObjectID)

		imageURL := fmt.Sprintf("/uploads/%s", fileID.Hex()) // Relative path

		token := c.Request.Header.Get("Authorization")
		update := bson.M{"$set": bson.M{"profile_image": imageURL}}
		err = updateUserFields(ctx, token, update)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Profile picture updated successfully", "image_url": imageURL})
	}
}

func updateUserFields(ctx context.Context, token string, fieldsToUpdate bson.M) error {
	// Include the current time in the update operation
	fieldsToUpdate["updated_at"] = time.Now()

	_, err := userCollection.UpdateOne(ctx, bson.M{"token": token}, bson.M{"$set": fieldsToUpdate})
	if err != nil {
		return err
	}
	return nil
}

func GetImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		imageId := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(imageId)
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		bucket, err := gridfs.NewBucket(
			database.Client.Database("image-server"),
		)
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		var buf bytes.Buffer
		dStream, err := bucket.DownloadToStream(objID, &buf)
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		log.Printf("File size to download: %v\n", dStream)
		contentType := http.DetectContentType(buf.Bytes())

		c.Writer.Header().Add("Content-Type", contentType)
		c.Writer.Header().Add("Content-Length", strconv.Itoa(len(buf.Bytes())))

		c.Writer.Write(buf.Bytes())
	}
}
