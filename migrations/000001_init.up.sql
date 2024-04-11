CREATE TABLE "currencies" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(100) NOT NULL,
  "code" varchar(5) NOT NULL,
  "max_supply" bigint DEFAULT 99999999999999,
  "desciption" text
);

CREATE TABLE "trade_volume_1h" (
  "id" bigserial PRIMARY KEY,
  "crypto_id" bigint NOT NULL,
  "trade_volume" float NOT NULL,
  "time_diff" datetime
);

CREATE TABLE "indices" (
  "id" bigserial PRIMARY KEY,
  "crypto_id" bigint NOT NULL,
  "price_index" integer NOT NULL,
  "volume_index" float NOT NULL,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "prices" (
  "id" bigserial PRIMARY KEY,
  "crypto_id" bigint NOT NULL,
  "value" money NOT NULL,
  "market_cap" integer NOT NULL,
  "time_diff" datetime,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "transaction_per_day" (
  "id" bigserial PRIMARY KEY,
  "crypto_id" bigint NOT NULL,
  "trans_value" bigint NOT NULL,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "supplies" (
  "id" bigserial PRIMARY KEY,
  "crypto_id" bigint NOT NULL,
  "total_supply" bigint,
  "created_at" timestamp DEFAULT (now())
);

ALTER TABLE "trade_volume_1h" ADD FOREIGN KEY ("crypto_id") REFERENCES "currencies" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "indices" ADD FOREIGN KEY ("crypto_id") REFERENCES "currencies" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "prices" ADD FOREIGN KEY ("crypto_id") REFERENCES "currencies" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "transaction_per_day" ADD FOREIGN KEY ("crypto_id") REFERENCES "currencies" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE "supplies" ADD FOREIGN KEY ("crypto_id") REFERENCES "currencies" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;

