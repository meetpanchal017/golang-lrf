package controller

import (
	"context"
	helper "golang-jwt-demo/helpers"
	"golang-jwt-demo/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// func Signup() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var user models.User
// 		if err := c.BindJSON(&user); err != nil {
// 			c.JSON(http.StatusBadRequest, models.Response{Meta: models.Meta{Message: err.Error()}})
// 			return
// 		}

// 		if err := validate.Struct(user); err != nil {
// 			c.JSON(http.StatusBadRequest, models.Response{Meta: models.Meta{Message: err.Error()}})
// 			return
// 		}

// 		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
// 		defer cancel()

// 		emailCount, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, models.Response{Meta: models.Meta{Message: "Error occurred while checking for the email"}})
// 			return
// 		}

// 		password, err := HashPassword(*user.Password)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, models.Response{Meta: models.Meta{Message: "Error occurred while hashing password"}})
// 			return
// 		}
// 		user.Password = &password

// 		phoneCount, err := userCollection.CountDocuments(ctx, bson.M{"Phone": user.PhoneNumber})
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, models.Response{Meta: models.Meta{Message: "Error occurred while checking for the phone"}})
// 			return
// 		}

// 		if emailCount > 0 || phoneCount > 0 {
// 			c.JSON(http.StatusConflict, models.Response{Meta: models.Meta{Message: "This email or phone already exists"}})
// 			return
// 		}

// 		user.CreatedAt = time.Now()
// 		user.UpdatedAt = time.Now()
// 		user.Id = primitive.NewObjectID()
// 		user.UserID = user.Id.Hex()

// 		token, refreshToken, err := helper.GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, *user.UserType, user.UserID)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, models.Response{Meta: models.Meta{Message: err.Error()}})
// 			return
// 		}
// 		user.Token = &token
// 		user.RefreshToken = &refreshToken

// 		_, err = userCollection.InsertOne(ctx, user)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, models.Response{Meta: models.Meta{Message: "User was not created"}})
// 			return
// 		}

// 		c.JSON(http.StatusOK, models.Response{Meta: models.Meta{Message: "User signed up successfully"}})
// 	}
// }

// func Login() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var user models.User
// 		if err := c.BindJSON(&user); err != nil {
// 			c.JSON(http.StatusBadRequest, models.Response{Meta: models.Meta{Message: err.Error()}})
// 			return
// 		}

// 		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
// 		defer cancel()

// 		var foundUser models.User
// 		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, models.Response{Meta: models.Meta{Message: "Email or password is incorrect"}})
// 			return
// 		}

// 		passIsValid, err := VerifyPassword(*user.Password, *foundUser.Password)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, models.Response{Meta: models.Meta{Message: err.Error()}})
// 			return
// 		}

// 		if !passIsValid {
// 			c.JSON(http.StatusBadRequest, models.Response{Meta: models.Meta{Message: "email or password is incorrect"}})
// 			return
// 		}

// 		if foundUser.Email == nil {
// 			c.JSON(http.StatusBadRequest, models.Response{Meta: models.Meta{Message: "user not found"}})
// 			return
// 		}

// 		token, refreshToken, err := helper.GenerateAllTokens(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, *foundUser.UserType, foundUser.UserID)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, models.Response{Meta: models.Meta{Message: err.Error()}})
// 			return
// 		}
// 		helper.UpdateAllTokens(token, refreshToken, foundUser.UserID)

// 		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.UserID}).Decode(&foundUser)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, models.Response{Meta: models.Meta{Message: err.Error()}})
// 			return
// 		}

// 		c.JSON(http.StatusOK, foundUser)
// 	}
// }

// func SendOTP() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		_, cancel := context.WithTimeout(context.Background(), appTimeOut)
// 		defer cancel()

// 		var requestBody struct {
// 			PhoneNumber string `json:"phone_number"`
// 		}

// 		if err := c.BindJSON(&requestBody); err != nil {
// 			c.JSON(http.StatusBadRequest, models.Response{Meta: models.Meta{Message: "Invalid request body"}})
// 			return
// 		}

// 		phoneNumber := requestBody.PhoneNumber

