package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Fs02/grimoire/changeset"
	"github.com/Fs02/grimoire/params"
)

// TodoTable definse table name for stroing todos in database.
const TodoTable = "todos"

// Todo respresent a record stored in todos table.
type Todo struct {
	ID        uint
	Title     string
	Order     int
	Completed bool
}

// MarshalJSON implement custom marshaller to marshal url.
func (todo Todo) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Order     int    `json:"order"`
		Completed bool   `json:"completed"`
		URL       string `json:"url"`
	}{
		ID:        todo.ID,
		Title:     todo.Title,
		Order:     todo.Order,
		Completed: todo.Completed,
		URL:       fmt.Sprint(os.Getenv("URL"), todo.ID),
	})
}

// ChangeTodo prepares data before database operation.
func ChangeTodo(todo interface{}, params params.Params) *changeset.Changeset {
	ch := changeset.Cast(todo, params, []string{"title", "order", "completed"})
	changeset.ValidateRequired(ch, []string{"title"})
	changeset.ValidateRange(ch, "title", 1, 255)

	return ch
}

// CreateTodo is similar to ChangeTodo, except it also fills some default values and used before insert operation.
func CreateTodo(params params.Params) *changeset.Changeset {
	ch := ChangeTodo(Todo{}, params)
	changeset.PutChange(ch, "completed", false)

	return ch
}
