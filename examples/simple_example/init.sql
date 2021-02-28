CREATE USER 'example_user'@'localhost' IDENTIFIED BY 'example_password';
CREATE DATABASE example;
GRANT ALL PRIVILEGES ON example.* TO 'example_user'@'localhost';
GRANT FILE ON *.* TO 'example_user'@'localhost';
USE example;

CREATE TABLE user (
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name varchar(255)
);

INSERT INTO user(name) VALUES ('john'),('harry'),('rite'),('Sophie');

CREATE TABLE secret (
    secret varchar(255)
);

INSERT INTO secret(secret) VALUES ('th!5!sSup3rS3cr3T!');