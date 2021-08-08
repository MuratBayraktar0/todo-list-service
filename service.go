package main

import (
	"time"

	"github.com/google/uuid"
)

//SERVICE
type TodoModel struct {
	ID        string    `Json:"id"`
	Content   string    `json:"content"`
	Done      bool      `json:"done"`
	CratedAt  time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
}

type TodoListModel struct {
	TodoList []TodoDTO `json:"todolist"`
}

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (service *Service) PostTodoService(todoDTO *TodoDTO) (*TodoDTO, error) {
	todoModel := ConvertTodoDTOtoModel(todoDTO)
	todoModel.ID = uuid.New().String()
	todoModel.CratedAt = time.Now().Round(time.Minute).UTC()
	todoModel.UpdatedAt = time.Now().Round(time.Minute).UTC()

	todoEntity, err := service.repository.AddTodoRepository(todoModel)
	if err != nil {
		return nil, err
	}

	return ConvertTodoEntitytoDTO(todoEntity), nil
}

func (service *Service) GetTodoService(id string) (*TodoDTO, error) {
	todoEntity, err := service.repository.GetTodoRepository(id)
	if err != nil {
		return nil, err
	}

	return ConvertTodoEntitytoDTO(todoEntity), nil
}

func (service *Service) GetTodoListService() (*TodoListDTO, error) {
	todoListEntity, err := service.repository.GetTodoListRepository()
	if err != nil {
		return nil, err
	}

	return ConvertTodoListEntitytoDTO(todoListEntity), nil
}

func (service *Service) UpdateTodoService(id string, todoDTO *TodoDTO) (*TodoDTO, error) {
	todoModel := ConvertTodoDTOtoModel(todoDTO)
	todoModel.UpdatedAt = time.Now().Round(time.Minute).UTC()
	todoEntity, err := service.repository.UpdateTodoRepository(id, todoModel)
	if err != nil {
		return nil, err
	}

	return ConvertTodoEntitytoDTO(todoEntity), nil
}

func (service *Service) DeleteTodoService(id string) error {
	err := service.repository.DeleteTodoRepository(id)
	if err != nil {
		return err
	}

	return nil
}

func ConvertTodoDTOtoModel(todoDTO *TodoDTO) *TodoModel {
	todoModel := TodoModel{
		ID:      todoDTO.ID,
		Content: todoDTO.Content,
		Done:    todoDTO.Done,
	}
	return &todoModel
}

func ConvertTodoEntitytoDTO(todoEntity *TodoEntity) *TodoDTO {
	todoDTO := TodoDTO{
		ID:      todoEntity.ID,
		Content: todoEntity.Content,
		Done:    todoEntity.Done,
	}
	return &todoDTO
}

func ConvertTodoListEntitytoDTO(todoListEntity *TodoListEntity) *TodoListDTO {
	todoListDTO := TodoListDTO{}
	for _, v := range todoListEntity.TodoList {
		todoListDTO.TodoList = append(todoListDTO.TodoList, *ConvertTodoEntitytoDTO(&v))
	}
	return &todoListDTO
}
