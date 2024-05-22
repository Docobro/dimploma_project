-- cryptowallet.currencies определение

CREATE TABLE cryptowallet.currencies
(

    `id` UUID,

    `name` String,

    `code` String,

    `max_supply` Int64 DEFAULT 9999,

    `description` String
)
ENGINE = MergeTree
PRIMARY KEY tuple(id)
ORDER BY id
SETTINGS index_granularity = 8192;


-- cryptowallet.indices определение

CREATE TABLE cryptowallet.indices
(

    `id` UUID,

    `created_at` DateTime DEFAULT now(),

    `crypto_id` UUID,

    `price_index` Float64
)
ENGINE = MergeTree
PRIMARY KEY tuple(id)
ORDER BY id
SETTINGS index_granularity = 8192;


-- cryptowallet.prices определение

CREATE TABLE cryptowallet.prices
(

    `id` UUID,

    `value` Decimal(18,
 2),

    `time_diff` DateTime,

    `created_at` DateTime DEFAULT now(),

    `crypto_id` UUID,

    `market_cap` Float64
)
ENGINE = MergeTree
PRIMARY KEY tuple(id)
ORDER BY id
SETTINGS index_granularity = 8192;


-- cryptowallet.supplies определение

CREATE TABLE cryptowallet.supplies
(

    `id` UUID,

    `crypto_id` UUID,

    `total_supply` Float64,

    `created_at` DateTime DEFAULT now()
)
ENGINE = MergeTree
PRIMARY KEY tuple(id)
ORDER BY id
SETTINGS index_granularity = 8192;


-- cryptowallet.trade_volume_1m определение

CREATE TABLE cryptowallet.trade_volume_1m
(

    `id` UUID,

    `crypto_id` UUID,

    `trade_volume` Float64,

    `created_at` DateTime DEFAULT now(),
)
ENGINE = MergeTree
PRIMARY KEY tuple(id)
ORDER BY id
SETTINGS index_granularity = 8192;


-- cryptowallet.pearsonpearson_price_volume определение

CREATE TABLE cryptowallet.pearson_correlation
(

    `id` UUID,

    `priceToVolume` Float64,

    `priceToMarketCap` Float64,

    `created_at` DateTime DEFAULT now(),

    `crypto_id` UUID,
)
ENGINE = MergeTree
PRIMARY KEY tuple(id)
ORDER BY id
SETTINGS index_granularity = 8192;
