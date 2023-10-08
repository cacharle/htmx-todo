package main

import (
	"embed"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

type TodoItem struct {
	Content string
	Done    bool
}

var todos []TodoItem = make([]TodoItem, 0)

func homeHandler(c *fiber.Ctx) error {
	return c.Status(200).Render("views/index", nil)
}

var todoElementsTemplate = template.New("todo-elements")

func listTodoHandler(c *fiber.Ctx) error {
	return c.Status(200).Render("views/todos/list", fiber.Map{"Todos": todos})
}

func addTodoHandler(c *fiber.Ctx) error {
	todoContent := strings.Clone(strings.TrimSpace(c.FormValue("content")))
	if todoContent != "" {
		todos = append(todos, TodoItem{Content: todoContent, Done: false})
	}
	return listTodoHandler(c)
}

func parseTodoIndex(indexString string) (int, error) {
	i, err := strconv.Atoi(indexString)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	if i > len(todos) {
		return 0, errors.New("index out of bound")
	}
	return i, nil
}

func deleteTodoHandler(c *fiber.Ctx) error {
	i, err := parseTodoIndex(c.Params("index"))
	if err != nil {
		log.Println(err)
		return listTodoHandler(c)
	}
	todos = append(todos[:i], todos[:i+1]...)
	return listTodoHandler(c)
}

func patchTodoHandler(c *fiber.Ctx) error {
	i, err := parseTodoIndex(c.Params("index"))
	if err != nil {
		log.Println(err)
		return listTodoHandler(c)
	}
	todoContent := strings.Clone(strings.TrimSpace(c.FormValue("content")))
	todoDone := c.FormValue("done", "off")
	if todoContent != "" {
		todos[i].Content = todoContent
	}
	if todoDone == "on" {
		todos[i].Done = true
	} else {
		todos[i].Done = false
	}
	return listTodoHandler(c)
}

func todoEditHandler(c *fiber.Ctx) error {
	i, err := parseTodoIndex(c.Params("index"))
	if err != nil {
		log.Println(err)
		return c.SendStatus(500)
	}
	return c.Render("views/todos/edit", fiber.Map{"Index": i, "Todo": todos[i]})
}

//go:embed views/*
var viewsFS embed.FS

func main() {
	todos = append(todos, TodoItem{Content: "Some task", Done: false})
	todos = append(todos, TodoItem{Content: "Some other task", Done: true})
	viewEngine := html.NewFileSystem(http.FS(viewsFS), ".html")
	app := fiber.New(fiber.Config{Views: viewEngine})
	app.Static("/", "./static")
	app.Use(logger.New())
	app.Get("/", homeHandler)
	app.Get("/front/todos", listTodoHandler)
	app.Post("/front/todos", addTodoHandler)
	app.Patch("/front/todos/:index", patchTodoHandler)
	app.Delete("/front/todos/:index", deleteTodoHandler)
	app.Get("/front/todos/edit/:index", todoEditHandler)
	log.Fatal(app.Listen(":8080"))
}
