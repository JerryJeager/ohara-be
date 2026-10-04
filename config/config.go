package config

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/genai"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Session *gorm.DB

func GetSession() *gorm.DB {
	return Session
}

var AI *genai.Client

func ConnectToDB() {
	var db *gorm.DB
	var err error

	//production DB config:::
	connectionString := os.Getenv("DATABASE_URL")
	db, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	Session = db.Session(&gorm.Session{SkipDefaultTransaction: true, PrepareStmt: false})
	if Session != nil {
		log.Print("success: created db session")
	}
}

func NewAIClient() {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	AI = client
}

func LoadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Print(err)
		log.Print("failed to load envirionment variables")
		// log.Fatal("failed to load environment variables")
	}
}
