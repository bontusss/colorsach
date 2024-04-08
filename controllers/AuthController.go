package controllers

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bontusss/colosach/models"
	"github.com/bontusss/colosach/services"
	"github.com/bontusss/colosach/utils"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthController struct {
	authService services.AuthService
	userService services.UserService
	ctx         context.Context
	collection  *mongo.Collection
	temp        *template.Template
}

func NewAuthController(authService services.AuthService, userService services.UserService, ctx context.Context, collection *mongo.Collection, temp *template.Template) AuthController {
	return AuthController{authService, userService, ctx, collection, temp}
}

func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var user *models.SignUpInput

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if user.Password != user.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	newUser, err := ac.authService.SignUpUser(user)
	fmt.Println("User registered")

	if err != nil {
		if strings.Contains(err.Error(), "email already exist") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "error", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	//generate verification code
	fmt.Println("Generating code")
	verificationCode := randstr.String(20)
	_, err = ac.collection.UpdateOne(ac.ctx, bson.M{"email": newUser.Email}, bson.M{"$set": bson.M{"verificationCode": verificationCode}})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("User Updated")
	var firstName = newUser.Name
	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	//send  email verification mail
	emailData := utils.EmailData{
		URL:       os.Getenv("PORT") + "/verify-email/" + verificationCode,
		FirstName: firstName,
		Subject:   "Your Colosach verification code",
	}
	fmt.Println("code generated")
	err = utils.SendEmail(newUser, &emailData, "verificationCode.html")
	fmt.Println("Starting to send mail")
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "success", "message": "There was an error sending email"})
		return
	}
	fmt.Println(verificationCode)
	fmt.Println(emailData)
	message := "We sent an email with a verification code to " + user.Email
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": message})
}

// SignInUser => Logging in user
func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var credentials *models.SignInInput

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user, err := ac.userService.FindUserByEmail(credentials.Email)

	// check if user email is verified
	if !user.Verified {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User email is not verified."})
		return
	}
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or password"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if err := utils.VerifyPassword(user.Password, credentials.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	// Generate Tokens
	tokenDuration, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRED_IN"))
	if err != nil {
		log.Fatal("Error parsing duration", err)
	}
	accessToken, err := utils.CreateToken(tokenDuration, user.ID, os.Getenv("ACCESS_TOKEN_PRIVATE_KEY"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refreshDuration, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_EXPIRED_IN"))
	if err != nil {
		log.Fatal("Error parsing duration", err)
	}
	refreshToken, err := utils.CreateToken(refreshDuration, user.ID, os.Getenv("REFRESH_TOKEN_PRIVATE_KEY"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	atma, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_MAXAGE"))
	if err != nil {
		log.Fatal("Error: ", err)
		return
	}

	rtma, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_MAXAGE"))
	if err != nil {
		log.Fatal("Error: ", err)
		return
	}
	ctx.SetCookie("access_token", accessToken, atma*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, rtma*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", atma*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	message := "could not refresh access token"

	cookie, err := ctx.Cookie("refresh_token")

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	sub, err := utils.ValidateToken(cookie, os.Getenv("REFRESH_TOKEN_PUBLIC_KEY"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user, err := ac.userService.FindUserById(fmt.Sprint(sub))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		return
	}

	tokenDuration, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRED_IN"))
	if err != nil {
		log.Fatal("Error parsing duration", err)
	}

	accessToken, err := utils.CreateToken(tokenDuration, user.ID, os.Getenv("ACCESS_TOKEN_PRIVATE_KEY"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	atma, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_MAXAGE"))
	if err != nil {
		log.Fatal("Error: ", err)
		return
	}

	ctx.SetCookie("access_token", accessToken, atma*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", atma*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

// Verify Email

func (ac *AuthController) VerifyEmail(ctx *gin.Context) {
	verificationCode := ctx.Params.ByName("verificationCode")
	//verificationCode := utils.Encode(code)

	filter := bson.D{{Key: "verificationCode", Value: verificationCode}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "verified", Value: true}}}, {Key: "$unset", Value: bson.D{{Key: "verificationCode", Value: ""}}}}
	result, err := ac.collection.UpdateOne(ac.ctx, filter, update)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "success", "message": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "success", "message": "Could not verify email address"})
		return
	}

	fmt.Println(result)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Email verified successfully"})

}

// ForgotPassword => Forgot Password
func (ac *AuthController) ForgotPassword(ctx *gin.Context) {
	var userCredential *models.ForgotPasswordInput

	if err := ctx.ShouldBindJSON(&userCredential); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	message := "You will receive a reset email if user with that email exist"

	user, err := ac.userService.FindUserByEmail(userCredential.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusOK, gin.H{"status": "fail", "message": message})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if !user.Verified {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Account not verified"})
		return
	}

	// Generate Verification Code
	passwordResetToken := randstr.String(20)

	// Update User in Database
	query := bson.D{{Key: "email", Value: strings.ToLower(userCredential.Email)}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "passwordResetToken", Value: passwordResetToken}, {Key: "passwordResetAt", Value: time.Now().Add(time.Minute * 15)}}}}
	result, err := ac.collection.UpdateOne(ac.ctx, query, update)

	if result.MatchedCount == 0 {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "success", "message": "There was an error sending email"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "success", "message": err.Error()})
		return
	}
	var firstName = user.Name

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// ? Send Email
	emailData := utils.EmailData{
		URL:       os.Getenv("CLIENT_ORIGIN") + "/reset-password/" + passwordResetToken,
		FirstName: firstName,
		Subject:   "Your password reset token (valid for 10min)",
	}

	err = utils.SendEmail(user, &emailData, "resetPassword.html")
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "success", "message": "There was an error sending email"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
}

func (ac *AuthController) ResetPassword(ctx *gin.Context) {
	passwordResetToken := ctx.Params.ByName("resetToken")
	var userCredential *models.ResetPasswordInput

	if err := ctx.ShouldBindJSON(&userCredential); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if userCredential.Password != userCredential.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	hashedPassword, _ := utils.HashPassword(userCredential.Password)

	// Update User in Database
	query := bson.D{{Key: "passwordResetToken", Value: passwordResetToken}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: hashedPassword}}}, {Key: "$unset", Value: bson.D{{Key: "passwordResetToken", Value: ""}, {Key: "passwordResetAt", Value: ""}}}}
	result, err := ac.collection.UpdateOne(ac.ctx, query, update)

	if result.MatchedCount == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "success", "message": "Token is invalid or has expired"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "success", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Password data updated successfully"})
}
