CREATE TABLE cryptowallet.currencies (
  id UUID PRIMARY KEY,
  name String NOT NULL,
  code String NOT NULL,
  max_supply Int64 DEFAULT 9999,
  description String
) ENGINE = MergeTree()
ORDER BY id;

CREATE TABLE cryptowallet.trade_volume_1h (
  id UUID PRIMARY KEY,
  crypto_id Int64 NOT NULL,
  trade_volume Float32 NOT NULL,
  time_diff DateTime
) ENGINE = MergeTree()
ORDER BY id;

CREATE TABLE cryptowallet.indices (
  id UUID PRIMARY KEY,
  crypto_id Int64 NOT NULL,
  price_index Int32 NOT NULL,
  volume_index Float32 NOT NULL,
  created_at DateTime DEFAULT now()
) ENGINE = MergeTree()
ORDER BY id;

CREATE TABLE cryptowallet.prices (
  id UUID PRIMARY KEY,
  crypto_id Int64 NOT NULL,
  value Decimal(18, 2) NOT NULL,
  market_cap Int32 NOT NULL,
  time_diff DateTime,
  created_at DateTime DEFAULT now()
) ENGINE = MergeTree()
ORDER BY id;

CREATE TABLE cryptowallet.transaction_per_day (
  id UUID PRIMARY KEY,
  crypto_id Int64 NOT NULL,
  trans_value Int64 NOT NULL,
  created_at DateTime DEFAULT now()
) ENGINE = MergeTree()
ORDER BY id;

CREATE TABLE cryptowallet.supplies (
  id UUID PRIMARY KEY,
  crypto_id Int64 NOT NULL,
  total_supply Int64,
  created_at DateTime DEFAULT now()
) ENGINE = MergeTree()
ORDER BY id;

ALTER TABLE cryptowallet.trade_volume_1h ADD COLUMN IF NOT EXISTS crypto_id Int64;
ALTER TABLE cryptowallet.indices ADD COLUMN IF NOT EXISTS crypto_id Int64;
ALTER TABLE cryptowallet.prices ADD COLUMN IF NOT EXISTS crypto_id Int64;
ALTER TABLE cryptowallet.transaction_per_day ADD COLUMN IF NOT EXISTS crypto_id Int64;
ALTER TABLE cryptowallet.supplies ADD COLUMN IF NOT EXISTS crypto_id Int64;

