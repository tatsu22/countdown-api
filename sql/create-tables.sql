CREATE DATABASE countdown_api;

\c countdown_api;

CREATE TABLE countdown_games (
    nums integer[],
    goal integer,
    equation_array text[],
    equation text,
    complete boolean,
    time_taken text,
    nodes_calculated integer
);