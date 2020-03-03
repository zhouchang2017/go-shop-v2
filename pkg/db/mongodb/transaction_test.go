package mongodb

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func stringify(doc interface{}, opts ...string) string {
	if len(opts) == 2 {
		b, _ := json.MarshalIndent(doc, opts[0], opts[1])
		return string(b)
	}
	b, _ := json.Marshal(doc)
	return string(b)
}

func TestTransactionCommit(t *testing.T) {

	replicaSetConfig := Config{
		Host:       "mongodb-primary",
		Port:       "30000",
		Database:   "go-shop-test",
		Username:   "root",
		Password:   "12345678",
		ReplicaSet: "rs0",
	}

	Connect(replicaSetConfig)
	defer Close()
	con := GetConFn()

	collectionExamples := "collection_examples"

	var err error

	var ctx = context.Background()
	var id = primitive.NewObjectID()
	var doc = bson.M{"_id": id, "hometown": "Atlanta", "year": int32(1998)}
	var result *mongo.UpdateResult
	var session mongo.Session
	var update = bson.D{{Key: "$set", Value: bson.D{{Key: "year", Value: int32(2000)}}}}

	collection := con.Collection(collectionExamples)
	if _, err = collection.InsertOne(ctx, doc); err != nil {
		t.Fatal(err)
	}

	if session, err = con.client.StartSession(); err != nil {
		t.Fatal(err)
	}
	if err = session.StartTransaction(); err != nil {
		t.Fatal(err)
	}
	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if result, err = collection.UpdateOne(sc, bson.M{"_id": id}, update); err != nil {
			t.Fatal(err)
		}
		if result.MatchedCount != 1 || result.ModifiedCount != 1 {
			t.Fatal("replace failed, expected 1 but got", result.MatchedCount)
		}

		if err = session.CommitTransaction(sc); err != nil {
			t.Fatal(err)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	session.EndSession(ctx)

	var v bson.M
	if err = collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&v); err != nil {
		t.Fatal(err)
	}
	if v["year"] != int32(2000) {
		t.Log(stringify(v))
		t.Fatal("expected 2000 but got", v["year"])
	}

	res, _ := collection.DeleteOne(ctx, bson.M{"_id": id})
	if res.DeletedCount != 1 {
		t.Fatal("delete failed, expected 1 but got", res.DeletedCount)
	}
}
