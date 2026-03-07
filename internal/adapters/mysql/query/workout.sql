-- name: CreateExercise :execresult
INSERT INTO exercises (id, name, description, muscle_group, difficulty)
VALUES (?, ?, ?, ?, ?);

-- name: GetExerciseByID :one
SELECT * FROM exercises WHERE id = ?;

-- name: ListExercises :many
SELECT * FROM exercises ORDER BY name ASC;

-- name: UpdateExercise :execresult
UPDATE exercises SET name = ?, description = ?, muscle_group = ?, difficulty = ? WHERE id = ?;

-- name: CreateWorkoutSession :execresult
INSERT INTO workout_sessions (id, user_id, session_date, duration_minutes, calories_burned, notes)
VALUES (?, ?, ?, ?, ?, ?);

-- name: GetWorkoutSessionByID :one
SELECT * FROM workout_sessions WHERE id = ?;

-- name: ListWorkoutSessionsByUser :many
SELECT * FROM workout_sessions WHERE user_id = ? ORDER BY session_date DESC;

-- name: UpdateWorkoutSession :execresult
UPDATE workout_sessions SET session_date = ?, duration_minutes = ?, calories_burned = ?, notes = ? WHERE id = ?;

-- name: CreateWorkoutExercise :execresult
INSERT INTO workout_exercises (id, workout_session_id, exercise_id, notes)
VALUES (?, ?, ?, ?);

-- name: ListWorkoutExercisesBySession :many
SELECT 
    we.*, 
    e.name as exercise_name,
    e.muscle_group
FROM workout_exercises we
JOIN exercises e ON we.exercise_id = e.id
WHERE we.workout_session_id = ?
ORDER BY we.created_at ASC;

-- name: GetWorkoutExerciseByID :one
SELECT * FROM workout_exercises WHERE id = ?;

-- name: UpdateWorkoutExercise :execresult
UPDATE workout_exercises SET exercise_id = ?, notes = ? WHERE id = ?;

-- name: CreateExerciseSet :execresult
INSERT INTO exercise_sets (workout_exercise_id, set_number, reps, weight, rest_seconds)
VALUES (?, ?, ?, ?, ?);

-- name: ListSetsByWorkoutExercise :many
SELECT * FROM exercise_sets 
WHERE workout_exercise_id = ? 
ORDER BY set_number ASC;

-- name: UpdateExerciseSet :execresult
UPDATE exercise_sets SET set_number = ?, reps = ?, weight = ?, rest_seconds = ? WHERE id = ?;

-- name: GetFullWorkoutSession :many
-- Fetching session, exercises, and sets in one join structure (if needed)
-- Note: This might require manual mapping in Go depending on how you want to handle the structure.
SELECT 
    ws.id as session_id,
    ws.session_date,
    ws.duration_minutes,
    we.id as workout_exercise_id,
    e.name as exercise_name,
    es.id as set_id,
    es.set_number,
    es.reps,
    es.weight
FROM workout_sessions ws
LEFT JOIN workout_exercises we ON ws.id = we.workout_session_id
LEFT JOIN exercises e ON we.exercise_id = e.id
LEFT JOIN exercise_sets es ON we.id = es.workout_exercise_id
WHERE ws.id = ?;
