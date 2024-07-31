CREATE TABLE access_level (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255)
);

CREATE TABLE users (
  id serial PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  hashed_password VARCHAR(255) NOT NULL,
  access_level INT,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (access_level) REFERENCES access_level (id)
);

CREATE TABLE banner (
  id serial PRIMARY KEY,
  feature_id INT NOT NULL,
  content JSON NOT NULL,
  version_number INT NOT NULL DEFAULT 1,
  original_banner_id INT DEFAULT NULL,
  version_status INT NOT NULL DEFAULT 1,
  is_active BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (id, version_number)
);

CREATE TABLE banner_tag (
  banner_id INT NOT NULL,
  tag_id INT NOT NULL,
  PRIMARY KEY (banner_id, tag_id),
  FOREIGN KEY (banner_id) REFERENCES banner (id) ON DELETE CASCADE
);

insert into access_level (name) values ('admin'), ('user');
insert into users (username, hashed_password, access_level) values ('admin', '3335e880d573848bc02e9550036e29e5d428d373', 1)
