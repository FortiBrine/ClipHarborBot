-- name: GetUserLanguage :one
SELECT language
FROM users
WHERE id = $1;

-- name: UpsertUserLanguage :exec
INSERT INTO users (id, language)
VALUES ($1, $2)
ON CONFLICT (id) DO UPDATE SET language = EXCLUDED.language;
