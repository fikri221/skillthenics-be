package workout

import (
	"time"
)

type Exercise struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	MuscleGroup *string `json:"muscle_group"`
	Difficulty  *string `json:"difficulty"`
}

type WorkoutSession struct {
	ID          string    `json:"id"`
	UserID      int32     `json:"user_id"`
	SessionDate time.Time `json:"session_date"`
	Duration    *int32    `json:"duration"`
	Calories    *int32    `json:"calories"`
	Notes       *string   `json:"notes"`
}

type WorkoutExercise struct {
	ID               string  `json:"id"`
	WorkoutSessionID string  `json:"workout_session_id"`
	ExerciseID       string  `json:"exercise_id"`
	ExerciseName     string  `json:"exercise_name"`
	MuscleGroup      *string `json:"muscle_group"`
	Notes            *string `json:"notes"`
}

type ExerciseSet struct {
	ID                string  `json:"id"`
	WorkoutExerciseID string  `json:"workout_exercise_id"`
	SetNumber         int32   `json:"set_number"`
	Reps              *int32  `json:"reps"`
	Weight            *string `json:"weight"`
	RestSeconds       *int32  `json:"rest_seconds"`
}

type FullWorkoutSession struct {
	Session         WorkoutSession          `json:"session"`
	Exercises       []WorkoutExercise       `json:"exercises"`
	ExerciseSetsMap map[int32][]ExerciseSet `json:"exercise_sets_map"`
}
