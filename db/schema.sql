-- スーパー支払い君.com データベーススキーマ
-- sqldef (psqldef) 用

-- ステータス型
-- pending=未処理, processing=処理中, paid=支払済, error=エラー
CREATE TYPE invoice_status AS ENUM ('pending', 'processing', 'paid', 'error');

-- 企業テーブル
CREATE TABLE companies (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,                -- 法人名
    representative_name VARCHAR(255) NOT NULL, -- 代表者名
    phone_number VARCHAR(16) NOT NULL,         -- 電話番号 (E.164形式: +81312345678)
    zip_code VARCHAR(10) NOT NULL,             -- 郵便番号
    address VARCHAR(500) NOT NULL,             -- 住所
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ユーザーテーブル（企業に紐づく）
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT NOT NULL REFERENCES companies(id) ON DELETE CASCADE, -- 所属企業ID
    name VARCHAR(255) NOT NULL,          -- 氏名
    email VARCHAR(255) NOT NULL UNIQUE,  -- メールアドレス
    password_hash VARCHAR(255) NOT NULL, -- パスワードハッシュ (bcrypt)
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_company_id ON users(company_id);
CREATE INDEX idx_users_email ON users(email);

-- 取引先テーブル（企業に紐づく）
CREATE TABLE vendors (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT NOT NULL REFERENCES companies(id) ON DELETE CASCADE, -- 所属企業ID
    name VARCHAR(255) NOT NULL,                -- 法人名
    representative_name VARCHAR(255) NOT NULL, -- 代表者名
    phone_number VARCHAR(16) NOT NULL,         -- 電話番号 (E.164形式: +81312345678)
    zip_code VARCHAR(10) NOT NULL,             -- 郵便番号
    address VARCHAR(500) NOT NULL,             -- 住所
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_vendors_company_id ON vendors(company_id);

-- 取引先銀行口座テーブル（取引先に紐づく）
CREATE TABLE vendor_bank_accounts (
    id BIGSERIAL PRIMARY KEY,
    vendor_id BIGINT NOT NULL REFERENCES vendors(id) ON DELETE CASCADE, -- 取引先ID
    bank_name VARCHAR(255) NOT NULL,           -- 銀行名
    branch_name VARCHAR(255) NOT NULL,         -- 支店名
    account_number VARCHAR(20) NOT NULL,       -- 口座番号
    account_holder_name VARCHAR(255) NOT NULL, -- 口座名義
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_vendor_bank_accounts_vendor_id ON vendor_bank_accounts(vendor_id);

-- 請求書テーブル（企業・取引先に紐づく）
CREATE TABLE invoices (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,              -- 企業ID
    vendor_id BIGINT NOT NULL REFERENCES vendors(id) ON DELETE RESTRICT,                -- 取引先ID
    vendor_bank_account_id BIGINT NOT NULL REFERENCES vendor_bank_accounts(id) ON DELETE RESTRICT, -- 振込先銀行口座ID
    issue_date DATE NOT NULL,                                    -- 発行日
    payment_amount BIGINT NOT NULL CHECK (payment_amount > 0),   -- 支払金額
    fee BIGINT NOT NULL CHECK (fee >= 0),                        -- 手数料 (payment_amount * fee_rate)
    fee_rate DECIMAL(5, 4) NOT NULL DEFAULT 0.04,                -- 手数料率 (デフォルト: 4%)
    tax BIGINT NOT NULL CHECK (tax >= 0),                        -- 消費税 (fee * tax_rate)
    tax_rate DECIMAL(5, 4) NOT NULL DEFAULT 0.10,                -- 消費税率 (デフォルト: 10%)
    total_amount BIGINT NOT NULL CHECK (total_amount > 0),       -- 請求金額 (payment_amount + fee + tax)
    due_date DATE NOT NULL,                                      -- 支払期日
    status invoice_status NOT NULL DEFAULT 'pending',            -- ステータス
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 支払期日での範囲検索用インデックス
CREATE INDEX idx_invoices_company_id ON invoices(company_id);
CREATE INDEX idx_invoices_due_date ON invoices(due_date);
CREATE INDEX idx_invoices_company_due_date ON invoices(company_id, due_date);
CREATE INDEX idx_invoices_status ON invoices(status);
