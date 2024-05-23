package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"boilerplate/server/handlers"
	"boilerplate/server/repositories"
	"boilerplate/server/restapi"
	"boilerplate/server/services"

	"github.com/go-chi/chi/v5"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql/driver"
	"go.uber.org/zap"
)

var (
	location, _ = time.LoadLocation("Asia/Tokyo")
	logger, _   = zap.NewProduction()
)

type DBValue struct {
	Password string `json:"password" validate:"required"`
	DBName   string `json:"dbname" validate:"required"`
	Host     string `json:"host" validate:"required"`
	User     string `json:"user" validate:"required"`
}

func main() {
	ctx := context.TODO()
	s := services.NewAWSService(ctx, logger)
	v, err := s.GetSecretValue(ctx, os.Getenv("SECRET_ID"))
	if err != nil {
		logger.Sugar().Errorf("cannot get secret value / %v", err)
		panic(fmt.Sprintf("cannot get secret value / %v", err))
	}
	var e DBValue
	if err := json.Unmarshal([]byte(*v.SecretString), &e); err != nil {
		logger.Sugar().Errorf("cannot marshal value / %v", err)
		panic(fmt.Sprintf("cannot marshal value / %v", err))
	}
	db := initDB(e)
	r := chi.NewRouter()
	r.Get("/hc", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("healthy")) })
	restapi.HandlerFromMuxWithBaseURL(buildHandlers(db), r, "/v1")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), r); err != nil {
		logger.Sugar().Errorf("cannot start server / %v", err)
		panic(fmt.Sprintf("cannot start server / %v", err))
	}
}

func buildHandlers(db *sql.DB) *handlers.Handler {
	var (
		h  = handlers.NewHelper(logger)
		ur = repositories.NewUserRepository(db)
	)

	return &handlers.Handler{
		UserHandler: handlers.NewUserHandler(h, ur),
	}
}

func initDB(e DBValue) *sql.DB {
	db, err := sql.Open("postgres", driver.PSQLBuildQueryString(
		e.User,
		e.Password,
		e.DBName,
		e.Host,
		5432,
		"disable",
	))
	if err != nil {
		logger.Sugar().Errorf("cannot open database / %v", err)
		panic(fmt.Sprintf("cannot open database / %v", err))
	}
	db.SetConnMaxIdleTime(5)
	db.SetMaxOpenConns(5)
	db.SetConnMaxLifetime(300 * time.Second)
	boil.SetDB(db)
	boil.SetLocation(location)
	if os.Getenv("ENV") == "dev" {
		boil.DebugMode = true
	}
	return db
}
