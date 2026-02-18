-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP NOT NULL DEFAULT NOW(),
    updatedAt TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Accounts table (one per user per currency)
CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    userId UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    currency VARCHAR(3) NOT NULL,
    balance DECIMAL(19, 4) NOT NULL DEFAULT 0,
    availableBalance DECIMAL(19, 4) NOT NULL DEFAULT 0,
    createdAt TIMESTAMP NOT NULL DEFAULT NOW(),
    updatedAt TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(userId, currency),
    CHECK (balance >= 0),
    CHECK (availableBalance >= 0)
);

CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts(userId);
CREATE INDEX IF NOT EXISTS idx_accounts_currency ON accounts(currency);

-- External Accounts (cached/validated accounts)
CREATE TABLE IF NOT EXISTS external_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    userId UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Account details
    accountNumber VARCHAR(255) NOT NULL,
    routingNumber VARCHAR(255),
    iban VARCHAR(34),
    swiftCode VARCHAR(11),
    bankName VARCHAR(255),
    accountHolderName VARCHAR(255),
    
    -- Currency and type
    currency VARCHAR(3) NOT NULL,
    accountType VARCHAR(50) NOT NULL,
    countryCode VARCHAR(2),
    
    -- Validation status
    validationStatus VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    validatedAt TIMESTAMP,
    validationExpiresAt TIMESTAMP,
    lastValidatedAt TIMESTAMP,
    
    -- Metadata
    nickname VARCHAR(100),
    isDefault BOOLEAN DEFAULT FALSE,
    metadata JSONB,
    
    -- Audit
    createdAt TIMESTAMP NOT NULL DEFAULT NOW(),
    updatedAt TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_external_accounts_user_id ON external_accounts(userId);
CREATE INDEX IF NOT EXISTS idx_external_accounts_validation_status ON external_accounts(validationStatus);
CREATE INDEX IF NOT EXISTS idx_external_accounts_hash ON external_accounts(accountNumber, routingNumber, currency);

-- Payment Requests (staging table for external payments)
CREATE TABLE IF NOT EXISTS payment_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    requestReference VARCHAR(255) UNIQUE NOT NULL,
    
    -- Payment details
    fromUserId UUID NOT NULL REFERENCES users(id),
    externalAccountId UUID REFERENCES external_accounts(id),
    toUserId UUID REFERENCES users(id),
    
    -- Amount details
    amount DECIMAL(19, 4) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    targetCurrency VARCHAR(3),
    
    -- External account details (NULL for internal payments)
    externalAccountNumber VARCHAR(255),
    externalRoutingNumber VARCHAR(255),
    externalIban VARCHAR(34),
    externalSwiftCode VARCHAR(11),
    externalBankName VARCHAR(255),
    externalAccountHolderName VARCHAR(255),
    externalCountryCode VARCHAR(2),
    
    -- Status and processing
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    paymentType VARCHAR(20) NOT NULL,
    
    -- Validation results
    validationResults JSONB,
    validationErrors JSONB,
    requiresManualReview BOOLEAN DEFAULT FALSE,
    
    -- Exchange rate (snapshot)
    exchangeRate DECIMAL(19, 8),
    convertedAmount DECIMAL(19, 4),
    
    -- Processing
    failureReason TEXT,
    processedAt TIMESTAMP,
    
    -- Metadata
    metadata JSONB,
    idempotencyKey VARCHAR(255) UNIQUE,
    
    createdAt TIMESTAMP NOT NULL DEFAULT NOW(),
    updatedAt TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CHECK (amount > 0),
    CHECK (
        -- Internal payment validation
        (paymentType = 'INTERNAL' AND 
         toUserId IS NOT NULL AND 
         externalAccountId IS NULL AND
         externalAccountNumber IS NULL AND
         externalIban IS NULL) OR
        -- External payment validation
        (paymentType = 'EXTERNAL' AND 
         toUserId IS NULL AND
         (externalAccountId IS NOT NULL OR 
          externalAccountNumber IS NOT NULL OR
          externalIban IS NOT NULL))
    )
);

