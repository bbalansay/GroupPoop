CREATE TABLE IF NOT EXISTS tblUser {
  ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  Email VARCHAR(128) NOT NULL,
  UserName VARCHAR(256) NOT NULL,
  PassHash VARCHAR(128) NOT NULL,
  FirstName VARCHAR(128) NOT NULL,
  LastName VARCHAR(128) NOT NULL,
  PhotoURL VARCHAR(512) NOT NULL,
  INDEX (Email, UserName)
}

CREATE TABLE IF NOT EXISTS tblBathroom {
  ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  Name VARCHAR(256) NOT NULL,
  Description VARCHAR(512) NOT NULL,
  Location VARCHAR(128) NOT NULL,
  NumSinks INT NOT NULL,
  NumToilets INT NOT NULL,
  NumUrinals INT NOT NULL,
  NumTrashCans INT NOT NULL,
  NumAirDryers INT NOT NULL,
  NumTowelDispensers INT NOT NULL
}

CREATE TABLE IF NOT EXISTS tblReview {
  ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  FOREIGN KEY (UserID) REFERENCES tblUser(ID) ON DELETE CASCADE NOT NULL,
  FOREIGN KEY (BathroomID) REFERENCES tblBathroom(ID) ON DELETE CASCADE NOT NULL,
  Content VARCHAR(512) NOT NULL,
  Time DATETIME NOT NULL
}

/*
INSERT INTO users (email, pass_hash, user_name, first_name, last_name, photo_url) VALUES ("admin@yfzhou.me", "password123", "admin", "first", "last", "photo_url");
SET @UID = LAST_INSERT_ID();
INSERT INTO channel (name, description, private, createdAt, creator, editedAt) VALUES ("general", "General channel for general discussion", FALSE, NOW(), @UID, NULL);
SET @CHID = LAST_INSERT_ID();
INSERT INTO message (channelID, body, createdAt, creator, editedAt) VALUES (@CHID, "Hello world! This is the general channel.", NOW(), @UID, NULL);
ContentContent
ALTER USER 'root' IDENTIFIED WITH mysql_native_password BY 'password123'
*/