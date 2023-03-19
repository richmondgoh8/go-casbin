CREATE SCHEMA test

CREATE TABLE test.roles (
  id SERIAL PRIMARY KEY,
  role_name varchar(255) NOT NULL,
  role_description varchar(255) NOT NULL
);

INSERT INTO 
    test.roles (role_name, role_description)
VALUES
    ('Admin','Takes care of everything'),
    ('Manager','Takes care of something');

CREATE TABLE test.permissions (
  id SERIAL PRIMARY KEY,
  role_id INT,
  policy_type varchar(255) NOT NULL,
  controller varchar(255) NOT NULL,
  action_type varchar(255) NOT NULL,
  CONSTRAINT role_id_fk FOREIGN KEY (role_id) REFERENCES test.roles(id)
);

INSERT INTO
    test.permissions (role_id, policy_type, controller, action_type)
VALUES
    (1, 'p', 'user','read'),
    (1, 'p', 'user','update');