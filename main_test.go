package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_TodoPost(t *testing.T) {
	Convey("Given to-do DTO", t, func() {
		repository := GetTestRepository()
		service := NewService(repository)
		api := NewAPI(service)

		todo := TodoDTO{
			Content: "To-do post request olustur.",
			Done:    false,
		}

		Convey("When create todo post request", func() {
			todoByte, err := json.Marshal(todo)
			So(err, ShouldBeNil)

			todoReader := bytes.NewReader(todoByte)
			request, err := http.NewRequest(http.MethodPost, "/todo", todoReader)
			So(err, ShouldBeNil)

			request.Header.Add("Content-Type", "application/json")
			request.Header.Set("Content-Length", strconv.Itoa(len(todoByte)))

			app := ServiceSetup(api)
			response, err := app.Test(request, 20000)
			So(err, ShouldBeNil)

			Convey("Then status code should be 201", func() {
				So(response.StatusCode, ShouldEqual, fiber.StatusCreated)

				Convey("Then to-do Should be returned", func() {
					responseBody, err := ioutil.ReadAll(response.Body)
					So(err, ShouldBeNil)

					returnedData := TodoDTO{}

					err = json.Unmarshal(responseBody, &returnedData)
					So(err, ShouldBeNil)

					So(returnedData.ID, ShouldNotBeEmpty)
					So(returnedData.Content, ShouldEqual, todo.Content)
					So(returnedData.Done, ShouldEqual, todo.Done)
					repository.DeleteTodoRepository(returnedData.ID)
				})
			})
		})
	})
}

func Test_TodoGet(t *testing.T) {
	Convey("Given to-do model in database", t, func() {
		repository := GetTestRepository()
		service := NewService(repository)
		api := NewAPI(service)

		todoID := uuid.New().String()
		todoModel := TodoModel{
			ID:        todoID,
			Content:   "To-do get request olustur.",
			Done:      false,
			CratedAt:  time.Now().Round(time.Minute).UTC(),
			UpdatedAt: time.Now().Round(time.Minute).UTC(),
		}

		repository.AddTodoRepository(&todoModel)

		Convey("When I get request", func() {
			request, _ := http.NewRequest(http.MethodGet, fmt.Sprint("/todo/", todoID), nil)

			request.Header.Add("Content-Type", "application/json")

			app := ServiceSetup(api)
			response, err := app.Test(request, 30000)
			So(err, ShouldBeNil)

			Convey("Then Status Code Should be 200", func() {
				So(response.StatusCode, ShouldEqual, fiber.StatusOK)

				Convey("Then to-do Should be returned", func() {
					responseBody, err := ioutil.ReadAll(response.Body)
					So(err, ShouldBeNil)

					returnedData := TodoDTO{}

					err = json.Unmarshal(responseBody, &returnedData)
					So(err, ShouldBeNil)

					So(returnedData.ID, ShouldNotBeEmpty)
					So(returnedData.ID, ShouldEqual, todoModel.ID)
					So(returnedData.Content, ShouldEqual, todoModel.Content)
					So(returnedData.Done, ShouldEqual, todoModel.Done)
				})
			})
		})
		repository.DeleteTodoRepository(todoID)
	})
}

