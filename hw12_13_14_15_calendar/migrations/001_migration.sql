-- +goose Up
CREATE TABLE eventsTable
(
    id INT PRIMARY KEY AUTO_INCREMENT,
    title varchar(255) NOT NULL,
    userID varchar(50) NOT NULL,
    description varchar(1500),
    dateStart datetime NOT NULL,
    dateStop datetime NOT NULL,
    eventMessageTimeDelta BIGINT
);

-- +goose Down
DROP TABLE eventsTable;