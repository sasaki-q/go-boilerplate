package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	sm "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewProduction()
)

type DBValue struct {
	Password string `json:"password" validate:"required"`
	DBName   string `json:"dbname" validate:"required"`
	Host     string `json:"host" validate:"required"`
	User     string `json:"user" validate:"required"`
}

func main() {
	if os.Getenv("ENV") == "DEV" {
		if err := godotenv.Load(".env"); err != nil {
			panic(fmt.Sprintf("err: cannot load envfile. / %v", err))
		}
	}

	ctx := context.Background()
	conn := getDBConn(ctx)

	m, err := migrate.New(
		"file://migrations",
		fmt.Sprintf(
			"postgres://%s:%s@%s:5432/%s?sslmode=disable",
			conn.User, conn.Password, conn.Host, conn.DBName),
	)

	if err != nil {
		panic(err)
	}

	if os.Getenv("TYPE") == "UP" {
		if err := m.Up(); err != nil {
			if e := deleteDirtyRecord(conn); e != nil {
				logger.Sugar().Errorf("err: cannot update dirty record. / %v", e)
			}
			panic(fmt.Sprintf("err: cannot up migration. / %v", err))
		}
		logger.Info("success: up migration")
	} else {
		if err := m.Down(); err != nil {
			if e := deleteDirtyRecord(conn); e != nil {
				logger.Sugar().Errorf("err: cannot update dirty record. / %v", e)
			}
			panic(fmt.Sprintf("err: cannot down migration. / %v", err))
		}
		logger.Info("success: down migration")
	}
}

func deleteDirtyRecord(e *DBValue) error {
	db, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
			e.Host, e.User, e.Password, e.DBName),
	)
	defer db.Close()

	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM schema_migrations WHERE dirty = True")
	return err
}

func getDBConn(ctx context.Context) *DBValue {
	secretID := os.Getenv("SECRET_ID")
	if secretID == "" {
		panic("err: please set secret id.")
	}
	cfg, err := newConfig(ctx)
	if err != nil {
		panic(fmt.Sprintf("err cannot load config. / %v", err))
	}

	v, err := sm.NewFromConfig(cfg).GetSecretValue(ctx, &sm.GetSecretValueInput{SecretId: &secretID})
	if err != nil {
		panic(fmt.Sprintf("err cannot get secret value. / %v", err))
	}

	var e DBValue
	if err := json.Unmarshal([]byte(*v.SecretString), &e); err != nil {
		panic(fmt.Sprintf("cannot marshal value / %v", err))
	}

	return &e
}

func newConfig(ctx context.Context) (aws.Config, error) {
	if os.Getenv("ENV") == "DEV" {
		var (
			key    = os.Getenv("AWS_ACCESS_KEY")
			secret = os.Getenv("AWS_SECRET_ACCESS_KEY")
		)
		if key == "" || secret == "" {
			panic("err: please set key & secret.")
		}
		p := credentials.NewStaticCredentialsProvider(key, secret, "")

		return config.LoadDefaultConfig(
			ctx,
			config.WithRegion("ap-northeast-1"),
			config.WithCredentialsProvider(p),
		)
	}

	return config.LoadDefaultConfig(ctx, config.WithRegion("ap-northeast-1"))
}
