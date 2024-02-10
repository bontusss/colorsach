# Colosach Backend Service
This is the service that powers the colosach website written in golang language and the gin framework.

## Stacks
1. [Gin Gonic](https://github.com/gin-gonic/gin) => Web framework
2. [MongoDB](https://mongodb.org) => Database
3. [JWT](https://github.com/golang-jwt/jwt/v5) => Authentication and authorization
4. [Cloudinary](https://cloudinary.com) => Image management services

## Endpoints
1. **GET** `/api/search` => Takes a name and color params and returns images from pexels containing the given name and color
2. **GET** `/api/health-checker` => Returns alive if server is up and dead if down (Redis is not set up for this feature, refactor please)
3. **POST** `/api/auth/register` => requires an email, name, password, passwordConfirm params, returns a 201 and sends an email to user email if successful, returns an error message if not successful.
4. **POST** `/api/auth/login `=> requires an email and password and returns a success message and a token if successful, returns an error message otherwise.
5. **GET** `/api/auth/refresh` => Returns a refresh token if successful, an error message otherwise.
6. **GET** `/api/auth/verify-email/:verificationCode` => Takes the code sent to the users email as params, returns a 200 and a message if successful, returns an error message otherwise.
7. **POST** `/api/auth/forgot-password` => Requires an email, returns a 201 and success message if successful, an error message otherwise.
8. **PATCH** `/api/auth/reset-password/:resetToken`
9. **GET** `/api/auth/logout`
