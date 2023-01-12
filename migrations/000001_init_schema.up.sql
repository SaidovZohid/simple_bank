-- * table for accounts
CREATE TABLE "accounts"(
    "id" BIGSERIAL NOT NULL,
    "owner" VARCHAR(255) NOT NULL,
    "balance" NUMERIC(18, 2) NOT NULL,
    "currency" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
-- * altering table accounts, making id primary key to make relation with another tables
ALTER TABLE
    "accounts" ADD PRIMARY KEY("id");

-- * table for transfers
CREATE TABLE "transfers"(
    "id" BIGSERIAL NOT NULL,
    "from_account_id" BIGINT NOT NULL,
    "to_account_id" BIGINT NOT NULL,
    "amount" NUMERIC(18, 2) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
-- * altering table tranfers, making id primary key to make relation with another tables
ALTER TABLE
    "transfers" ADD PRIMARY KEY("id");
    
-- * table for entries 
CREATE TABLE "entries"(
    "id" BIGSERIAL NOT NULL,
    "account_id" BIGINT NOT NULL,
    "amount" NUMERIC(18, 2) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
-- * altering table entries, making id primary key to make relation with another tables
ALTER TABLE
    "entries" ADD PRIMARY KEY("id");
-- * adding foreign key and making relation with accounts's id
ALTER TABLE
    "transfers" ADD CONSTRAINT "transfers_to_account_id_foreign" FOREIGN KEY("to_account_id") REFERENCES "accounts"("id");
-- * adding foreign key and making relation with accounts's id
ALTER TABLE
    "transfers" ADD CONSTRAINT "transfers_from_account_id_foreign" FOREIGN KEY("from_account_id") REFERENCES "accounts"("id");
-- * adding foreign key and making relation with accounts's id
ALTER TABLE
    "entries" ADD CONSTRAINT "entries_account_id_foreign" FOREIGN KEY("account_id") REFERENCES "accounts"("id");