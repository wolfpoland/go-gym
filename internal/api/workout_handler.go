package api

import (
	"encoding/json"
	"go-server/internal/store"
	"go-server/internal/utils"
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

	err = wh.workoutStore.DeleteWorkout(workoutId)
	if err != nil {
		wh.logger.Printf("Error: while executing DeleteWorkout %v", err)
		utils.WriterJSON(resWriter, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriterJSON(resWriter, http.StatusOK, utils.Envelope{"removedElement": workoutId})
}
