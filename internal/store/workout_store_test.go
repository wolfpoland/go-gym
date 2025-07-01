package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres sslmode=disable port=5433")
	if err != nil {
		t.Fatal("opening test db %w", err)
	}

	err = Migrate(db, "../../migrations/")
	if err != nil {
		t.Fatal("cannot do migration %w", err)
	}

	_, err = db.Exec("TRUNCATE workouts, workout_entries CASCADE")
	if err != nil {
		t.Fatal("cannot TRUNCATE %w", err)
	}

	return db
}

func TestCreateWorkout(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewPostgresWorkoutStore(db)

	tests := []struct {
		name    string
		workout *Workout
		wantErr bool
	}{
		{
			name: "valid workout",
			workout: &Workout{
				Title:           "Push day",
				Description:     "upper body day",
				DurationSeconds: 3600,
				CaloriesBurned:  400,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Bench press",
						Sets:         3,
						Reps:         IntPtr(10),
						Weight:       FloatPtr(100),
						Notes:        "",
						OrderIndex:   0,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Workout with invalid entries",
			workout: &Workout{
				Title:           "Push day",
				Description:     "upper body day",
				DurationSeconds: 3600,
				CaloriesBurned:  400,
				Entries: []WorkoutEntry{
					{
						ExerciseName:    "Bench press",
						Sets:            3,
						Reps:            IntPtr(10),
						DurationSeconds: IntPtr(180),
						Weight:          FloatPtr(100),
						Notes:           "",
						OrderIndex:      0,
					},
					{
						ExerciseName:    "Plank",
						Sets:            3,
						Reps:            IntPtr(10),
						DurationSeconds: IntPtr(180),
						Weight:          FloatPtr(100),
						Notes:           "keep form",
						OrderIndex:      0,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdWorkout, err := store.CreateWorkout(tt.workout)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.workout.Title, createdWorkout.Title)
			assert.Equal(t, tt.workout.CaloriesBurned, createdWorkout.CaloriesBurned)
			assert.Equal(t, tt.workout.Description, createdWorkout.Description)

			retrieved, err := store.GetWorkoutByID(int64(createdWorkout.ID))
			require.NoError(t, err)

			assert.Equal(t, createdWorkout.ID, retrieved.ID)
			assert.Equal(t, createdWorkout.Title, retrieved.Title)
			assert.Equal(t, len(tt.workout.Entries), len(retrieved.Entries))

			for i := range retrieved.Entries {
				assert.Equal(t, tt.workout.Entries[i].ExerciseName, retrieved.Entries[i].ExerciseName)
				assert.Equal(t, tt.workout.Entries[i].Sets, retrieved.Entries[i].Sets)
				assert.Equal(t, tt.workout.Entries[i].OrderIndex, retrieved.Entries[i].OrderIndex)
			}
		})
	}
}

func IntPtr(i int) *int {
	return &i
}

func FloatPtr(i float64) *float64 {
	return &i
}
