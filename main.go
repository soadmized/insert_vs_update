package main

import (
	"context"
	"log"
)

const (
	documents = 1500000
	batchSize = 1000
)

func main() {
	ctx := context.Background()

	db := newDbConn(ctx)
	defer db.client.Disconnect(ctx)

	err := db.coll.Drop(ctx)
	if err != nil {
		log.Fatal(err)
	}
	db.insertManyDocs(ctx)

	err = db.coll.Drop(ctx)
	if err != nil {
		log.Fatal(err)
	}
	db.updateManyDocs(ctx)

}
