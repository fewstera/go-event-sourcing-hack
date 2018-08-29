CREATE TABLE `event` (
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
