package main

import (
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
)

func getTokenHandler(ctx iris.Context) {
	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, _ := token.SignedString([]byte("My Secret"))

	ctx.HTML(`Token: ` + tokenString + `<br/><br/>
    <a href="/secured?token=` + tokenString + `">/secured?token=` + tokenString + `</a>`)
}

func myAuthenticatedHandler(ctx iris.Context) {
	user := ctx.Values().Get("jwt").(*jwt.Token)

	ctx.Writef("This is an authenticated request\n")
	ctx.Writef("Claim content:\n")

	foobar := user.Claims.(jwt.MapClaims)
	for key, value := range foobar {
		ctx.Writef("%s = %s", key, value)
	}
}

func main(){
	app := iris.New()

	j := jwt.New(jwt.Config{
		// Extract by "token" url parameter.
		Extractor: jwt.FromParameter("token"),

		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("My Secret"), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	app.Get("/", getTokenHandler)
	app.Get("/secured", j.Serve, myAuthenticatedHandler)
	app.Listen(":8080")
}