func Test_TodoListGet(t *testing.T) {
	Convey("Given to-do model in database", t, func() {
		repository := GetTestRepository()
		service := NewService(repository)
		api := NewAPI(service)

		todoID1 := uuid.New().String()
		todoID2 := uuid.New().String()

		todoModel1 := TodoModel{
			ID:        todoID1,
			Content:   "To-do get list request olustur.",
			Done:      false,
			Index:     0,
			CratedAt:  time.Now().Round(time.Minute).UTC(),
			UpdatedAt: time.Now().Round(time.Minute).UTC(),
		}
		todoModel2 := TodoModel{
			ID:        todoID2,
			Content:   "To-do post list request olustur.",
			Done:      false,
			Index:     1,
			CratedAt:  time.Now().Round(time.Minute).UTC(),
			UpdatedAt: time.Now().Round(time.Minute).UTC(),
		}

		repository.AddTodoRepository(&todoModel1)
		repository.AddTodoRepository(&todoModel2)

		Convey("When I get request", func() {
			request, _ := http.NewRequest(http.MethodGet, "/todo", nil)

			request.Header.Add("Content-Type", "application/json")

			app := ServiceSetup(api)
			response, err := app.Test(request, 20000)
			So(err, ShouldBeNil)

			Convey("Then Status Code Should be 200", func() {
				So(response.StatusCode, ShouldEqual, fiber.StatusOK)

				Convey("Then to-do Should be returned", func() {
					responseBody, err := ioutil.ReadAll(response.Body)
					So(err, ShouldBeNil)

					returnedData := TodoListDTO{}

					err = json.Unmarshal(responseBody, &returnedData)
					So(err, ShouldBeNil)
					So(len(returnedData.TodoList), ShouldEqual, 2)
					So(returnedData.TodoList[0].ID, ShouldNotBeEmpty)
					So(returnedData.TodoList[1].ID, ShouldNotBeEmpty)
					So(returnedData.TodoList[0].ID, ShouldEqual, todoID2)
					So(returnedData.TodoList[1].ID, ShouldEqual, todoID1)
					So(returnedData.TodoList[0].Content, ShouldEqual, todoModel2.Content)
					So(returnedData.TodoList[1].Content, ShouldEqual, todoModel1.Content)
					So(returnedData.TodoList[0].Done, ShouldEqual, todoModel2.Done)
					So(returnedData.TodoList[1].Done, ShouldEqual, todoModel1.Done)
					So(returnedData.TodoList[0].Index, ShouldEqual, todoModel2.Index)
					So(returnedData.TodoList[1].Index, ShouldEqual, todoModel1.Index)
				})
			})
		})
		repository.DeleteTodoRepository(todoID1)
		repository.DeleteTodoRepository(todoID2)
	})
}

func Test_TodoUpdate(t *testing.T) {
	Convey("Given to-do model in database", t, func() {
		repository := GetTestRepository()
		service := NewService(repository)
		api := NewAPI(service)

		todoID := uuid.New().String()
		todoModel := TodoModel{
			ID:        todoID,
			Content:   "To-do put request olustur.",
			Done:      false,
			CratedAt:  time.Now().Round(time.Minute).UTC(),
			UpdatedAt: time.Now().Round(time.Minute).UTC(),
		}

		todoDTO := TodoDTO{
			Content: "To-do put request olustur. Update edildi",
		}

		repository.AddTodoRepository(&todoModel)
		Convey("When I get request", func() {
			todoByte, _ := json.Marshal(todoDTO)
			todoReader := bytes.NewReader(todoByte)
			request, _ := http.NewRequest(http.MethodPut, fmt.Sprint("/todo/", todoID), todoReader)

			request.Header.Add("Content-Type", "application/json")

			app := ServiceSetup(api)
			response, err := app.Test(request, 30000)
			So(err, ShouldBeNil)

			Convey("Then Status Code Should be 200", func() {
				So(response.StatusCode, ShouldEqual, fiber.StatusOK)

				Convey("Then to-do Should be returned", func() {
					responseBody, err := ioutil.ReadAll(response.Body)
					So(err, ShouldBeNil)

					returnedData := TodoDTO{}

					err = json.Unmarshal(responseBody, &returnedData)
					So(err, ShouldBeNil)

					So(returnedData.ID, ShouldNotBeEmpty)
					So(returnedData.ID, ShouldEqual, todoModel.ID)
					So(returnedData.Content, ShouldEqual, todoDTO.Content)
					So(returnedData.Done, ShouldEqual, todoDTO.Done)
				})
			})
		})
		repository.DeleteTodoRepository(todoID)
	})
}

