-- +----------------------------------------------------------------------------
-- | Database
-- +----------------------------------------------------------------------------

CREATE DATABASE IF NOT EXISTS `test` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
USE test;

-- +----------------------------------------------------------------------------
-- | Tables
-- +----------------------------------------------------------------------------

CREATE TABLE `user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `dataSourceId` int(10) unsigned NOT NULL,
  `name` varchar(64) COLLATE utf8_bin NOT NULL,
  `added` datetime NOT NULL,
  `statusId` tinyint(3) unsigned NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `dataSourceId_name_UNIQUE` (`dataSourceId`,`name`),
  KEY `fk_user_status1_idx` (`statusId`),
  KEY `added_INDEX` (`added`),
  KEY `fk_user_dataSource1_idx` (`dataSourceId`)
) ENGINE=InnoDB;

-- +----------------------------------------------------------------------------
-- | Functions
-- +----------------------------------------------------------------------------

DELIMITER $$

CREATE DEFINER=gary@localhost FUNCTION hello_world() RETURNS text
BEGIN
    RETURN 'Hello World';
END$$

DELIMITER ;

-- +----------------------------------------------------------------------------
-- | Procedures
-- +----------------------------------------------------------------------------

DELIMITER `foo`

CREATE DEFINER=gary@localhost PROCEDURE GetAllProducts()
BEGIN
    SELECT *  FROM users;
ENDfoo

DELIMITER ;

-- +----------------------------------------------------------------------------
-- | Triggers
-- +----------------------------------------------------------------------------

DELIMITER 'hello'

CREATE DEFINER=`gary`@`localhost` TRIGGER test
AFTER INSERT ON user
FOR EACH ROW
BEGIN
    INSERT INTO user (name, added, statusId) VALUES ("trigger", NOW(), 2);
ENDhello

DELIMITER ;
