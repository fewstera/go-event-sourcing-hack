CREATE TABLE IF NOT EXISTS `event` (
	`position` bigint NOT NULL AUTO_INCREMENT,
	`timestamp` DATETIME DEFAULT CURRENT_TIMESTAMP,
	`stream_category` varchar(24) NOT NULL,
	`stream_id` varchar(40) NOT NULL,
	`event_number` int NOT NULL,
	`event_type` varchar(24) NOT NULL,
	`data` json,
   PRIMARY KEY (`position`),
   KEY `by_category_and_id` (`stream_category`, `stream_id`),
   CONSTRAINT `unique_event_number` UNIQUE (`stream_category`, `stream_id`, `event_number`)
);
