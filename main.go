/*

Written by Ikwechegh `Ukandu <ikwecheghu@gmail.com>
Date: 9th september 2023

*/

package main

import (
	"context"
	"fmt"
	"github.com/bontusss/colosach/config"
	"github.com/bontusss/colosach/controllers"
	"github.com/bontusss/colosach/routes"
	"github.com/bontusss/colosach/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"html/template"
	"log"
	"net/http"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoclient *mongo.Client
	//redisclient *redis.Client

	userService         services.UserService
	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	authCollection      *mongo.Collection
	authService         services.AuthService
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	temp *template.Template
)

func init() {
	temp = template.Must(template.ParseGlob("templates/*.html"))
	loadConfig, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	ctx = context.TODO()

	// Connect to MongoDB
	mongoconn := options.Client().ApplyURI(loadConfig.DBUri)
	mongoclient, err := mongo.Connect(ctx, mongoconn)

	if err != nil {
		panic(err)
	}

	if err := mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB successfully connected...")

	// Connect to Redis
	//redisclient = redis.NewClient(&redis.Options{
	//	Addr: loadConfig.RedisUri,
	//})
	//
	//if _, err := redisclient.Ping(ctx).Result(); err != nil {
	//	panic(err)
	//}
	//
	//err = redisclient.Set(ctx, "test", "Welcome to Golang with Redis and MongoDB", 0).Err()
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("Redis client connected successfully...")

	// Collections
	authCollection = mongoclient.Database("golang_mongodb").Collection("users")
	userService = services.NewUserServiceImpl(authCollection, ctx)
	authService = services.NewAuthService(authCollection, ctx)
	AuthController = controllers.NewAuthController(authService, userService, ctx, authCollection, temp)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(userService)
	UserRouteController = routes.NewRouteUserController(UserController)

	server = gin.Default()
}

func main() {
	loadConfig, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load config", err)
	}

	defer func(mongoclient *mongo.Client, ctx context.Context) {
		err := mongoclient.Disconnect(ctx)
		if err != nil {
			log.Println("error disconnecting mongo")
		}
	}(mongoclient, ctx)

	//value, err := redisclient.Get(ctx, "test").Result()
	//
	//if errors.Is(err, redis.Nil) {
	//	fmt.Println("key: test does not exist")
	//} else if err != nil {
	//	panic(err)
	//}

	corsConfig := cors.DefaultConfig()
	//todo update origins
	corsConfig.AllowOrigins = []string{loadConfig.Origin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/health-checker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "alive"})
	})

	router.GET("/search", services.SearchPexel)

	AuthRouteController.AuthRoute(router, userService)
	UserRouteController.UserRoute(router, userService)
	log.Fatal(server.Run(":" + loadConfig.Port))
}
