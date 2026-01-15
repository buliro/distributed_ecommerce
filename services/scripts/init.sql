-- Create the database
CREATE DATABASE storedb;

--Create the user and grant priviledges
CREATE USER storedb_user WITH PASSWORD 'storedb_password';
GRANT ALL PRIVILEGES ON DATABASE storedb TO storedb_user;
