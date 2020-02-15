CREATE TABLE IF NOT EXISTS managers
(
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    name    TEXT    NOT NULL,
    login   TEXT    NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT ,
    name TEXT NOT NULL ,
    login TEXT NOT NULL ,
    password TEXT NOT NULL,
    account TEXT NOT NULL ,
    balance INTEGER DEFAULT 0,
    phone TEXT NOT NULL
);

CREATE TABLE services (
                          id INTEGER PRIMARY KEY AUTOINCREMENT,
                          name TEXT NOT NULL ,
                          price INTEGER NOT NULL
);

CREATE TABLE atm_list (
                          id INTEGER PRIMARY KEY AUTOINCREMENT,
                          name TEXT NOT NULL ,
                          address TEXT NOT NULL
);

CREATE TABLE sales (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    service_id INTEGER,
    price INTEGER
);

DROP TABLE managers;
DROP TABLE users;
DROP TABLE services;
DROP TABLE atm_list;
DROP TABLE sales;

UPDATE users SET balance = 100 WHERE id = 7;
UPDATE users SET balance = balance - 10 WHERE id = 7;
SELECT balance FROM users WHERE id = 1;