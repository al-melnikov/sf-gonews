package mongo

import (
	"GoNews/pkg/storage"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Storage Хранилище данных.
type Storage struct {
	Db *mongo.Client
}

const (
	databaseName   = "go_news" // имя учебной БД
	collectionName = "posts"   // имя коллекции в учебной БД
)

// New Конструктор, принимает строку подключения к БД.
func New(ctx context.Context, constr string) (*Storage, error) {
	mongoOpts := options.Client().ApplyURI(constr)
	client, err := mongo.Connect(ctx, mongoOpts)
	if err != nil {
		log.Fatal(err)
	}
	// не забываем закрывать ресурсы
	s := Storage{
		Db: client,
	}
	return &s, nil
}

func (mg *Storage) Posts() ([]storage.Post, error) {
	collection := mg.Db.Database(databaseName).Collection(collectionName)
	filter := bson.D{}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	var data []storage.Post
	for cur.Next(context.Background()) {
		var l storage.Post
		err := cur.Decode(&l)
		if err != nil {
			return nil, err
		}
		data = append(data, l)
	}
	return data, cur.Err()
}

func (mg *Storage) AddPost(p storage.Post) error {
	collection := mg.Db.Database(databaseName).Collection(collectionName)
	_, err := collection.InsertOne(context.Background(), p)
	if err != nil {
		return err
	}
	return nil
}

// пример запроса curl
/*
curl -X PUT http://localhost:8080/posts -H "Content-Type: application/json" -d '{"id":2,"title":"title 22","content":"new text","author_id":2,"name":"","created_at":12,"published_at":26}
*/
// Ищет такой же title, меняет все остальное
func (mg *Storage) UpdatePost(p storage.Post) error {
	collection := mg.Db.Database(databaseName).Collection(collectionName)
	filter := bson.D{{Key: "title", Value: p.Title}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "content", Value: p.Content}}}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

// удаляет с совпадающим title
func (mg *Storage) DeletePost(p storage.Post) error {
	collection := mg.Db.Database(databaseName).Collection(collectionName)
	filter := bson.D{{Key: "title", Value: p.Title}}
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}
