-- name: GetInvoiceByID :one
SELECT * FROM invoices WHERE id = $1;

-- name: GetInvoiceByIDAndCompanyID :one
SELECT * FROM invoices WHERE id = $1 AND company_id = $2;

-- name: GetInvoicesByCompanyID :many
SELECT * FROM invoices WHERE company_id = $1 ORDER BY due_date DESC, id DESC;

-- name: GetInvoicesByCompanyIDAndDateRange :many
SELECT * FROM invoices
WHERE company_id = $1
  AND due_date >= $2
  AND due_date <= $3
ORDER BY due_date ASC, id ASC;

-- name: CreateInvoice :one
INSERT INTO invoices (
    company_id,
    vendor_id,
    vendor_bank_account_id,
    issue_date,
    payment_amount,
    fee,
    fee_rate,
    tax,
    tax_rate,
    total_amount,
    due_date,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
) RETURNING *;

-- name: UpdateInvoiceStatus :one
UPDATE invoices SET
    status = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: CountInvoicesByCompanyIDAndDateRange :one
SELECT COUNT(*) FROM invoices
WHERE company_id = $1
  AND due_date >= $2
  AND due_date <= $3;
