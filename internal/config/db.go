package config

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // Postgres driver
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	PostgresDB *sql.DB
	MongoDB    *mongo.Database
	RedisClient *redis.Client
)

func InitDB() {
	initPostgres()
	initMongo()
	initRedis()
}

func initPostgres() {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=academ_aide port=5432 sslmode=disable"
	}
	var err error
	PostgresDB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Failed to connect to Postgres:", err)
	}
	if err = PostgresDB.Ping(); err != nil {
		log.Fatal("Failed to ping Postgres:", err)
	}
	log.Println("Connected to PostgreSQL")
}

func initMongo() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}
	MongoDB = client.Database("academ_aide")
	log.Println("Connected to MongoDB")
}

func initRedis() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr: addr,
		Password: "",
		DB: 0,
	})
	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Failed to ping Redis:", err)
	}
	log.Println("Connected to Redis")
}
