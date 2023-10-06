package main

import (
	"log"
	"strings"
	"text/template"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var todos []string = make([]string, 0)

func homeHandler(c *fiber.Ctx) error {
	return c.Status(200).Render("views/index.html", fiber.Map{"Todos": todos})
}

var todoElementsTemplate = template.New("todo-elements")

func addTodoHandler(c *fiber.Ctx) error {
	todo := strings.Clone(strings.TrimSpace(c.FormValue("todo")))
	if todo != "" {
		todos = append(todos, todo)
	}
	return c.Status(200).Render("views/todo_elements.html", todos)
}

func main() {
	template.Must(template.ParseGlob("views/*"))
	app := fiber.New()
	app.Static("/static", "./static")
	app.Use(logger.New())
	app.Get("/", homeHandler)
	app.Post("/front/add-todo", addTodoHandler)
	log.Fatal(app.Listen(":8080"))
}
