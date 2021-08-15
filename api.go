package main

import (
	"github.com/gofiber/fiber/v2"
)

//API
type TodoDTO struct {
	ID      string  `Json:"id"`
	Content string  `json:"content"`
	Done    bool    `json:"done"`
	Index   float64 `json:"index"`
}

type TodoListDTO struct {
	TodoList []TodoDTO `json:"todolist"`
}

type Api struct {
	service *Service
}

func NewAPI(service *Service) *Api {
	return &Api{
		service: service,
	}
}

func (api *Api) PostTodoApi(ctx *fiber.Ctx) error {
	todoDTO := TodoDTO{}
	ctx.BodyParser(&todoDTO)
	returnedData, err := api.service.PostTodoService(&todoDTO)

	switch err {
	case nil:
		ctx.Status(fiber.StatusCreated)
		ctx.JSON(returnedData)
		return nil
	case fiber.ErrConflict:
		ctx.Status(fiber.StatusConflict)
		return err
	case fiber.ErrBadRequest:
		ctx.Status(fiber.StatusBadRequest)
		return err
	default:
		ctx.Status(fiber.StatusInternalServerError)
		return err
	}
}

func (api *Api) GetTodoApi(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	returnedData, err := api.service.GetTodoService(id)

	switch err {
	case nil:
		ctx.Status(fiber.StatusOK)
		ctx.JSON(returnedData)
		return nil

	case fiber.ErrBadRequest:
		ctx.Status(fiber.StatusBadRequest)
		return err

	default:
		ctx.Status(fiber.StatusInternalServerError)
		return err
	}
}

func (api *Api) GetTodoListApi(ctx *fiber.Ctx) error {
	returnedData, err := api.service.GetTodoListService()

	switch err {
	case nil:
		ctx.Status(fiber.StatusOK)
		ctx.JSON(returnedData)
		return nil

	case fiber.ErrBadRequest:
		ctx.Status(fiber.StatusBadRequest)
		return err

	default:
		ctx.Status(fiber.StatusInternalServerError)
		return err
	}
}

func (api *Api) PutTodoApi(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	todoDTO := TodoDTO{}
	ctx.BodyParser(&todoDTO)
	returnedData, err := api.service.UpdateTodoService(id, &todoDTO)

	switch err {
	case nil:
		ctx.Status(fiber.StatusOK)
		ctx.JSON(returnedData)
		return nil

	case fiber.ErrBadRequest:
		ctx.Status(fiber.StatusBadRequest)
		return err

	default:
		ctx.Status(fiber.StatusInternalServerError)
		return err
	}
}

func (api *Api) PutSortApi(ctx *fiber.Ctx) error {
	currentId := ctx.Query("currentid")
	backId := ctx.Query("backid")
	frontId := ctx.Query("frontid")

	returnedData, err := api.service.UpdateTodoSortService(currentId, backId, frontId)

	switch err {
	case nil:
		ctx.Status(fiber.StatusOK)
		ctx.JSON(returnedData)
		return nil

	case fiber.ErrBadRequest:
		ctx.Status(fiber.StatusBadRequest)
		return err

	default:
		ctx.Status(fiber.StatusInternalServerError)
		return err
	}
}

func (api *Api) DeleteTodoApi(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	err := api.service.DeleteTodoService(id)

	switch err {
	case nil:
		ctx.Status(fiber.StatusOK)
		return nil

	case fiber.ErrBadRequest:
		ctx.Status(fiber.StatusBadRequest)
		return err

	default:
		ctx.Status(fiber.StatusInternalServerError)
		return err
	}
}
