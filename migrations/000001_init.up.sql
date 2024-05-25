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

INSERT INTO cryptowallet.currencies (id,name,code,max_supply,description) VALUES
 ('43dc5ab9-f8ec-48f6-8ae7-348096691835','USDT','4',389423168412894,'Dollars'),
 ('97f5148a-f8ec-48f6-8ae7-348096691835','TON','2',9999,''),
 ('835a9c62-caaf-4891-91a7-4f9137f3f815','BTC','1',389423168412894,'BITOC'),
 ('43dc5ab9-3ff7-4303-ae06-9aafe0114822','ETH','3',9999,'');

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

    `predict_time` DateTime,

    `created_at` DateTime DEFAULT now(),

    `crypto_id` UUID,

    `market_cap` Float64,

    `predict` Float64
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


-- cryptowallet.volatilities определение

CREATE TABLE cryptowallet.volatilities
(

    `id` UUID,

    `crypto_id` UUID,

    `volatility` Float64,

    `created_at` DateTime DEFAULT now()
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

    `priceToVolatility` Float64

    `created_at` DateTime DEFAULT now(),

    `crypto_id` UUID,
)
ENGINE = MergeTree
PRIMARY KEY tuple(id)
ORDER BY id
SETTINGS index_granularity = 8192;


-- cryptowallet.predictions определение

CREATE TABLE cryptowallet.predictions
(

    `id` UUID,

    `value` Float64,

    `updated_at` DateTime DEFAULT now(),

    `crypto_id` UUID,
)
ENGINE = MergeTree
PRIMARY KEY tuple(id)
ORDER BY id
SETTINGS index_granularity = 8192;

