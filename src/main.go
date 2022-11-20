package main

import (
    "context"
    "fmt"
    "net/http"
    "net/url"
    "os"
    "time"

    "github.com/go-redis/redis"
    "github.com/gorilla/mux"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

    "go.elastic.co/ecszap"
    "go.uber.org/zap"
)

var client *mongo.Client
var collection *mongo.Collection

type Tweet struct {
    ID       int64  `json:"_id,omitempty" bson:"_id,omitempty"`
    FullText string `json:"full_text,omitempty" bson:"full_text,omitempty"`
    User     struct {
        ScreenName string `json:"screen_name" bson:"screen_name"`
    } `json:"user,omitempty" bson:"user,omitempty"`
}

func GetTweetsEndpoint(response http.ResponseWriter, request *http.Request)    {}
func SearchTweetsEndpoint(response http.ResponseWriter, request *http.Request) {}

func main() {
    fmt.Println("Starting the application...")

    ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
    defer cancel()
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
    if err != nil {
        fmt.Println(err.Error())
    }
    defer func() {
        if err = client.Disconnect(ctx); err != nil {
            panic(err)
        }
    }()
    err = client.Database("synonyms").CreateCollection(context.Background(), "tweets")
    if err != nil {
        fmt.Println(err)
    }
    redisUrl, _ := url.Parse(os.Getenv("REDIS_URI"))
    redisPassword, _ := redisUrl.User.Password()
    redisDB := 0
    redisOptions := redis.Options{
        Addr:     redisUrl.Host,
        Password: redisPassword,
        DB:       redisDB,
    }
    clientRedis := redis.NewClient(&redisOptions)

    pong, err := clientRedis.Ping().Result()
    fmt.Println(pong, err)
    router := mux.NewRouter()
    router.HandleFunc("/tweets", GetTweetsEndpoint).Methods("GET")
    router.HandleFunc("/search", SearchTweetsEndpoint).Methods("GET")

    encoderConfig := ecszap.NewDefaultEncoderConfig()
    core := ecszap.NewCore(encoderConfig, os.Stdout, zap.DebugLevel)
    logger := zap.New(core, zap.AddCaller())
    logger = logger.With(zap.String("app", "myapp")).With(zap.String("environment", "psm"))

    http.ListenAndServe(":8000", router)
}
