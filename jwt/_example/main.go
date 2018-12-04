// iris provides some basic middleware, most for your learning courve.
// You can use any net/http compatible middleware with iris.FromStd wrapper.
//
// JWT net/http video tutorial for golang newcomers: https://www.youtube.com/watch?v=dgJFeqeXVKw
//
// This middleware is the only one cloned from external source: https://github.com/auth0/go-jwt-middleware
// (because it used "context" to define the user but we don't need that so a simple iris.FromStd wouldn't work as expected.)
package main

// $ go get -u github.com/dgrijalva/jwt-go
// $ go run main.go

import (
	"fmt"
	"time"

	"github.com/teamlint/iris"

	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/teamlint/middleware/jwt"
)

var jwtHandler *jwtmiddleware.Middleware

func myHandler(ctx iris.Context) {
	// user := ctx.Values().Get("jwt").(*jwt.Token)
	token := jwtHandler.Get(ctx)

	fmt.Println("myHandler")
	ctx.Writef("This is an authenticated request\n")
	ctx.Writef("Claim content:\n")

	var userID string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if val, ok := claims["user_id"]; ok {
			userID = val.(string)
		}
	} else {
		userID = "user_id is null"
	}
	var exp int64
	fmt.Println("claims exp get")
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Printf("claims[exp]:%v", claims["exp"])
		if val, ok := claims["exp"]; ok {
			exp = int64(val.(float64))
		}
	}
	if exp > 0 {
		expTime := time.Unix(exp, 0)
		ctx.Writef("exp:%v\n", expTime)
	}

	ctx.Writef("raw:%v\nsign:%s\nvalid:%v\nuser_id:%v", token.Raw, token.Signature, token.Valid, userID)
}

func main() {
	app := iris.New()
	key := []byte("My Secret")

	jwtHandler = jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return key, nil
		},
		// Expiration: true,
		Debug: false,
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
		ErrorHandler: func(ctx iris.Context, msg string) {
			data := struct {
				Success bool
				Msg     string
			}{
				Success: false,
				Msg:     "错误处理:" + msg,
			}
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(data)
		},
	})

	// app.Use(jwtHandler.Serve)

	app.Get("/", func(ctx iris.Context) {
		// Create a new token object, specifying signing method and the claims
		// you would like it to contain.
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": "1001",
			"exp":     time.Now().Add(time.Minute * 3).Unix(),
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, _ := token.SignedString(key)
		ctx.Application().Logger().Debugf("token:%v", tokenString)
		ctx.Header("Authorization", "Bearer "+tokenString)
		ctx.WriteString(tokenString)

	})
	app.Get("/ping", jwtHandler.Serve, myHandler)
	app.Run(iris.Addr("localhost:3001"))
} // don't forget to look ../jwt_test.go to seee how to set your own custom claims
