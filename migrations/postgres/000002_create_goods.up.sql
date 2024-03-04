CREATE TABLE IF NOT EXISTS goods (
    id SERIAL NOT NULL,
    project_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(1000) NULL,
    priority INT NOT NULL,
    removed BOOL NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY(id, project_id)
);


CREATE INDEX IF NOT EXISTS idx_goods_name ON goods (name);