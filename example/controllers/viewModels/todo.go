package viewModels

import "errors"

type AddTodo struct {
	Text string `json:"text" example:"Go get groceries"`
}

func (todo AddTodo) Validation() error {
	if len(todo.Text) == 0 {
		return errors.New("text is empty")
	}

	return nil
}
