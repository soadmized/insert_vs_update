package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type productDoc struct {
	ID   uuid.UUID `bson:"_id"`
	Name string    `bson:"name"`
}

type productDocs []productDoc

type dbConn struct {
	client *mongo.Client
	db     *mongo.Database
	coll   *mongo.Collection
}

func newDbConn(ctx context.Context) dbConn {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017").SetRegistry(BsonRegistry()))
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("super-testing")
	coll := db.Collection("insert-update-test")

	return dbConn{
		client: client,
		db:     db,
		coll:   coll,
	}

}

func (db *dbConn) insertManyDocs(ctx context.Context) {
	docs := prepareDocs()
	batch := make([]interface{}, 0, batchSize)

	start := time.Now()
	for _, d := range docs {
		batch = append(batch, d)

		if len(batch) >= batchSize {
			_, err := db.coll.InsertMany(ctx, batch)
			if err != nil {
				log.Fatal(errors.Wrap(err, "insert many"))
			}

			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		_, err := db.coll.InsertMany(ctx, batch)
		if err != nil {
			log.Fatal(errors.Wrap(err, "insert many out of for loop"))
		}
	}

	log.Print("INSERT MANY TIME: ", -start.Sub(time.Now()))
}

func (db *dbConn) updateManyDocs(ctx context.Context) {
	docs := prepareDocs()
	prDocs := make(productDocs, 0, len(docs))

	for _, v := range docs {
		c := v.(productDoc)

		prDocs = append(prDocs, c)
	}

	models := mongoWriteModels(prDocs)
	batch := make([]mongo.WriteModel, 0, batchSize)

	start := time.Now()
	for _, m := range models {
		batch = append(batch, m)

		if len(batch) >= batchSize {
			_, err := db.coll.BulkWrite(ctx, batch)
			if err != nil {
				log.Fatal(errors.Wrap(err, "update many"))
			}

			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		_, err := db.coll.BulkWrite(ctx, batch)
		if err != nil {
			log.Fatal(errors.Wrap(err, "update many out of for loop"))
		}
	}

	log.Print("BULK WRITE TIME: ", -start.Sub(time.Now()))
}

func prepareDocs() []interface{} {
	docs := make([]interface{}, 0, documents)

	for i := 0; i < documents; i++ {
		doc := productDoc{
			ID:   uuid.New(),
			Name: fmt.Sprintf("Name_%d", i),
		}

		docs = append(docs, doc)
	}

	return docs
}

func mongoWriteModels(docs productDocs) []mongo.WriteModel {
	models := make([]mongo.WriteModel, 0, len(docs))

	for _, v := range docs {
		filter := bson.M{"_id": v.ID}
		model := mongo.NewUpdateManyModel().
			SetFilter(filter).
			SetUpsert(true).
			SetUpdate(bson.M{"$set": bson.M{"_id": v.ID, "name": v.Name}})

		models = append(models, model)
	}

	return models
}
