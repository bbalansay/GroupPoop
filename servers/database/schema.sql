CREATE TABLE IF NOT EXISTS tblUser (
  ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  Email VARCHAR(512) NOT NULL UNIQUE,
  UserName VARCHAR(256) NOT NULL UNIQUE,
  PassHash VARCHAR(128) NOT NULL,
  FirstName VARCHAR(128) NOT NULL,
  LastName VARCHAR(128) NOT NULL,
  PhotoURL VARCHAR(512) NOT NULL,
  INDEX (Email, UserName)
);

CREATE TABLE IF NOT EXISTS tblBathroom (
  ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  Name VARCHAR(256) NOT NULL,
  Description VARCHAR(512) NOT NULL,
  Location VARCHAR(128) NOT NULL,
  Gender VARCHAR(128) NOT NULL,
  NumSinks INT NOT NULL,
  NumToilets INT NOT NULL,
  NumUrinals INT NOT NULL,
  NumTrashCans INT NOT NULL,
  NumAirDryers INT NOT NULL,
  NumTowelDispensers INT NOT NULL
);

CREATE TABLE IF NOT EXISTS tblReview (
  ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  UserID INT NOT NULL,
  FOREIGN KEY (UserID) REFERENCES tblUser(ID) ON DELETE CASCADE,
  BathroomID INT NOT NULL,
  FOREIGN KEY (BathroomID) REFERENCES tblBathroom(ID) ON DELETE CASCADE,
  Score INT NOT NULL,
  Content VARCHAR(512) NOT NULL,
  CreatedAt DATETIME NOT NULL,
  EditedAt DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS tblFavorites (
  ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  UserID INT NOT NULL,
  FOREIGN KEY (UserID) REFERENCES tblUser(ID) ON DELETE CASCADE,
  BathroomID INT NOT NULL,
  FOREIGN KEY (BathroomID) REFERENCES tblBathroom(ID) ON DELETE CASCADE
);

INSERT INTO tblBathroom (Name, Description, Location, Gender, NumSinks, NumToilets, NumUrinals, NumTrashCans, NumAirDryers, NumTowelDispensers)
VALUES ("Men's Mary Gates 4th Floor", "The Mecca of bathrooms, big and not frequented", "Mary Gates Hall", "Masculine", 4, 2, 5, 2, 2, 2);
INSERT INTO tblBathroom (Name, Description, Location, Gender, NumSinks, NumToilets, NumUrinals, NumTrashCans, NumAirDryers, NumTowelDispensers)
VALUES ("Women's Mary Gates 4th Floor", "The Mecca of bathrooms, big and not frequented", "Mary Gates Hall", "Feminine", 4, 4, 0, 2, 2, 2);
INSERT INTO tblBathroom (Name, Description, Location, Gender, NumSinks, NumToilets, NumUrinals, NumTrashCans, NumAirDryers, NumTowelDispensers)
VALUES ("Men's Ode 1st Floor", "MY NOSTRILS!!!", "Odegaard Library", "Masculine", 5, 4, 5, 3, 2, 2);
INSERT INTO tblBathroom (Name, Description, Location, Gender, NumSinks, NumToilets, NumUrinals, NumTrashCans, NumAirDryers, NumTowelDispensers)
VALUES ("Women's Mary Gates 4th Floor", "MY NOSTRILS!!!", "Odegaard Library", "Feminine", 5, 6, 0, 3, 2, 2);
INSERT INTO tblBathroom (Name, Description, Location, Gender, NumSinks, NumToilets, NumUrinals, NumTrashCans, NumAirDryers, NumTowelDispensers)
VALUES ("Gender Neutral Smith 2nd Floor", "The door locks but there's two stalls...", "Smith Hall", "Gender Neutral", 1, 2, 0, 1, 1, 1);
INSERT INTO tblBathroom (Name, Description, Location, Gender, NumSinks, NumToilets, NumUrinals, NumTrashCans, NumAirDryers, NumTowelDispensers)
VALUES ("Gender Neutral Gowen 2nd Floor", "Never actually been in there, always been occupied.", "Gowen Hall", "Gender Neutral", 1, 1, 1, 1, 1, 1);

-- Set password for nodejs app-- not an ideal way.
alter user root identified with mysql_native_password by 'password123';
flush privileges;




