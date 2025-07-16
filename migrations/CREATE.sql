CREATE TABLE Block
(
	Id INT NOT NULL PRIMARY KEY,
	BlockerId int not null,
	BlockeeId int not null
);

CREATE TABLE ControlAppCmd
(
	Id INT NOT NULL PRIMARY KEY,
	SenderId int not null DEFAULT 0,
	SubId int not null,
	CmdId int not null,
    GroupRefId INT NULL
);

CREATE TABLE CommandList
(
	CmdId INT NOT NULL PRIMARY KEY,
	Content TEXT not null, 
    SendDate timestamp NOT NULL DEFAULT CURRENT_DATE 
);

CREATE TABLE CmdLog
(
	Id INT NOT NULL PRIMARY KEY,
	SenderID int not null DEFAULT 0,
	SubId int not null,
	Content TEXT not null, 
    GroupRefId INT NULL
);

CREATE TABLE ChatLog
(
	Id INT NOT NULL PRIMARY KEY,
	ReceiverID int,
	SenderID int,
	ChatTxt TEXT,
	Date_Add timestamp default CURRENT_DATE
);

CREATE TABLE Invites
(
	Id INT NOT NULL PRIMARY KEY,
	SubId int not null,
	DomId int not null
);

CREATE TABLE SubContent
(
	Id INT NOT NULL PRIMARY KEY, 
    Sub_User_Id INT NULL, 
    IFrame_Content TEXT NULL, 
    Date_add timestamp default CURRENT_DATE
);

CREATE TABLE Report
(
	Id INT NOT NULL PRIMARY KEY,
	ReporterId int not null,
	Reportee int not null,
	ReportedCommand TEXT
);

CREATE TABLE Relationship
(
	Id INT NOT NULL PRIMARY KEY, 
    DomID INT NULL, 
    SubID INT NULL
	
);

CREATE TABLE Groups
(
	Id INT NOT NULL PRIMARY KEY,
	GroupName varchar(50) NOT NULL, 
    RefId INT NOT NULL
);

CREATE TABLE UsersGroups
(
	Id INT NOT NULL PRIMARY KEY,
	UserId INT NOT NULL,
	GroupRefId INT NOT NULL
);

CREATE TABLE Users
(
	Id INT NOT NULL PRIMARY KEY, 
    ScreenName VARCHAR(50) NOT NULL, 
    LoginName VARCHAR(50) NOT NULL, 
    Password VARCHAR(50) NOT NULL, 
    Role VARCHAR(50) NULL,
    RandOpt BOOLEAN , 
    AnonCmd BOOLEAN , 
    Varified BOOLEAN NOT NULL DEFAULT FALSE, 
    VarifiedCode INT NULL DEFAULT RANDOM(), 
    LoginDate timestamp NOT NULL DEFAULT CURRENT_DATE, 
    ThumbsUp INT NULL DEFAULT 0
);

CREATE TABLE UserData
(
	Id INT NOT NULL PRIMARY KEY,
	CommandId INT,
	NoViews INT,
	NoThumbsUp INT,
	NOThumbsDown INT,
	Content TEXT
);

CREATE TABLE SubReport
(
	Id INT NOT NULL PRIMARY KEY, 
    Sub_User_Id INT NULL, 
    Dom_User_Id INT NULL, 
    Event_Datetime timestamp NULL, 
    Event_Desc VARCHAR(50) NULL
);

CREATE TABLE VarificationReq
(
	Id INT NOT NULL PRIMARY KEY,
	SubID int not null,
	VarificationCode int null, 
    Count INT NULL DEFAULT 0, 
    VerifyDate TIMESTAMP NULL, 
    Sent BOOLEAN NULL DEFAULT FALSE
);