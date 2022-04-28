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

package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/suecodelabs/cnfuzz/example/controllers/viewModels"
	"github.com/suecodelabs/cnfuzz/example/model"
	"github.com/suecodelabs/cnfuzz/example/util"
)

func AddTodoController(routes *gin.RouterGroup) {
	var todoGroup = routes.Group("/todo")
	todoGroup.GET(":id", getTodo)
	todoGroup.GET("", listTodos)
	todoGroup.GET("search", searchTodo)
	todoGroup.POST("", createTodo)
}

// getTodo godoc
// @Summary Get a Todo item
// @Description Get a single Todo item
// @Security ApiKeyAuth
// @Tags todos
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} model.Todo
// @Failure 400 {object} util.HttpError
// @Failure 404 {object} util.HttpError
// @Router /api/todo/{id} [GET]
func getTodo(ctx *gin.Context) {
	var id = ctx.Params.ByName("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		util.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	todo, err := model.FindTodo(aid)
	if err != nil {
		util.NewError(ctx, http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, todo)
}

// listTodo godoc
// @Summary List Todo items
// @Description Get all the Todo items
// @Security ApiKeyAuth
// @Tags todos
// @Produce json
// @Param user header string true "Username or something"
// @Success 200 {array} model.Todo
// @Router /api/todo [GET]
func listTodos(ctx *gin.Context) {
	uname := ctx.GetHeader("user")
	log.Printf("User %s wants a list of all todos", uname)
	todos := model.GetTodos()
	ctx.JSON(http.StatusOK, todos)
}

// searchTodo godoc
// @Summary Find a Todo item
// @Description Find Todo items that contain a substring
// @Security ApiKeyAuth
// @Tags todos
// @Produce json
// @Param substring query string true "Substring of a todo item"
// @Success 200 {array} model.Todo
// @Failure 400 {object} util.HttpError
// @Failure 404 {object} util.HttpError
// @Router /api/todo/search [GET]
func searchTodo(ctx *gin.Context) {
	substring := ctx.Query("substring")

	todo, err := model.SearchTodos(substring)
	if err != nil {
		util.NewError(ctx, http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, todo)
}

// createTodo godoc
// @Summary Create a Todo item
// @Description Add a new todo item to the list
// @Security ApiKeyAuth
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body viewModels.AddTodo true "Create todo"
// @Success 200 {object} model.Todo
// @Failure 400 {object} util.HttpError
// @Failure 404 {object} util.HttpError
// @Router /api/todo [POST]
func createTodo(ctx *gin.Context) {
	var todoVm viewModels.AddTodo
	if err := ctx.ShouldBindJSON(&todoVm); err != nil {
		util.NewError(ctx, http.StatusBadRequest, err)
		return
	}
	if err := todoVm.Validation(); err != nil {
		util.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	todo := model.Todo{Text: todoVm.Text}
	todoId := model.InsertTodo(todo)
	todo.Id = todoId
	ctx.JSON(http.StatusOK, todo)
}
