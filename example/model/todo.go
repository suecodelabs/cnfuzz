// Copyright 2022 Sue B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
