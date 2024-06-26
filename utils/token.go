package utils

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/bontusss/colosach/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// CreateToken is a function that generates a JWT token with the given time to live, payload, and private key.
// It returns the generated token and any error that occurred during the process.
// Example:
// token, err := CreateToken(time.Hour, "payload", "private_key")

func CreateToken(ttl time.Duration, payload interface{}, privateKey string) (string, error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("could not decode key: %w", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)

	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["sub"] = payload
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)

	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

func ValidateToken(token string, publicKey string) (interface{}, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		log.Println("error validate token 1: could not decode public key", err)
		return nil, fmt.Errorf("could not decode: %w", err)
	}

	log.Printf("Decoded Public Key: %s\n", string(decodedPublicKey))

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		log.Println("error validate token 2: failed to parse public key", err)
		return nil, fmt.Errorf("validate1: parse key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		log.Println("error validate token 3: token parsing or signature verification failed", err)
		return nil, fmt.Errorf("validate2: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate3: invalid token")
	}

	return claims["sub"], nil
}

func SetToken(user *models.DBResponse, ctx *gin.Context) {
	tokenDuration, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_EXPIRED_IN"))
	if err != nil {
		log.Fatal("Error parsing duration", err)
	}
	accessToken, err := CreateToken(tokenDuration, user.ID, os.Getenv("ACCESS_TOKEN_PRIVATE_KEY"))
	if err != nil {
		log.Println("error settoken 1: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refreshDuration, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_EXPIRED_IN"))
	if err != nil {
		log.Fatal("Error parsing duration", err)
	}
	refreshToken, err := CreateToken(refreshDuration, user.ID, os.Getenv("REFRESH_TOKEN_PRIVATE_KEY"))
	if err != nil {
		log.Println("error settoken 2: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	atma, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_MAXAGE"))
	if err != nil {
		log.Fatal("Error settoken 3: ", err)
	}

	rtma, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_MAXAGE"))
	if err != nil {
		log.Fatal("Error settoken 4: ", err)
	}
	ctx.SetCookie("access_token", accessToken, atma*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, rtma*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", atma*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}
