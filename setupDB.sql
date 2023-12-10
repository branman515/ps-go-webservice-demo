CREATE DATABASE readinglist;

CREATE ROLE readinglist WITH LOGIN PASSWORD 'pa$$w0rd';

CREATE TABLE IF NOT EXISTS books (
    id bigserial PRIMARY KEY,  
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    published integer NOT NULL,
    pages integer NOT NULL,
    genres text[] NOT NULL,
    version integer NOT NULL DEFAULT 1,
    rating FLOAT NOT NULL
);