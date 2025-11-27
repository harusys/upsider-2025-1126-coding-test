-- name: GetVendorBankAccountByID :one
SELECT * FROM vendor_bank_accounts WHERE id = $1;

-- name: GetVendorBankAccountsByVendorID :many
SELECT * FROM vendor_bank_accounts WHERE vendor_id = $1 ORDER BY id;

-- name: GetVendorBankAccountByIDAndVendorID :one
SELECT * FROM vendor_bank_accounts WHERE id = $1 AND vendor_id = $2;

-- name: CreateVendorBankAccount :one
INSERT INTO vendor_bank_accounts (
    vendor_id,
    bank_name,
    branch_name,
    account_number,
    account_holder_name
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: UpdateVendorBankAccount :one
UPDATE vendor_bank_accounts SET
    bank_name = $2,
    branch_name = $3,
    account_number = $4,
    account_holder_name = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;
