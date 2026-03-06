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
INSERT INTO workout_sessions (user_id, session_date, duration_minutes, calories_burned, notes)
VALUES (?, ?, ?, ?, ?);

-- name: GetWorkoutSession :one
SELECT * FROM workout_sessions WHERE id = ?;

-- name: ListWorkoutSessionsByUser :many
SELECT * FROM workout_sessions WHERE user_id = ? ORDER BY session_date DESC;

-- name: CreateWorkoutExercise :execresult
INSERT INTO workout_exercises (workout_session_id, exercise_id, notes)
VALUES (?, ?, ?);

-- name: ListWorkoutExercisesBySession :many
SELECT 
    we.*, 
    e.name as exercise_name,
    e.muscle_group
FROM workout_exercises we
JOIN exercises e ON we.exercise_id = e.id
WHERE we.workout_session_id = ?
ORDER BY we.created_at ASC;

-- name: CreateExerciseSet :execresult
INSERT INTO exercise_sets (workout_exercise_id, set_number, reps, weight, rest_seconds)
VALUES (?, ?, ?, ?, ?);

-- name: ListSetsByWorkoutExercise :many
SELECT * FROM exercise_sets 
WHERE workout_exercise_id = ? 
ORDER BY set_number ASC;

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
