-- name: GetCompanyByID :one
SELECT * FROM companies WHERE id = $1;

-- name: CreateCompany :one
INSERT INTO companies (
    name,
    representative_name,
    phone_number,
    zip_code,
    address
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: UpdateCompany :one
UPDATE companies SET
    name = $2,
    representative_name = $3,
    phone_number = $4,
    zip_code = $5,
    address = $6,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;
