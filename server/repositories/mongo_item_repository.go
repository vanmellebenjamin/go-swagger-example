package repositories

import (
	"context"
	"errors"
	"flightAPI/server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strings"
)

type MongoItemRepository struct {
	client         *mongo.Client
	dbName         string
	collectionName string
}

type itemMongo struct {
	ID          int    `bson:"_id"`
	Completed   bool   `bson:"completed"`
	Description string `bson:"description"`
}

func NewMongoItemRepository(client *mongo.Client, dbName string, collectionName string) *MongoItemRepository {
	repository := new(MongoItemRepository)
	repository.client = client
	repository.dbName = dbName
	repository.collectionName = collectionName
	return repository
}

func (mongoItemRepository *MongoItemRepository) AddItem(item models.Item) (*models.Item, error) {
	collection := mongoItemRepository.getCollection()

	mongoItem, err := toMongoItem(item)
	if err != nil {
		return nil, err
	}

	_, err = collection.InsertOne(context.TODO(), mongoItem)
	if err != nil {
		if strings.Contains(err.Error(), "E11000") {
			return nil, errors.New("already_exist")
		}
		panic(err)
		return nil, err
	}

	return &item, nil
}

func (mongoItemRepository *MongoItemRepository) FindItem(ID int32) (*models.Item, error) {
	collection := mongoItemRepository.getCollection()

	var itemMongo itemMongo
	err := collection.FindOne(context.TODO(), bson.M{"_id": ID}).Decode(&itemMongo)
	if err != nil {
		log.Fatal(err)
	}

	item, err := toItem(itemMongo)
	if err != nil {
		log.Fatal(err)
	}

	return item, nil
}

func (mongoItemRepository *MongoItemRepository) DeleteItem(ID int32) error {
	collection := mongoItemRepository.getCollection()
	deleteResult, err := collection.DeleteOne(context.TODO(), bson.M{"_id": ID})
	if err != nil {
		log.Fatal(err)
		return err
	}
	if deleteResult.DeletedCount == 0 {
		return errors.New("not_found")
	}
	return nil
}

func (mongoItemRepository *MongoItemRepository) FindItems(from int32, limit int32) ([]*models.Item, error) {
	collection := mongoItemRepository.getCollection()

	limit64 := int64(limit)
	skip := int64(from)
	findOptions := options.FindOptions{
		Limit: &limit64,
		Skip:  &skip,
		Sort:  bson.D{{"_id", 1}},
	}
	cursor, err := collection.Find(context.TODO(), bson.D{}, &findOptions)
	defer cursor.Close(nil)
	if err != nil {
		return nil, err
	}

	var items []*models.Item
	var itemMongo itemMongo
	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&itemMongo)
		if err != nil {
			return nil, err
		}
		item, err := toItem(itemMongo)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (mongoItemRepository *MongoItemRepository) UpdateItem(item models.Item) (*models.Item, error) {
	collection := mongoItemRepository.client.Database(mongoItemRepository.dbName).Collection(mongoItemRepository.collectionName)
	mongoItem, err := toMongoItem(item)
	if err != nil {
		return nil, err
	}
	replaceResult, err := collection.ReplaceOne(context.TODO(), bson.M{"_id": item.ID}, mongoItem)
	if err != nil {
		return nil, err
	}
	if replaceResult.MatchedCount == 0 {
		return nil, errors.New("not_found")
	}
	return &item, nil
}

func (mongoItemRepository *MongoItemRepository) getCollection() *mongo.Collection {
	collection := mongoItemRepository.client.Database(mongoItemRepository.dbName).Collection(mongoItemRepository.collectionName)
	return collection
}

func toMongoItem(item models.Item) (*itemMongo, error) {
	itemMongo := new(itemMongo)
	itemMongo.ID = int(item.ID)
	itemMongo.Description = *item.Description
	itemMongo.Completed = *item.Completed
	return itemMongo, nil
}

func toItem(itemMongo itemMongo) (*models.Item, error) {
	return &models.Item{
		ID:          int32(itemMongo.ID),
		Description: &itemMongo.Description,
		Completed:   &itemMongo.Completed,
	}, nil
}