func Test_TodoSortUpdate(t *testing.T) {
	Convey("Given to-do models in database", t, func() {
		repository := GetTestRepository()
		service := NewService(repository)
		api := NewAPI(service)

		todoID1 := uuid.New().String()
		todoModel1 := TodoModel{
			ID:        todoID1,
			Content:   "To-do put request olustur.",
			Done:      false,
			Index:     0,
			CratedAt:  time.Now().Round(time.Minute).UTC(),
			UpdatedAt: time.Now().Round(time.Minute).UTC(),
		}

		todoID2 := uuid.New().String()
		todoModel2 := TodoModel{
			ID:        todoID2,
			Content:   "To-do put request olustur.",
			Done:      false,
			Index:     1,
			CratedAt:  time.Now().Round(time.Minute).UTC(),
			UpdatedAt: time.Now().Round(time.Minute).UTC(),
		}

		todoID3 := uuid.New().String()
		todoModel3 := TodoModel{
			ID:        todoID3,
			Content:   "To-do put request olustur.",
			Done:      false,
			Index:     2,
			CratedAt:  time.Now().Round(time.Minute).UTC(),
			UpdatedAt: time.Now().Round(time.Minute).UTC(),
		}

		repository.AddTodoRepository(&todoModel1)
		repository.AddTodoRepository(&todoModel2)
		repository.AddTodoRepository(&todoModel3)
		Convey("When I get request", func() {
			request, _ := http.NewRequest(http.MethodPut, fmt.Sprint("/sort/", todoID1, "/", todoID3, "/", todoID2), nil)

			request.Header.Add("Content-Type", "application/json")

			app := ServiceSetup(api)
			response, err := app.Test(request, 30000)
			So(err, ShouldBeNil)

			Convey("Then Status Code Should be 200", func() {
				So(response.StatusCode, ShouldEqual, fiber.StatusOK)

				returnedData, _ := repository.GetTodoListRepository()

				Convey("Then to-do Should be returned", func() {
					So(len(returnedData.TodoList), ShouldEqual, 3)
					So(returnedData.TodoList[0].Index, ShouldEqual, todoModel2.Index)
					So(returnedData.TodoList[1].Index, ShouldEqual, ((float64(todoModel2.Index) + float64(todoModel1.Index)) / 2))
					So(returnedData.TodoList[2].Index, ShouldEqual, todoModel1.Index)
				})
			})
		})
		repository.DeleteTodoRepository(todoID1)
		repository.DeleteTodoRepository(todoID2)
		repository.DeleteTodoRepository(todoID3)
	})
}

func Test_TodoDelete(t *testing.T) {
	Convey("Given to-do model in database", t, func() {
		repository := GetTestRepository()
		service := NewService(repository)
		api := NewAPI(service)

		todoID := uuid.New().String()
		todoModel := TodoModel{
			ID:        todoID,
			Content:   "To-do delete request olustur.",
			Done:      false,
			CratedAt:  time.Now().Round(time.Minute).UTC(),
			UpdatedAt: time.Now().Round(time.Minute).UTC(),
		}

		repository.AddTodoRepository(&todoModel)

		Convey("When I delete request", func() {
			request, _ := http.NewRequest(http.MethodDelete, fmt.Sprint("/todo/", todoID), nil)

			request.Header.Add("Content-Type", "application/json")

			app := ServiceSetup(api)
			response, err := app.Test(request, 30000)
			So(err, ShouldBeNil)

			Convey("Then Status Code Should be 200", func() {
				So(response.StatusCode, ShouldEqual, fiber.StatusOK)
			})
		})
	})
}

func GetTestRepository() *Repository {

	config := ServiceConfig{
		Port:       ":8080",
		MongoDBURL: "mongodb://localhost:27017",
	}
	repository := NewRepository(config.MongoDBURL)

	return repository
}
