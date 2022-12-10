package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databseName        = "logs"
	logsCollectionName = "logs"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Data      string    `json:"data" bson:"data"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database(databseName).Collection(logsCollectionName)

	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error inserting into logs:", err)
	}

	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	col := client.Database(databseName).Collection(logsCollectionName)

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := col.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("finding all docs error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry
	for cursor.Next(ctx) {
		var item LogEntry
		err := cursor.Decode(&item)
		if err != nil {
			log.Println("error decoding log into slice", err)
			return nil, err
		}

		logs = append(logs, &item)
	}

	return logs, nil
}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	col := client.Database(databseName).Collection(logsCollectionName)

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry LogEntry
	err = col.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	col := client.Database(databseName).Collection(logsCollectionName)

	if err := col.Drop(ctx); err != nil {
		return err
	}

	return nil
}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	col := client.Database(databseName).Collection(logsCollectionName)

	docID, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return nil, err
	}

	result, err := col.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: l.Name},
				{Key: "data", Value: l.Data},
				{Key: "updated_at", Value: time.Now()},
			}},
		},
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
