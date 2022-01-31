CREATE TABLE `authentication` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `token` varchar(255) NOT NULL,
  `account_id` varchar(255) NOT NULL,
  `expiry_time` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE `member` (
  `id` varchar(255) NOT NULL,
  `merchant_id` varchar(255) NOT NULL,
  `email_address` varchar(255) NOT NULL,
  `name` varchar(100) DEFAULT NULL,
  `address` varchar(255) DEFAULT NULL,
  `password` varchar(255) NOT NULL,
  `role` varchar(20) NOT NULL DEFAULT 'user',
  `is_deleted` tinyint(4) NOT NULL DEFAULT '0',
  `created_time` bigint(20) NOT NULL DEFAULT '0',
  `updated_time` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `member_email_address_uindex` (`email_address`),
  KEY `member_merchant_id_fk` (`merchant_id`),
  KEY `member_created_time_index` (`created_time`),
  KEY `member_update_time_index` (`updated_time`),
  CONSTRAINT `member_merchant_id_fk` FOREIGN KEY (`merchant_id`) REFERENCES `merchant` (`id`)
);

CREATE TABLE `merchant` (
  `id` varchar(255) NOT NULL,
  `name` varchar(100) DEFAULT NULL,
  `address` varchar(255) DEFAULT NULL,
  `is_deleted` tinyint(4) NOT NULL DEFAULT '0',
  `created_time` bigint(20) NOT NULL DEFAULT '0',
  `updated_time` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `merchant_created_time_index` (`created_time`),
  KEY `merchant_updated_time_index` (`updated_time`)
);

INSERT INTO merchant (id, name, address, is_deleted, created_time, updated_time) VALUES ('c8b49264-93cb-4b19-bad3-81953cf5317e', 'supermerchant', 'myaddress', 0, 1643391877, 1643391877);
INSERT INTO member (id, merchant_id, email_address, name, address, password, role, is_deleted, created_time, updated_time) VALUES ('f55b6262-9ee8-4f07-8855-493b6b5cacb1', 'c8b49264-93cb-4b19-bad3-81953cf5317e', 'superadmin@merchant.com', 'myname', 'myaddress', '$2a$10$Dkw3zYsqPn5GosAH0ZZLSO3KZ1anschXrciltSV7d/lW.C56qwIWC', 'superadmin', 0, 1643433607, 0);