package main

import (
	"context"
	"log"
	"time"
)

const documents = 50000

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

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
