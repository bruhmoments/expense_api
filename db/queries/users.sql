-- name: CreateUser :one
INSERT INTO users (username, password) 
VALUES ($1, $2)
RETURNING id, username;

-- name: GetUserByUsername :one
SELECT id, username, password 
FROM users 
WHERE username = $1 
LIMIT 1;
