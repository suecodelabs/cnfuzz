package model

import (
	"errors"
	"strings"
)

var todos = []Todo{}

type Todo struct {
	Id   int    `json:"id" example:"1" format:"int32"`
	Text string `json:"text" example:"Go get groceries"`
}

func GetTodos() []Todo {
	return todos
}

func FindTodo(id int) (Todo, error) {
	for _, todo := range todos {
		if id == todo.Id {
			return todo, nil
		}
	}
	return Todo{}, errors.New("todo doesn't exist")
}

func SearchTodos(substring string) ([]Todo, error) {
	if len(substring) == 0 {
		return nil, errors.New("substring is empty")
	}

	var found []Todo
	for _, todo := range todos {
		if strings.Contains(strings.ToLower(todo.Text), strings.ToLower(substring)) {
			found = append(found, todo)
		}
	}
	return found, nil
}

func InsertTodo(todo Todo) int {
	lastId := 0
	if len(todos) > 0 {
		lastId = todos[len(todos)-1].Id
	}
	todo.Id = lastId + 1
	todos = append(todos, todo)
	return todo.Id
}
