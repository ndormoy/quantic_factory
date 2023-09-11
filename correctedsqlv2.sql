-- Exported from QuickDBD: https://www.quickdatabasediagrams.com/
-- NOTE! If you have used non-SQL datatypes in your design, you will have to change these here.

CREATE TABLE `ChannelType` (
    -- UNSIGNED
    `ChannelTypeID` smallint UNSIGNED AUTO_INCREMENT NOT NULL ,
    `Name` varchar(30)  NOT NULL ,
    PRIMARY KEY (
        `ChannelTypeID`
    )
);

CREATE TABLE `EventType` (
    -- UNSIGNED
    `EventTypeID` smallint UNSIGNED AUTO_INCREMENT NOT NULL ,
    `Name` varchar(30)  NOT NULL ,
    PRIMARY KEY (
        `EventTypeID`
    )
);

CREATE TABLE `Content` (
    -- UNSIGNED
    `ContentID` int UNSIGNED AUTO_INCREMENT NOT NULL ,
    -- UNSIGNED
    `ClientContentID` bigint UNSIGNED NOT NULL ,
    `InsertDate` timestamp  NOT NULL ,
    PRIMARY KEY (
        `ContentID`
    )
);

CREATE TABLE `Customer` (
    -- AUTOINCREMENT UNSIGNED
    `CustomerID` bigint UNSIGNED AUTO_INCREMENT NOT NULL ,
    -- UNSIGNED
    `ClientCustomerID` bigint UNSIGNED NOT NULL ,
    `InsertDate` timestamp  NOT NULL ,
    PRIMARY KEY (
        `CustomerID`
    )
);

CREATE TABLE `CustomerData` (
    -- AUTOINCREMENT UNSIGNED
    `CustomerChannelID` bigint UNSIGNED AUTO_INCREMENT NOT NULL ,
    -- UNSIGNED
    `CustomerID` bigint UNSIGNED NOT NULL ,
    -- UNSIGNED
    `ChannelTypeID` smallint  NOT NULL ,
    `ChannelValue` varchar(600)  NOT NULL ,
    `InsertDate` timestamp  NOT NULL ,
    PRIMARY KEY (
        `CustomerChannelID`
    )
);

CREATE TABLE `CustomerEvent` (
    -- AUTOINCREMENT UNSIGNED
    `EventID` bigint UNSIGNED AUTO_INCREMENT NOT NULL ,
    -- UNSIGNED
    `ClientEventID` bigint  NOT NULL ,
    `InsertDate` timestamp  NOT NULL ,
    PRIMARY KEY (
        `EventID`
    )
);

CREATE TABLE `CustomerEventData` (
    -- AUTOINCREMENT UNSIGNED
    `EventDataId` bigint UNSIGNED AUTO_INCREMENT NOT NULL ,
    -- UNSIGNED
    `EventID` bigint UNSIGNED NOT NULL ,
    -- UNSIGNED
    `ContentID` int UNSIGNED NOT NULL ,
    -- UNSIGNED
    `CustomerID` bigint UNSIGNED NOT NULL ,
    -- UNSIGNED
    `EventTypeID` smallint UNSIGNED NOT NULL ,
    `EventDate` timestamp  NOT NULL ,
    -- UNSIGNED
    `Quantity` smallint UNSIGNED NOT NULL ,
    `InsertDate` timestamp  NOT NULL ,
    PRIMARY KEY (
        `EventDataId`
    )
);



CREATE TABLE `ContentPrice` (
    -- UNSIGNED
    `ContentPriceID` mediumint UNSIGNED AUTO_INCREMENT NOT NULL ,
    -- UNSIGNED
    `ContentID` int UNSIGNED NOT NULL ,
    `Price` decimal(8,2) UNSIGNED NOT NULL ,
    `Currency` char(3)  NOT NULL ,
    `InsertDate` timestamp  NOT NULL ,
    PRIMARY KEY (
        `ContentPriceID`
    )
);


);

ALTER TABLE `CustomerData` ADD CONSTRAINT `fk_CustomerData_CustomerID` FOREIGN KEY(`CustomerID`)
REFERENCES `Customer` (`CustomerID`);

ALTER TABLE `CustomerData` ADD CONSTRAINT `fk_CustomerData_ChannelTypeID` FOREIGN KEY(`ChannelTypeID`)
REFERENCES `ChannelType` (`ChannelTypeID`);

ALTER TABLE `CustomerEventData` ADD CONSTRAINT `fk_CustomerEventData_EventID` FOREIGN KEY(`EventID`)
REFERENCES `CustomerEvent` (`EventID`);

ALTER TABLE `CustomerEventData` ADD CONSTRAINT `fk_CustomerEventData_ContentID` FOREIGN KEY(`ContentID`)
REFERENCES `Content` (`ContentID`);

ALTER TABLE `CustomerEventData` ADD CONSTRAINT `fk_CustomerEventData_CustomerID` FOREIGN KEY(`CustomerID`)
REFERENCES `Customer` (`CustomerID`);

ALTER TABLE `CustomerEventData` ADD CONSTRAINT `fk_CustomerEventData_EventTypeID` FOREIGN KEY(`EventTypeID`)
REFERENCES `EventType` (`EventTypeID`);

ALTER TABLE `ContentPrice` ADD CONSTRAINT `fk_ContentPrice_ContentID` FOREIGN KEY(`ContentID`)
REFERENCES `Content` (`ContentID`);

