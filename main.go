/*

Written by Ikwechegh `Ukandu <ikwecheghu@gmail.com: https://github.com/bontusss>
Date: 9th september 2023

*/

package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/bontusss/colosach/controllers"
	"github.com/bontusss/colosach/routes"
	"github.com/bontusss/colosach/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoclient *mongo.Client
	//redisclient *redis.Client

	// for users
	userService         services.UserService
	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	// for auths
	authCollection      *mongo.Collection
	authService         services.AuthService
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	// for libraries
	libService         services.LibraryService
	libController      controllers.LibraryController
	libCollection      *mongo.Collection
	libRouteController routes.LibraryRouteController

	// for Images
	imageService         services.ImageService
	ImageController      controllers.ImageController
	imageCollection      *mongo.Collection
	ImageRouteController routes.ImageRouteController

	// Templates
	temp *template.Template
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading app.env file", err)
	}
	temp = template.Must(template.ParseGlob("templates/*.html"))
	ctx = context.TODO()
	// Connect to MongoDB
	fmt.Println("mongo DBurl is", os.Getenv("MONGODB_URI"))
	mongoconn := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
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
	authCollection = mongoclient.Database("colosach").Collection("users")
	userService = services.NewUserServiceImpl(authCollection, ctx)
	authService = services.NewAuthService(authCollection, ctx)
	AuthController = controllers.NewAuthController(authService, userService, ctx, authCollection, temp)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(userService)
	UserRouteController = routes.NewRouteUserController(UserController)

	libCollection = mongoclient.Database("colosach").Collection("libraries")
	libService = services.NewLibService(libCollection, ctx)
	libController = controllers.NewLibraryController(libService)
	libRouteController = routes.NewLibRouteController(libController)

	imageCollection = mongoclient.Database("colosach").Collection("images")
	imageService = services.NewImageService(imageCollection, ctx)
	ImageController = controllers.NewImageController(imageService)
	ImageRouteController = routes.NewImageRouteController(ImageController)

	server = gin.Default()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading app.env file")
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
	// corsConfig.AllowOrigins = []string{os.Getenv("CLIENT_ORIGIN"), os.Getenv("CLIENT_URL")}
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Authorization")
	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/health-checker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "alive"})
	})

	router.POST("/search", services.SearchPexel)

	AuthRouteController.AuthRoute(router, userService)
	UserRouteController.UserRoute(router, userService)
	libRouteController.LibraryRoute(router, userService)
	ImageRouteController.ImageRoute(router, userService)
	log.Fatal(server.Run(":" + os.Getenv("PORT")))
}
