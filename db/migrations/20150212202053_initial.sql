
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE records (
  id integer NOT NULL,
  source character varying(255) NOT NULL,
  latest character varying(255) NOT NULL,
  modified timestamp with time zone NOT NULL
);


CREATE SEQUENCE records_id_seq
  START WITH 1
  INCREMENT BY 1
  NO MINVALUE
  NO MAXVALUE
  CACHE 1;


ALTER SEQUENCE records_id_seq OWNED BY records.id;


ALTER TABLE ONLY records ALTER COLUMN id SET DEFAULT nextval('records_id_seq'::regclass);


ALTER TABLE ONLY records
  ADD CONSTRAINT records_pkey PRIMARY KEY (id),
  ADD CONSTRAINT records_source_uniq UNIQUE (source);


CREATE INDEX records_source_idx ON records USING btree (source);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE records CASCADE;
