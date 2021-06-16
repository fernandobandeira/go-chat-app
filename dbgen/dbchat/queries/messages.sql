-- messages.sql

-- name: GetMessages :many
SELECT
    "id",
    "author",
    "text",
    "datetime"
FROM
    "messages"
ORDER BY
    "datetime"
LIMIT 50
OFFSET 0
;

-- name: AddMessage :one
INSERT INTO
    "messages"
(
    "author",
    "text"
)
VALUES
(
    @author::text,
    @text::text
)
RETURNING *;