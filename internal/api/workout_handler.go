package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"go-server/internal/store"
	"go-server/internal/utils"
	"go-server/middleware"
	"log"
	"net/http"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
		logger:       logger,
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutById(resWriter http.ResponseWriter, request *http.Request) {
	workoutId, err := utils.ReadID(request)
	if err != nil {
		wh.logger.Printf("Error: while reading param %v", err)
		utils.WriterJSON(resWriter, http.StatusNotFound, utils.Envelope{"error": "invalid workout id"})
		return
	}

	workout, err := wh.workoutStore.GetWorkoutByID(workoutId)
	if err != nil {
		wh.logger.Printf("Error: while executing GetWorkoutByID %v", err)
		utils.WriterJSON(resWriter, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriterJSON(resWriter, http.StatusOK, utils.Envelope{"workout": workout})
}

func (wh *WorkoutHandler) HandleCreateWorkout(resWriter http.ResponseWriter, request *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(request.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("Error: while decoding request body %v", err)
		utils.WriterJSON(resWriter, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	currentUser := middleware.GetUser(request)
	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriterJSON(resWriter, http.StatusBadRequest, utils.Envelope{"error": "cannot find user related to workout"})
		return
	}

	workout.UserID = currentUser.ID

	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		wh.logger.Printf("Error: while executing CreateWorkout %v", err)
		utils.WriterJSON(resWriter, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriterJSON(resWriter, http.StatusOK, utils.Envelope{"workout": createdWorkout})
}

func (wh *WorkoutHandler) HandleUpdateWorkout(resWriter http.ResponseWriter, request *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(request.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("Error: while decoding request body %v", err)
		utils.WriterJSON(resWriter, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	currentUser := middleware.GetUser(request)
	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriterJSON(resWriter, http.StatusBadRequest, utils.Envelope{"error": "cannot find user related to workout"})
		return
	}

	workoutOwner, err := wh.workoutStore.GetWorkoutOwner(int64(workout.ID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriterJSON(resWriter, http.StatusNotFound, utils.Envelope{"error": "workout not exist"})
			return
		}
		utils.WriterJSON(resWriter, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if workoutOwner != currentUser.ID {
		utils.WriterJSON(resWriter, http.StatusForbidden, utils.Envelope{"error": "not authorized to perform this action"})
		return
	}

	workout.UserID = currentUser.ID

	err = wh.workoutStore.UpdateWorkout(&workout)
	if err != nil {
		wh.logger.Printf("Error: while executing UpdateWorkout %v", err)
		utils.WriterJSON(resWriter, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriterJSON(resWriter, http.StatusOK, utils.Envelope{"workout": workout})
}

func (wh *WorkoutHandler) HandleDelete(resWriter http.ResponseWriter, request *http.Request) {
	workoutId, err := utils.ReadID(request)
	if err != nil {
		wh.logger.Printf("Error: while reading param %v", err)
		utils.WriterJSON(resWriter, http.StatusNotFound, utils.Envelope{"error": "invalid workout id"})
		return
	}

	currentUser := middleware.GetUser(request)
	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriterJSON(resWriter, http.StatusBadRequest, utils.Envelope{"error": "cannot find user related to workout"})
		return
	}

	workoutOwner, err := wh.workoutStore.GetWorkoutOwner(workoutId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriterJSON(resWriter, http.StatusNotFound, utils.Envelope{"error": "workout not exist"})
			return
		}
		utils.WriterJSON(resWriter, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if workoutOwner != currentUser.ID {
		utils.WriterJSON(resWriter, http.StatusForbidden, utils.Envelope{"error": "not authorized to perform this action"})
		return
	}

	err = wh.workoutStore.DeleteWorkout(workoutId)
	if err != nil {
		wh.logger.Printf("Error: while executing DeleteWorkout %v", err)
		utils.WriterJSON(resWriter, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriterJSON(resWriter, http.StatusOK, utils.Envelope{"removedElement": workoutId})
}
