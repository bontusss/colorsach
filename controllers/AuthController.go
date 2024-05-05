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

// login, token generation, email verification, password reset, and logout.
// It interacts with services and the database to handle these operations securely and efficiently.

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

// SignUpUser is a function that handles the user sign-up process.
// It first validates the user input, checks if the passwords match, and then proceeds to create a new user in the database.
// If the user is successfully created, it generates a verification code, updates the user's record in the database with the verification code,
// and sends a verification email to the user. If any error occurs during this process, it returns an appropriate error message.
// Example:
// POST /api/signup
// Request Body:
//
//	{
//	  "username": "newuser",
//	  "email": "newuser@example.com",
//	  "password": "password",
//	  "passwordConfirm": "password"
//	}
//
// Response:
//
//	{
//	  "status": "success",
//	  "message": "We sent an email with a verification code to newuser@example.com"
//	}
// @Summary Register User
// @Description Register a user
// @Tags auth
// @Accept json
// @Produce json
// @Param SignUpInput body models.SignUpInput true "SignUpInput"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/auth/register [post]
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

	err := ac.collection.FindOne(ac.ctx, bson.M{"username": user.Username}).Decode(&user)
	if err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "username already exists"})
		return
	}

	fmt.Println("signing up with: ", user.Email, user.Password, user.PasswordConfirm)

	newUser, err := ac.authService.SignUpUser(user)
	fmt.Println("User registered")

	if err != nil {
		if strings.Contains(err.Error(), "email already exist") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "error", "message": err.Error()})
			// log.Fatal(err)
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})

		return
	}

	if !user.Verified {
		//generate verification code
		// fmt.Println("Generating code")
		verificationCode := randstr.String(20)
		_, err = ac.collection.UpdateOne(ac.ctx, bson.M{"email": newUser.Email}, bson.M{"$set": bson.M{"verificationCode": verificationCode}})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			// log.Fatal(err)
		}

		randImage, err := services.GetRandomImage()
		if err != nil {
			// is it necessary to stop process because of this?
			log.Println("error getting image for email template")
		}

		//send  email verification mail
		emailData := utils.EmailData{
			URL:              os.Getenv("CLIENT_URL") + verificationCode,
			Username:         newUser.Username,
			Subject:          "Your Colosach Verification Code",
			Year:             time.Now().Year(),
			BannerImageUrl:   randImage.Src.Original,
			PhotographerName: randImage.Photographer,
			PhotographerUrl:  randImage.PhotographerURL,
			ImageUrl:         randImage.URL,
			AvgColor:         randImage.AvgColor,
			Logo:             "/assets/logo.png",
			Instagram:        "/assets/InstagramLogo.svg",
			Linkedin:         "https://www.flaticon.com/free-icon/linkedin_2504923?term=linkedin&page=1&position=11&origin=search&related_id=2504923",
			Arrow:            "/assets/ArrowRight.svg",
			X:                "/assets/X.svg",
		}
		// fmt.Println("code generated")
		err = utils.SendEmail(newUser, &emailData, "verificationCode.tmpl")
		// fmt.Println("Starting to send mail")
		if err != nil {
			fmt.Println(err)
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": "There was an error sending email"})
			// log.Fatal(err)
		}
		// fmt.Println(verificationCode)
		// fmt.Println(emailData)
		message := "We sent an email with a verification code to " + user.Email
		ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": message})

	} else {
		ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Registered successfully"})
	}

}
// SignInUser is a function to handle user sign in
// It takes a gin context as a parameter
// Example of how to use the function
//
//	func main() {
//		router := gin.Default()
//		authController := controllers.NewAuthController()
//		router.POST("/signup", authController.SignUpUser)
//		router.Run(":8080")
//	}
// @Summary SignInUser
// @Description SignInUser
// @Tags auth
// @Accept  json
// @Produce  json
// @Param user body models.SignInInput true "SignInInput"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/auth/login [post]
func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var credentials *models.SignInInput

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		log.Fatal(err)
	}

	fmt.Println("logging in with: ", credentials.Email, credentials.Password)

	user, err := ac.userService.FindUserByEmail(credentials.Email)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or password"})
			// log.Fatal(err)
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or password"})
		// log.Fatal(err)
	}

	// check if user email is verified
	if !user.Verified {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User email is not verified."})
		log.Fatal("user email is not verified")
	}


	fmt.Println(user.Password)
	fmt.Println(credentials.Password)
	// hashedPassword, _ := utils.HashPassword(credentials.Password)
	if err := utils.VerifyPassword(user.Password, credentials.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Passwordd"})
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
	}

	rtma, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_MAXAGE"))
	if err != nil {
		log.Fatal("Error: ", err)
	}
	ctx.SetCookie("access_token", accessToken, atma*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, rtma*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", atma*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

// @Summary RefreshAccessToken
// @Description RefreshAccessToken
// @Tags auth
// @Accept json
// @Produce json
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/auth/refresh [get]
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

// @Summary LogoutUser
// @Description LogoutUser
// @Tags auth
// @Accept json
// @Produce json
// @Success 200
// @Failure 401
// @Router /api/auth/logout [get]
func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

// Verify Email
// @Summary VerifyEmail
// @Description VerifyEmail
// @Tags auth
// @Accept json
// @Produce json
// @Param verificationCode path string true "Verification Code"
// @Success 200
// @Failure 403
// @Router /api/auth/verify-email/{verificationCode} [get]
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
// @Summary ForgotPassword
// @Description ForgotPassword
// @Tags auth
// @Accept json
// @Produce json
// @Param ForgotPasswordInput body models.ForgotPasswordInput true "ForgotPasswordInput"
// @Success 200
// @Failure 401
// @Failure 403
// @Failure 502
// @Router /api/auth/forgot-password [post]
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

	// ? Send Email
	emailData := utils.EmailData{
		URL:      os.Getenv("CLIENT_ORIGIN") + "/reset-password/" + passwordResetToken,
		Username: user.Username,
		Subject:  "Your password reset token (valid for 10min)",
	}

	err = utils.SendEmail(user, &emailData, "resetPassword.html")
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "success", "message": "There was an error sending email"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
}

// ResetPassword => Reset Password
// @Summary ResetPassword
// @Description ResetPassword
// @Tags auth
// @Accept json
// @Produce json
// @Param resetToken path string true "Reset Token"
// @Param ResetPasswordInput body models.ResetPasswordInput true "ResetPasswordInput"
// @Success 200
// @Failure 400
// @Failure 403
// @Failure 500
// @Router /api/auth/reset-password/{resetToken} [post]
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
