BEGIN;

CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(26) NOT NULL,
    username VARCHAR(12) NOT NULL,
    email VARCHAR(25) NOT NULL,
    passhash VARCHAR NOT NULL,
    role VARCHAR(25) NOT NULL
);

CREATE TABLE IF NOT EXISTS luxury_items (
    id VARCHAR(26) NOT NULL,
    brand CHAR(26) NOT NULL,
    price INT NOT NULL,
    ownerid VARCHAR(26) NOT NULL,
    tokenURI VARCHAR(64) NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(26) NOT NULL,
    txHash VARCHAR(64),
    txType VARCHAR(24) NOT NULL,
    fromAddress VARCHAR(64),
    toAddress VARCHAR(64)
);

COMMIT;