// 		fmt.Print("MeetDebug check phone number here - ", phoneNumber)
// 		_, err := helper.TwilioSendOTP(phoneNumber)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, models.Response{Meta: models.Meta{Message: err.Error()}})
// 			return
// 		}
// 		c.JSON(http.StatusOK, models.Response{Meta: models.Meta{Message: "Otp sent successfully"}})
// 	}
// }

// func VerifyOTP() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		_, cancel := context.WithTimeout(context.Background(), appTimeOut)
// 		defer cancel()

// 		var requestBody struct {
// 			PhoneNumber string `json:"phone_number"`
// 			Code        string `json:"code"`
// 		}

// 		if err := c.BindJSON(&requestBody); err != nil {
// 			c.JSON(http.StatusBadRequest, models.Response{Meta: models.Meta{Message: "Invalid request body"}})
// 			return
// 		}

// 		phoneNumber := requestBody.PhoneNumber
// 		code := requestBody.Code

// 		err := helper.TwilioVerifyOTP(phoneNumber, code)

// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, models.Response{Meta: models.Meta{Message: err.Error()}})
// 			return
// 		}
// 		// phoneCount, err := userCollection.CountDocuments(ctx, bson.M{"Phone": user.Phone})
// 		// if err != nil {
// 		// 	c.JSON(http.StatusInternalServerError, models.Response{Meta: models.Meta{Message: "Error occurred while checking for the phone"}})
// 		// 	return
// 		// }
// 		c.JSON(http.StatusOK, models.Response{Meta: models.Meta{Message: "Otp verified successfully"}})
// 	}
// }

func LoginWithPhoneNumber() gin.HandlerFunc {
	return func(c *gin.Context) {

		_, cancel := context.WithTimeout(context.Background(), appTimeOut)
		defer cancel()
		var requestBody struct {
			PhoneNumber string `json:"phone_number"`
		}
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Meta: models.Meta{Message: "Invalid request body"}})
			return
		}
		_, err := helper.TwilioSendOTP(requestBody.PhoneNumber)

		if err != nil {
			c.JSON(http.StatusBadRequest, models.Response{Meta: models.Meta{Message: err.Error()}})
			return
		}
		c.JSON(http.StatusOK, models.Response{Meta: models.Meta{Message: "Otp sent successfully"}})

	}
}

func VerifyOTPNew() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), appTimeOut)
		defer cancel()

		var requestBody struct {
			PhoneNumber string `json:"phone_number"`
			Code        string `json:"code"`
		}

		err := c.BindJSON(&requestBody)
		if handleServerError(c, err) {
			return
		}

		phoneNumber := requestBody.PhoneNumber
		code := requestBody.Code

		err = helper.TwilioVerifyOTP(phoneNumber, code)
		if handleServerError(c, err) {
			return
		}

		phoneCount, err := userCollection.CountDocuments(c, bson.M{"phone_number": requestBody.PhoneNumber})
		if handleServerError(c, err) {
			return
		}

		var user models.User

		if phoneCount > 0 {
			println("MeetDebug - User exist - ", user.PhoneNumber)
			err := userCollection.FindOne(c, bson.M{"phone_number": user.PhoneNumber}).Decode(&user)

			if handleServerError(c, err) {
				return
			}
		} else {
			user.Id = primitive.NewObjectID()
			user.Uid = user.Id.Hex()
			user.CreatedAt = time.Now()
			user.UpdatedAt = time.Now()
			user.PhoneNumber = &requestBody.PhoneNumber
			token, refreshToken, err := helper.GenerateTokensNew(*user.PhoneNumber, user.Uid)
			if handleServerError(c, err) {
				return
			}
			user.Token = &token
			user.RefreshToken = &refreshToken
			_, err = userCollection.InsertOne(c, user)
			if handleServerError(c, err) {
				return
			}

		}
		c.JSON(http.StatusOK, models.Response{Data: user, Meta: models.Meta{Message: "User Login successfully!"}})
	}
}

func handleServerError(c *gin.Context, err error) bool {
	if err != nil {
		message := "Oops! Something went wrong"
		if err.Error() != "" {
			message = err.Error()
		}
		c.JSON(http.StatusInternalServerError, models.Response{Meta: models.Meta{Message: message}})
		return true
	}
	return false
}
