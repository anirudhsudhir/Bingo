-- Create a database for the application
CREATE DATABASE bingo CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE bingo;

-- Create a `snips` table.
CREATE TABLE snips (
id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
title VARCHAR(100) NOT NULL,
content TEXT NOT NULL,
created DATETIME NOT NULL,
expires DATETIME NOT NULL
);
-- Create an index
CREATE INDEX idx_snips_created ON snips(created);

-- Create a new user for accessing the database from the pastebin
CREATE USER 'web'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE ON bingo.* TO 'web'@'localhost'; 
-- Replace 'pass' with a password of your choice
ALTER USER 'web'@'localhost' IDENTIFIED BY 'pass';

