package routes

import (
	"go-server/internal/app"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	router := chi.NewRouter()

	router.Group(func(router chi.Router) {
		router.Use(app.UserMiddleware.Authenticate)

		router.With(app.UserMiddleware.RequireUser).Get("/workout/{id}", app.WorkoutHandler.HandleGetWorkoutById)

		router.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)
	
		router.Put("/workout", app.WorkoutHandler.HandleUpdateWorkout)
	
		router.Delete("/workout/{id}", app.WorkoutHandler.HandleDelete)
	})

	router.Get("/health", app.HealthCheck)
	router.Post("/user", app.UserHandler.HandleCreateUser)
	router.Post("/token/authentication", app.TokenHandler.HandleCreateToken)

	return router
}
