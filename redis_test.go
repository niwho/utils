package common

import (
	"fmt"
	"os"
	"testing"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var redisClient *RedisClient
var mongoClient *MongoDb

func TestMain(m *testing.M) {
	var err error
	redisClient = NewRedisClient("10.116.76.211:6370", 15)
	mongoClient, err = InitMongo("miveshow", "miveshow123", "miveshow", []string{"10.116.76.211:27017"})
	//mongoClient, err = InitMongo("miveshow", "XE93G0uw7YSPHPeO", "miveshow", []string{"10.116.75.228:27017", "10.116.75.239:27017", "10.116.75.233:27017"})

	fmt.Println("InitMongo", mongoClient, err)
	os.Exit(m.Run())
}

func TestSetInt64(t *testing.T) {
	var val int64 = 11234
	err := redisClient.SetInt64("wtt_test", val, time.Hour)
	fmt.Printf("TestSetInt64 val=%v, err=%v\n", val, err)
}

func TestGetInt64(t *testing.T) {
	val, err := redisClient.GetInt64("key-with-expire-time")
	fmt.Printf("TestSetInt64 val=%v, err=%v\n", val, err)
}

func TestHGetInt(t *testing.T) {
	err := redisClient.HSet("hget-t1", "aa", 12)
	fmt.Printf("TestHSetInt val=%v, err=%v\n", "", err)
	val, err := redisClient.HGetInt64("hget-t1", "aa")
	fmt.Printf("TestHGetInt  val=%v, err=%v\n", val, err)
}

func TestZrange(t *testing.T) {
	vals, err := redisClient.ZRevRange("_1003_god0000@list", 0, -1)
	fmt.Println("TestZrange", vals, err)
}

func TestCreateIndexAndShard(t *testing.T) {
	//err := mongoClient.Insert("chat_messages", "")
	index := mgo.Index{
		Key:        []string{"session_id", "type", "-mid"},
		Unique:     true,
		DropDups:   true,
		Background: true, // See notes.
		Sparse:     true,
	}
	err := mongoClient.session.DB(mongoClient.db).C("chat_messages").EnsureIndex(index)
	fmt.Println("index", err)
	return
	var result interface{}

	err = mongoClient.session.DB(mongoClient.db).Run(bson.D{{"enableSharding", "mivechat"}}, &result)
	fmt.Println("enableshard~~", err, result)
	err = mongoClient.session.DB(mongoClient.db).Run(bson.D{{"shardCollection", "mivechat.chat_messages"}, {"key",
		bson.M{"session_id": "hashed"}}}, &result)
	fmt.Println("shard~~", err, result)
}

func TestCreateIndex2(t *testing.T) {
	//err := mongoClient.Insert("chat_messages", "")
	index := mgo.Index{
		Key:        []string{"type", "session_id", "uid"},
		Unique:     true,
		DropDups:   true,
		Background: true, // See notes.
		Sparse:     true,
	}
	err := mongoClient.session.DB(mongoClient.db).C("chat_state").EnsureIndex(index)
	fmt.Println("index222", err)
	return
}

func TestIp(t *testing.T) {
	ip := GetLocalIP()
	fmt.Println("GetLocalIP:", ip)
}