CREATE INDEX IF NOT EXISTS idx_payment_requests_from_user_id ON payment_requests(fromUserId);
CREATE INDEX IF NOT EXISTS idx_payment_requests_to_user_id ON payment_requests(toUserId);
CREATE INDEX IF NOT EXISTS idx_payment_requests_status ON payment_requests(status);
CREATE INDEX IF NOT EXISTS idx_payment_requests_idempotency_key ON payment_requests(idempotencyKey);

-- Payments (final processed payments)
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    paymentReference VARCHAR(255) UNIQUE NOT NULL,
    requestId UUID REFERENCES payment_requests(id),
    
    -- Payment details
    fromUserId UUID NOT NULL REFERENCES users(id),
    toUserId UUID REFERENCES users(id),
    externalAccountId UUID REFERENCES external_accounts(id),
    
    -- Amount details
    amount DECIMAL(19, 4) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    exchangeRate DECIMAL(19, 8),
    convertedAmount DECIMAL(19, 4),
    convertedCurrency VARCHAR(3),
    
    -- Status
    status VARCHAR(20) NOT NULL,
    paymentType VARCHAR(20) NOT NULL,
    
    -- External payment details (denormalized for audit)
    externalAccountNumber VARCHAR(255),
    externalBankName VARCHAR(255),
    
    -- Processing
    failureReason TEXT,
    processedAt TIMESTAMP,
    completedAt TIMESTAMP,
    
    -- Metadata
    metadata JSONB,
    
    createdAt TIMESTAMP NOT NULL DEFAULT NOW(),
    updatedAt TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CHECK (amount > 0)
);

CREATE INDEX IF NOT EXISTS idx_payments_from_user_id ON payments(fromUserId);
CREATE INDEX IF NOT EXISTS idx_payments_to_user_id ON payments(toUserId);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);
CREATE INDEX IF NOT EXISTS idx_payments_request_id ON payments(requestId);

-- Ledger Entries (double-entry bookkeeping)
CREATE TABLE IF NOT EXISTS ledger_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    paymentId UUID NOT NULL REFERENCES payments(id),
    accountId UUID NOT NULL REFERENCES accounts(id),
    entryType VARCHAR(20) NOT NULL,
    amount DECIMAL(19, 4) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    balanceAfter DECIMAL(19, 4) NOT NULL,
    createdAt TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ledger_entries_payment_id ON ledger_entries(paymentId);
CREATE INDEX IF NOT EXISTS idx_ledger_entries_account_id ON ledger_entries(accountId);
CREATE INDEX IF NOT EXISTS idx_ledger_entries_created_at ON ledger_entries(createdAt);

-- Payment Holds (for pending transactions)
CREATE TABLE IF NOT EXISTS payment_holds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    accountId UUID NOT NULL REFERENCES accounts(id),
    requestId UUID REFERENCES payment_requests(id),
    paymentId UUID REFERENCES payments(id),
    amount DECIMAL(19, 4) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    expiresAt TIMESTAMP,
    releasedAt TIMESTAMP,
    createdAt TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payment_holds_account_id ON payment_holds(accountId);
CREATE INDEX IF NOT EXISTS idx_payment_holds_request_id ON payment_holds(requestId);

-- Exchange Rates (snapshot for audit)
CREATE TABLE IF NOT EXISTS exchange_rates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fromCurrency VARCHAR(3) NOT NULL,
    toCurrency VARCHAR(3) NOT NULL,
    rate DECIMAL(19, 8) NOT NULL,
    effectiveAt TIMESTAMP NOT NULL,
    source VARCHAR(50) NOT NULL,
    createdAt TIMESTAMP NOT NULL DEFAULT NOW(),
    
    UNIQUE(fromCurrency, toCurrency, effectiveAt)
);

CREATE INDEX IF NOT EXISTS idx_exchange_rates_effective_at ON exchange_rates(effectiveAt);
CREATE INDEX IF NOT EXISTS idx_exchange_rates_currencies ON exchange_rates(fromCurrency, toCurrency);
