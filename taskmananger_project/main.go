package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

var db *sql.DB

func main() {
	db, _ = sql.Open("sqlite", "tasks.db")
	createTable()

	r := gin.Default()
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	api := r.Group("/api")
	{
		api.GET("/tasks", getTasks)
		api.POST("/tasks", createTask)
		api.PUT("/tasks/:id/status", updateStatus)
		api.DELETE("/tasks/:id", deleteTask)
	}

	r.Run(":8080")
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		status TEXT,
		created_at DATETIME
	);`
	db.Exec(query)
}

func getTasks(c *gin.Context) {
	rows, _ := db.Query("SELECT id, title, status, created_at FROM tasks ORDER BY created_at DESC")
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		rows.Scan(&t.ID, &t.Title, &t.Status, &t.CreatedAt)
		tasks = append(tasks, t)
	}
	c.JSON(http.StatusOK, tasks)
}

func createTask(c *gin.Context) {
	var input struct {
		Title string `json:"title"`
	}
	if err := c.BindJSON(&input); err != nil || input.Title == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	db.Exec(
		"INSERT INTO tasks(title, status, created_at) VALUES (?, ?, ?)",
		input.Title, "todo", time.Now(),
	)
	c.Status(http.StatusCreated)
}

func updateStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var input struct {
		Status string `json:"status"`
	}
	c.BindJSON(&input)
	db.Exec("UPDATE tasks SET status = ? WHERE id = ?", input.Status, id)
	c.Status(http.StatusOK)
}

func deleteTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	db.Exec("DELETE FROM tasks WHERE id = ?", id)
	c.Status(http.StatusOK)
}
