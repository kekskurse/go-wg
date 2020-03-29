CREATE DATABASE wg;
USE wg;
CREATE TABLE tickets (
    id int NOT NULL AUTO_INCREMENT,
    ticketID varchar(255) NOT NULL,
    status ENUM('new', 'approved', 'denied', 'rejected'),
    publicKey varchar(100) NOT NULL,
    publicIP varchar(250) NOT NULL,
    hostname varchar(250),
    internIpv4 varchar(30),

    PRIMARY KEY (id),
    UNIQUE(ticketID),
    UNIQUE(publicKey)
);


