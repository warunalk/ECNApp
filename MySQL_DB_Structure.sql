CREATE DATABASE `ecnapplication` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;

CREATE TABLE `access_level` (
  `access_level_id` int NOT NULL,
  `access_level` int DEFAULT NULL,
  `access_level_desc` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`access_level_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='levels of access';


CREATE TABLE `user` (
  `user_id` int unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(100) NOT NULL,
  `email` varchar(100) NOT NULL,
  `fname` varchar(100) NOT NULL,
  `lname` varchar(100) NOT NULL,
  `password` varchar(100) NOT NULL,
  `createddate` datetime NOT NULL,
  `lastmoddate` datetime NOT NULL,
  `lastloggeddate` datetime NOT NULL,
  `isadmin` tinyint NOT NULL DEFAULT '0',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=latin1;


CREATE TABLE `user_access_level` (
  `user_access_level_id` int NOT NULL,
  `user_id` int NOT NULL,
  `access_level_id` int NOT NULL,
  KEY `fk_access_level` (`user_access_level_id`),
  CONSTRAINT `fk_access_level` FOREIGN KEY (`user_access_level_id`) REFERENCES `access_level` (`access_level_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='user access level config';







