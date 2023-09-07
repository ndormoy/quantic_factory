-- Exported from QuickDBD: https://www.quickdatabasediagrams.com/
-- NOTE! If you have used non-SQL datatypes in your design, you will have to change these here.


CREATE TABLE `Customer` (
    -- AUTOINCREMENT UNSIGNED
    `CustomerID` bigint  NOT NULL ,
    -- UNSIGNED
    `ClientCustomerID` bigint  NOT NULL ,
    `InsertDate` timestamp  NOT NULL ,
    PRIMARY KEY (
        `CustomerID`
    )
);

CREATE TABLE `CustomerData` (
    -- AUTOINCREMENT UNSIGNED
    `CustomerChannelID` bigint  NOT NULL ,
    -- UNSIGNED
    `CustomerID` bigint  NOT NULL ,
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
    `EventID` bigint  NOT NULL ,
    -- UNSIGNED
    `ClientEventID` bigint  NOT NULL ,
    `InsertDate` timestamp  NOT NULL ,
    PRIMARY KEY (
        `EventID`
    )
);

CREATE TABLE `CustomerEventData` (
    -- AUTOINCREMENT UNSIGNED
    `EventDataId` bigint  NOT NULL ,
    -- UNSIGNED
    `EventID` bigint  NOT NULL ,
    -- UNSIGNED
    `ContentID` int  NOT NULL ,
    -- UNSIGNED
    `CustomerID` bigint  NOT NULL ,
    -- UNSIGNED
    `EventTypeID` smallint  NOT NULL ,
    `EventDate` timestamp  NOT NULL ,
    -- UNSIGNED
    `Quantity` smallint  NOT NULL ,
    `InsertDate` timestamp  NOT NULL ,
    PRIMARY KEY (
        `EventDataId`
    )
);

CREATE TABLE `Content` (
    -- UNSIGNED
    `ContentID` int AUTO_INCREMENT NOT NULL ,
    -- UNSIGNED
    `ClientContentID` bigint  NOT NULL ,
    `InsertDate` timestamp  NOT NULL ,
    PRIMARY KEY (
        `ContentID`
    )
);

CREATE TABLE `ContentPrice` (
    -- UNSIGNED
    `ContentPriceID` mediumint AUTO_INCREMENT NOT NULL ,
    -- UNSIGNED
    `ContentID` int  NOT NULL ,
    `Price` decimal(8,2)  NOT NULL ,
    `Currency` char(3)  NOT NULL ,
    `InsertDate` timestamp  NOT NULL ,
    PRIMARY KEY (
        `ContentPriceID`
    )
);

CREATE TABLE `ChannelType` (
    -- UNSIGNED
    `ChannelTypeID` smallint AUTO_INCREMENT NOT NULL ,
    `Name` varchar(30)  NOT NULL ,
    PRIMARY KEY (
        `ChannelTypeID`
    )
);

CREATE TABLE `EventType` (
    -- UNSIGNED
    `EventTypeID` smallint AUTO_INCREMENT NOT NULL ,
    `Name` varchar(30)  NOT NULL ,
    PRIMARY KEY (
        `EventTypeID`
    )
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

