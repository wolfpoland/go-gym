package app

import (
	"database/sql"
	"fmt"
	"go-server/internal/api"
	"go-server/internal/store"
	"go-server/middleware"
	"go-server/migrations"
	"log"
	"net/http"
	"os"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	UserHandler    *api.UserHandler
	TokenHandler   *api.TokensHandler
	UserMiddleware middleware.UserMiddleware
	DB             *sql.DB
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FC, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	workoutStore := store.NewPostgresWorkoutStore(pgDB)
	userStore := store.NewPostgresUserStore(pgDB)
	tokenStore := store.NewPostgresTokenStore(pgDB)
	userMiddleware := middleware.UserMiddleware{UserStore: userStore}

	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)

	app := &Application{
		DB:             pgDB,
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		UserHandler:    userHandler,
		TokenHandler:   tokenHandler,
		UserMiddleware: userMiddleware,
	}

	return app, nil
}

func (a *Application) HealthCheck(resWriter http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(resWriter, "Status is available\n")
}
