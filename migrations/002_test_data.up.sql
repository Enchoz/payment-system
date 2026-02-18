-- Insert test users
INSERT INTO users (id, email, name, createdAt, updatedAt) VALUES
    ('11111111-1111-1111-1111-111111111111', 'alice@example.com', 'Alice Smith', NOW(), NOW()),
    ('22222222-2222-2222-2222-222222222222', 'bob@example.com', 'Bob Jones', NOW(), NOW()),
    ('33333333-3333-3333-3333-333333333333', 'charlie@example.com', 'Charlie Brown', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert accounts for Alice (USD, EUR, GBP)
INSERT INTO accounts (id, userId, currency, balance, availableBalance, createdAt, updatedAt) VALUES
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111', 'USD', 10000.0000, 10000.0000, NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '11111111-1111-1111-1111-111111111111', 'EUR', 5000.0000, 5000.0000, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-cccccccccccc', '11111111-1111-1111-1111-111111111111', 'GBP', 3000.0000, 3000.0000, NOW(), NOW())
ON CONFLICT (userId, currency) DO NOTHING;

-- Insert accounts for Bob (USD, EUR, GBP)
INSERT INTO accounts (id, userId, currency, balance, availableBalance, createdAt, updatedAt) VALUES
    ('dddddddd-dddd-dddd-dddd-dddddddddddd', '22222222-2222-2222-2222-222222222222', 'USD', 5000.0000, 5000.0000, NOW(), NOW()),
    ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', '22222222-2222-2222-2222-222222222222', 'EUR', 2500.0000, 2500.0000, NOW(), NOW()),
    ('ffffffff-ffff-ffff-ffff-ffffffffffff', '22222222-2222-2222-2222-222222222222', 'GBP', 1500.0000, 1500.0000, NOW(), NOW())
ON CONFLICT (userId, currency) DO NOTHING;

-- Insert accounts for Charlie (USD, EUR, GBP)
INSERT INTO accounts (id, userId, currency, balance, availableBalance, createdAt, updatedAt) VALUES
    ('11111111-1111-1111-1111-111111111112', '33333333-3333-3333-3333-333333333333', 'USD', 2000.0000, 2000.0000, NOW(), NOW()),
    ('11111111-1111-1111-1111-111111111113', '33333333-3333-3333-3333-333333333333', 'EUR', 1000.0000, 1000.0000, NOW(), NOW()),
    ('11111111-1111-1111-1111-111111111114', '33333333-3333-3333-3333-333333333333', 'GBP', 500.0000, 500.0000, NOW(), NOW())
ON CONFLICT (userId, currency) DO NOTHING;
