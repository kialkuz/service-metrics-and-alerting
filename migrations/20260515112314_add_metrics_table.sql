-- +goose Up
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS public.metrics (
	id serial4 NOT NULL,
	"type" varchar(50) NOT NULL,
	"name" varchar(50) NOT NULL,
	value double precision NULL,
	delta int8 NULL,
	CONSTRAINT users_pkey PRIMARY KEY (id)
);

-- +goose Down
SELECT 'down SQL query';
DROP TABLE IF EXISTS public.metrics;
