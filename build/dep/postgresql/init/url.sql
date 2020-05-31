CREATE TABLE url_map (
  url_id SERIAL NOT NULL,
  url_short text NOT NULL,
  url_long text NOT NULL,

  CONSTRAINT pk_url PRIMARY KEY (url_id),
  CONSTRAINT uq_url_short UNIQUE (url_short)
);

CREATE INDEX ix_url_short ON url_map (url_short);
