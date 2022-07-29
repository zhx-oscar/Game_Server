package DB

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/rand"
	"testing"
)

//func TestMain(m *testing.M) {
//	// replace mongo
//	testReplaceMongo()
//	testReplaceRedis()
//
//	m.Run()
//}

func testReplaceMongo() {
	clientInfo := options.Client()
	clientInfo.ApplyURI("mongodb://192.168.23.194:27017/")

	client, err := mongo.Connect(context.Background(), clientInfo)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		fmt.Println(err)
		return
	}

	MongoDB = client.Database("testdb")
	if MongoDB == nil {
		fmt.Println(err)
		return
	}
}

func testReplaceRedis() {
	//redisOptions := &redis.Options{}
	//RedisDB = redis.NewClient(redisOptions)
	//if _, err := RedisDB.Ping().Result(); err != nil {
	//	fmt.Println(err)
	//}
}

type TestStu struct {
	ID primitive.ObjectID `bson:"_id"`
	A  string             `bson:"a"`
	B  int64              `bson:"b"`
	C  bson.M             `bson:"c"`
}

func BenchmarkInsert(b *testing.B) {
	b.ResetTimer()

	t := TestStu{}

	for i := 0; i < b.N; i++ {
		t.ID = primitive.NewObjectID()
		t.A = primitive.NewObjectID().Hex()
		t.B = rand.Int63n(500000000)
		t.C = bson.M{
			primitive.NewObjectID().Hex(): primitive.NewObjectID().Hex(),
			primitive.NewObjectID().Hex(): primitive.NewObjectID().Hex(),
			primitive.NewObjectID().Hex(): primitive.NewObjectID().Hex(),
			primitive.NewObjectID().Hex(): primitive.NewObjectID().Hex(),
			primitive.NewObjectID().Hex(): primitive.NewObjectID().Hex(),
		}

		if _, err := MongoDB.Collection("test_insert").InsertOne(context.Background(), t); err != nil {
			panic(err)
		}
	}
}

func BenchmarkFind(b *testing.B) {
	cursor, err := MongoDB.Collection("test_insert").Find(context.Background(), bson.M{}, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.Background())

	ts := []TestStu{}

	if err := cursor.All(context.Background(), &ts); err != nil {
		panic(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		idx := rand.Int31n(int32(len(ts)))
		id := ts[idx].ID

		if rv := MongoDB.Collection("test_insert").FindOne(context.Background(), bson.M{"_id": id}); rv.Err() != nil {
			panic(rv.Err())
		} else {
			t := TestStu{}
			if err := rv.Decode(&t); err != nil {
				panic(err)
			}
		}
	}
}
