package main

import (
	"log"
	"strings"
	"text/template"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

type TodoItem struct {
	Content string
	Done bool
}

var todos []TodoItem = make([]TodoItem, 0)

func homeHandler(c *fiber.Ctx) error {
	return c.Status(200).Render("index", nil)
}

var todoElementsTemplate = template.New("todo-elements")

func addTodoHandler(c *fiber.Ctx) error {
	todo_content := strings.Clone(strings.TrimSpace(c.FormValue("todo")))
	if todo_content != "" {
		todos = append(todos, TodoItem{Content: todo_content, Done: false})
	}
	return c.Status(200).Render("todo_elements", fiber.Map{"Todos": todos})
}

func main() {
	todos = append(todos, TodoItem{Content: "Some task", Done: false})
	todos = append(todos, TodoItem{Content: "Some other task", Done: true})
	viewEngine := html.New("./views", ".html")
	template.Must(template.ParseGlob("views/*"))
	// log.Fatal(tpl.ExecuteTemplate(os.Stdout, "todo_elements.html", todos))
	app := fiber.New(fiber.Config{Views: viewEngine})
	app.Static("/static", "./static")
	app.Use(logger.New())
	app.Get("/", homeHandler)
	app.Post("/front/add-todo", addTodoHandler)
	log.Fatal(app.Listen(":8080"))
}
