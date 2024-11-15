create database car_management_application;
use car_management_application;

drop database car_management_application;
CREATE TABLE users (
user_id INT PRIMARY KEY auto_increment,
first_name VARCHAR(100) NOT NULL,
last_name VARCHAR(100) NOT NULL,
email VARCHAR(255) NOT NULL UNIQUE,
password VARCHAR(255) NOT NULL
);
CREATE TABLE cars (
    car_id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT,
    car_name VARCHAR(255) NOT NULL,
    tags VARCHAR(255), 
    description TEXT NOT NULL,
    car_type VARCHAR(50),
    car_company VARCHAR(100),
    dealer VARCHAR(100),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);
CREATE TABLE car_images (
    image_id INT PRIMARY KEY AUTO_INCREMENT,
    car_id INT,
    image_url VARCHAR(255) NOT NULL, 
    FOREIGN KEY (car_id) REFERENCES cars(car_id)
);


select *from users;
delete from users where user_id =5;

ALTER Table users
ADD created_at datetime default current_timestamp;

ALTER Table users
ADD updated_at datetime default current_timestamp;

ALTER Table users
ADD deleted_at datetime default current_timestamp;

