package database

import (
	"context"
	"fmt"
	"log/slog"

	"task/app/api"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, log *slog.Logger) map[string]*mongo.Collection {
	// Create client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Error(fmt.Sprintf("Connection Error: %d", err))
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Error(fmt.Sprintf("Connection Error: %d", err))
	}

	log.Info("Connected to MongoDB!")

	collections := map[string]*mongo.Collection{
		"Expressions": client.Database("calc").Collection("expressions"),
		"Operations":  client.Database("calc").Collection("operations"),
		"Agents":      client.Database("calc").Collection("agents"),
		"Tasks":       client.Database("calc").Collection("tasks"),
	}

	return collections
}

// =========================================================================================================================
// =========================================================================================================================

// Запись или обновление данных в базе данных
func SetOperations(ctx context.Context, collection *mongo.Collection, operations []api.OperationDuration) error {
	for _, op := range operations {
		filter := bson.M{"operationname": op.OperationName}
		update := bson.M{"$set": op}
		_, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
		if err != nil {
			return err
		}
	}
	return nil
}

// Получение всех операций из базы данных
func GetOperations(ctx context.Context, collection *mongo.Collection) ([]api.OperationDuration, error) {
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var retrievedOperations []api.OperationDuration
	for cursor.Next(ctx) {
		var op api.OperationDuration
		err := cursor.Decode(&op)
		if err != nil {
			return nil, err
		}
		retrievedOperations = append(retrievedOperations, op)
	}
	return retrievedOperations, nil
}

// =========================================================================================================================
// =========================================================================================================================

func Add(ctx context.Context, collection *mongo.Collection, expression api.Expression) error {
	filter := bson.M{"id": expression.ExpressionData}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		// Если элемент не найден, создаем новый
		update := bson.M{"$set": expression}
		_, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
		return err
	}
	return nil
}

func GetByID(ctx context.Context, collection *mongo.Collection, id string) (api.Expression, error) {
	var expression api.Expression
	filter := bson.M{"id": id}
	err := collection.FindOne(ctx, filter).Decode(&expression)
	return expression, err
}

func GetAll(ctx context.Context, collection *mongo.Collection) ([]api.Expression, error) {
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var expressions []api.Expression
	err = cursor.All(ctx, &expressions)
	return expressions, err
}

func Update(ctx context.Context, collection *mongo.Collection, id string, newExpression api.Expression) error {
	filter := bson.M{"id": id}
	_, err := collection.ReplaceOne(ctx, filter, newExpression)
	return err
}

// =========================================================================================================================
// =========================================================================================================================

func AddAgent(ctx context.Context, collection *mongo.Collection, agent api.Agent) error {
	filter := bson.M{"id": agent.ID}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		// Если элемент не найден, создаем новый
		update := bson.M{"$set": agent}
		_, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
		return err
	}
	return nil
}

func GetAllAgents(ctx context.Context, collection *mongo.Collection) ([]api.Agent, error) {
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var agents []api.Agent
	err = cursor.All(ctx, &agents)
	return agents, err
}

func UpdateAgent(ctx context.Context, collection *mongo.Collection, newAgent api.Agent) error {
	filter := bson.M{"id": newAgent.ID}
	_, err := collection.ReplaceOne(ctx, filter, newAgent)
	return err
}
