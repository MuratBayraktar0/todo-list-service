package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TodoEntity struct {
	ID        string    `bson:"_id"`
	Content   string    `bson:"content"`
	Done      bool      `bson:"done"`
	Index     float64   `bson:"index"`
	CratedAt  time.Time `bson:"createdat"`
	UpdatedAt time.Time `bson:"updatedat"`
}

type TodoListEntity struct {
	TodoList []TodoEntity `bson:"todolist"`
}

type Repository struct {
	client *mongo.Client
}

func NewRepository(dbUrl string) *Repository {
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	clientOptions := options.Client().ApplyURI(dbUrl)
	client, _ := mongo.Connect(ctx, clientOptions)
	return &Repository{client}
}

func (repository *Repository) AddTodoRepository(todoModel *TodoModel) (*TodoEntity, error) {
	collection := repository.client.Database("todo").Collection("todolist")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	todoEntity := ConvertTodoModeltoEntity(todoModel)
	_, err := collection.InsertOne(ctx, todoEntity)

	if err != nil {
		return nil, err
	}

	return repository.GetTodoRepository(todoEntity.ID)
}

func (repository *Repository) GetTodoRepository(id string) (*TodoEntity, error) {
	collection := repository.client.Database("todo").Collection("todolist")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	todoEntity := TodoEntity{}
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&todoEntity)

	if err != nil {
		return nil, err
	}
	return &todoEntity, nil
}

func (repository *Repository) GetTodoListRepository() (*TodoListEntity, error) {
	collection := repository.client.Database("todo").Collection("todolist")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"index", -1}})
	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	todoListEntity := TodoListEntity{}
	for cursor.Next(ctx) {
		todoEntity := TodoEntity{}
		cursor.Decode(&todoEntity)
		todoListEntity.TodoList = append(todoListEntity.TodoList, todoEntity)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return &todoListEntity, nil
}

func (repository *Repository) UpdateTodoRepository(id string, todoModel *TodoModel) (*TodoEntity, error) {
	collection := repository.client.Database("todo").Collection("todolist")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	todoEntity := ConvertTodoModeltoEntity(todoModel)

	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"_id", id}}
	update := bson.M{
		"$set": bson.M{
			"content":   todoEntity.Content,
			"done":      todoEntity.Done,
			"updatedat": todoEntity.UpdatedAt,
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update, opts)

	if err != nil {
		return nil, err
	}
	return repository.GetTodoRepository(id)
}

func (repository *Repository) UpdateTodoSortRepository(currentId string, newIndex float64) (*TodoEntity, error) {
	collection := repository.client.Database("todo").Collection("todolist")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"_id", currentId}}
	update := bson.M{
		"$set": bson.M{
			"index": newIndex,
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update, opts)

	if err != nil {
		return nil, err
	}
	return repository.GetTodoRepository(currentId)
}

func (repository *Repository) DeleteTodoRepository(id string) error {
	collection := repository.client.Database("todo").Collection("todolist")
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		return err
	}
	return nil
}

func ConvertTodoModeltoEntity(todoModel *TodoModel) *TodoEntity {
	todoEntity := TodoEntity{
		ID:        todoModel.ID,
		Content:   todoModel.Content,
		Done:      todoModel.Done,
		Index:     todoModel.Index,
		CratedAt:  todoModel.CratedAt,
		UpdatedAt: todoModel.UpdatedAt,
	}
	return &todoEntity
}
