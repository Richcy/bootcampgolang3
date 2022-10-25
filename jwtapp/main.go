package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

func main() {
	app := fiber.New()

	// Login route
	app.Post("/login", login)

	// Unauthenticated route
	app.Get("/", accessible)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("mysecretpassword"),
	}))

	// Restricted Routes
	app.Get("/restricted", restricted)
	app.Post("/upload", UploadFile)	

	app.Listen(":3000")
}

func login(c *fiber.Ctx) error {
	user := c.FormValue("user")
	pass := c.FormValue("pass")

	// Throws Unauthorized error
	if user != "john" || pass != "doe" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Create the Claims
	exp := time.Now().Add(time.Hour * 72)
	claims := jwt.MapClaims{
		"name":  "John Doe",
		"admin": true,
		"exp":   exp.Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("mysecretpassword"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"token": t,
		"expired": exp.Format("2006-01-02 15:04:05"),
	})
}

func accessible(c *fiber.Ctx) error {
	return c.SendString("Accessible")
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}


func UploadFile(c *fiber.Ctx) error {
	// Parse the multipart form:
	if form, err := c.MultipartForm(); err == nil {
		// => *multipart.Form
	
		// Get all files from "documents" key:
		files := form.File["file"]
		// => []*multipart.FileHeader
	
		// Loop through files:
		for _, file := range files {
		  fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
		  // => "tutorial.pdf" 360641 "application/pdf"
	
		  // Save the files to disk:
		  if err := c.SaveFile(file, fmt.Sprintf("./upload/%s", file.Filename)); err != nil {
			return err
		  }
		  return c.SendString("Succeed.. " + file.Filename)
		}
		return err
	  }
	  return c.SendStatus(400)
}