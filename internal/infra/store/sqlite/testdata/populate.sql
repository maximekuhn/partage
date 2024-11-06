-- This file contains SQL queries to populate the database with some initial data.
-- It can be helpful when running queries directly using SQLite3.
-- To import data into a SQLite3 file, run (in the current directory):
-- $ sqlite3 <DB_FILE_PATH.sqlite3> < populate.sql

-- Clear content of all tables
DELETE FROM user;
DELETE FROM partage_group;
DELETE FROM partage_group_user;

-- Create users
INSERT INTO user (id, nickname, email, created_at) 
VALUES ('87822DC7-9DB9-4565-9301-87006FE1C35B', 'john', 'john@acme.com', '2013-10-07T08:23:19');

INSERT INTO user (id, nickname, email, created_at) 
VALUES ('5D70BC4C-5841-4BEB-93D5-61AF4BE09478', 'jane', 'jane@acme.com', '2013-10-07T08:23:19');

INSERT INTO user (id, nickname, email, created_at) 
VALUES ('3E70370F-FD25-4D3E-B5BB-777E49A4DABC', 'bob', 'bob@acme.com', '2013-10-07T08:23:19');

-- Create groups
INSERT INTO partage_group (id, name, owner, created_at)
VALUES ('65C509DD-D1BC-4E8C-834E-57E30D100A2D', 'group 1', '87822DC7-9DB9-4565-9301-87006FE1C35B', '2014-10-07T08:23:19');

-- Add users into groups
INSERT INTO partage_group_user (group_id, user_id)
VALUES ('65C509DD-D1BC-4E8C-834E-57E30D100A2D', '87822DC7-9DB9-4565-9301-87006FE1C35B');

INSERT INTO partage_group_user (group_id, user_id)
VALUES ('65C509DD-D1BC-4E8C-834E-57E30D100A2D', '5D70BC4C-5841-4BEB-93D5-61AF4BE09478');

