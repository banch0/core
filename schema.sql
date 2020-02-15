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