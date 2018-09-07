CREATE TABLE IF NOT EXISTS `event` (
	`position` bigint NOT NULL AUTO_INCREMENT,
	`timestamp` DATETIME DEFAULT CURRENT_TIMESTAMP,
	`stream_category` varchar(24) NOT NULL,
	`stream_id` varchar(40) NOT NULL,
	`event_number` int NOT NULL,
	`event_type` varchar(24) NOT NULL,
	`data` json,
   primary key (`position`),
   CONSTRAINT unique_event_number UNIQUE (stream_id, event_number)
);

CREATE TABLE IF NOT EXISTS `snapshot` (
	`position` bigint NOT NULL,
	`created` DATETIME DEFAULT CURRENT_TIMESTAMP,
	`status` varchar(24) NOT NULL,
	`location` varchar(24) DEFAULT NULL,
	`status_last_updated` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   primary key (`position`)
);
