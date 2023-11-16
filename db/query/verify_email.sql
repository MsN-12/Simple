-- name: CreateVerifyEmail :one
INSERT INTO verify_emails(
                           username,
                           email,
                           secret_code
) values (
          $1,$2,$3
) RETURNING *;