CREATE TABLE req_history (
  req_id SERIAL NOT NULL,
  url_short TEXT NOT NULL,
  req_time BIGINT NOT NULL,

  CONSTRAINT pk_req PRIMARY KEY (req_id),
  CONSTRAINT fk_url FOREIGN KEY (url_short) REFERENCES url_map (url_short)
);

CREATE INDEX ix_req_url ON req_history (url_short)
