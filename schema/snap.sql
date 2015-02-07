SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';

DROP SCHEMA IF EXISTS `snap_config` ;
CREATE SCHEMA IF NOT EXISTS `snap_config` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci ;
USE `snap_config` ;

-- -----------------------------------------------------
-- Table `snap_config`.`initialisedDatabases`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `snap_config`.`initialisedDatabases` ;

CREATE TABLE IF NOT EXISTS `snap_config`.`initialisedDatabases` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL,
  `initialized` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `databaseName_UNIQUE` (`name` ASC))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `snap_config`.`revisions`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `snap_config`.`revisions` ;

CREATE TABLE IF NOT EXISTS `snap_config`.`revisions` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `databaseId` INT UNSIGNED NOT NULL,
  `revision` INT UNSIGNED NOT NULL DEFAULT 1,
  `upSql` TEXT NOT NULL,
  `downSql` TEXT NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_revisions_initialisedDatabases_idx` (`databaseId` ASC),
  CONSTRAINT `fk_revisions_initialisedDatabases`
    FOREIGN KEY (`databaseId`)
    REFERENCES `snap_config`.`initialisedDatabases` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
