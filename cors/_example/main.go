package main

import (
	"github.com/teamlint/iris"

	"github.com/teamlint/middleware/cors"
)

func main() {
	app := iris.New()

	opts := cors.Options{
		AllowedOrigins: []string{"*"}, // allows everything, use that to change the hosts.
		AllowedHeaders: []string{"*"},
		// AllowedMethods:   []string{"GET", "POST", "PUT", "HEAD"},
		AllowCredentials: true,
		Debug:            true,
	}

	crs := cors.New(opts)
	// app.Use(cors.AllowAll())
	app.Use(crs)
	app.AllowMethods(iris.MethodOptions)

	// v1 := app.Party("/api/v1", crs).AllowMethods(iris.MethodOptions) // <- important for the preflight.
	v1 := app.Party("/api/v1")
	{
		v1.Get("/home", func(ctx iris.Context) {
			ctx.WriteString("Hello from /home")
		})
		v1.Get("/about", func(ctx iris.Context) {
			ctx.WriteString("Hello from /about")
		})
		v1.Post("/send", func(ctx iris.Context) {
			ctx.WriteString("sent")
		})
		v1.Put("/send", func(ctx iris.Context) {
			ctx.WriteString("updated")
		})
		v1.Delete("/send", func(ctx iris.Context) {
			ctx.WriteString("deleted")
		})
	}

	app.Run(iris.Addr(":8082"))
}
