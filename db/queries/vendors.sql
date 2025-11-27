-- name: GetVendorByID :one
SELECT * FROM vendors WHERE id = $1;

-- name: GetVendorsByCompanyID :many
SELECT * FROM vendors WHERE company_id = $1 ORDER BY id;

-- name: GetVendorByIDAndCompanyID :one
SELECT * FROM vendors WHERE id = $1 AND company_id = $2;

-- name: CreateVendor :one
INSERT INTO vendors (
    company_id,
    name,
    representative_name,
    phone_number,
    zip_code,
    address
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: UpdateVendor :one
UPDATE vendors SET
    name = $2,
    representative_name = $3,
    phone_number = $4,
    zip_code = $5,
    address = $6,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;
