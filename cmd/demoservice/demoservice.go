package main

import (
	"demo/internal/data"
	"demo/internal/env"
	"demo/internal/logger"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

//loginHandler checks user's credentials, authenticates the user and then generates a jwt authorization token to the
//user with a configurable expiry date.
func loginHandler(ctx iris.Context) {
	userName := ctx.URLParam("username")
	passwordHash := ctx.URLParam("password")

	rejectReason := "bad user credentials"
	user, err := data.GetUser(userName, passwordHash)
	if err != nil {
		logger.Log.Sugar().Errorf("user validation failed with error %e", err)
		rejectReason = fmt.Sprintf("user validation faile with error %e", err)
	}

	err = nil
	if user.ID == 0 {
		_, err = ctx.JSON(iris.Map{
			"status":  http.StatusUnauthorized,
			"message": rejectReason,
		})
	} else {
		tokenExpiry, _ := strconv.Atoi(env.GetEnv("TOKEN_EXPIRY").(string))
		//tokenExpiry := int64(3600)
		token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username":  user.Name,
			"plan": user.Plan,
			"exp": time.Now().Unix() + int64(tokenExpiry),
		})

		tokenString, _ := token.SignedString([]byte("demoappdemoappdemoapp"))
		_, err = ctx.JSON(iris.Map{
			"status":  http.StatusOK,
			"message": "authentication successfully",
			"token":   tokenString,
		})

	}
	if err != nil {
		logger.Log.Sugar().Errorf("error passing json to context %e", err)
	}

}

//myAuthenticatedHandler is invoked when the user supplied token is valid and has not expired.
//the token contains the username, token expiry date and user plan.
//The received post data is passed to data.ProcessRequest.
func myAuthenticatedHandler(ctx iris.Context) {
	user := ctx.Values().Get("jwt").(*jwt.Token)
	foobar := user.Claims.(jwt.MapClaims)

	var username, plan string

	for key, value := range foobar {
		if key == "username" {
			username = value.(string)
		}
		if key == "plan" {
			plan = value.(string)
		}
		//ctx.Writef("%s = %s\n", key, value)
	}

	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		logger.Log.Sugar().Errorf("failed to read data from post body with error %e", err)
	}
	response, err := data.ProcessRequest(body, username, plan)
	var rejectReason string
	if err != nil {
		logger.Log.Sugar().Errorf("failed to write process request with error %e", err)
		rejectReason = fmt.Sprintf("failed to write process request with error %e", err)
		_, err = ctx.JSON(iris.Map{
			"status":  http.StatusUnauthorized,
			"message": rejectReason,
		})
		if err != nil {
			logger.Log.Sugar().Errorf("failed to return response with error %e", err)
			_, err = ctx.JSON(iris.Map{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
		}
	}else{
		if response.Type == "xml" {
			requestStr, _ := xml.Marshal(response)
			_, err = ctx.Writef(string(requestStr))
		}else{
			requestStr, _ := json.Marshal(response)
			_, err = ctx.Writef(string(requestStr))
		}
		if err != nil {
			logger.Log.Sugar().Errorf("failed to return response with error %e", err)
		}
	}

}

func main() {
	app := iris.New()

	j := jwt.New(jwt.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("demoappdemoappdemoapp"), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	env.LoadEnv()
	logger.Log.Info("starting service at port")
	defer data.CloseConnection()
	app.Get("/login", loginHandler)
	app.Post("/fetch", j.Serve, myAuthenticatedHandler)
	_ = app.Listen(":8080")
}
