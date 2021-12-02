create database pp default charset utf8;
create user pp_user@localhost identified by 'password';
grant all on pp.* to pp_user@localhost identified by 'password';
flush privileges;

USE pp;
CREATE TABLE `applications` (
  `id` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `slug` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `company_name` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `logo_detail` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `logo_list` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `category_code` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `category_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `rate` float NOT NULL,
  `users_rated` int(11) NOT NULL,
  `users_downloaded` int(11) NOT NULL,
  `website` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `link` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `is_offline` tinyint(4) NOT NULL DEFAULT '0',
  `is_gdrive` tinyint(4) NOT NULL DEFAULT '0',
  `is_google` tinyint(4) NOT NULL DEFAULT '0',
  `is_android` tinyint(4) NOT NULL DEFAULT '0',
  `is_free` tinyint(4) NOT NULL DEFAULT '0',
  `faq_page` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `version` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `last_updated` date NOT NULL,
  `created_at` datetime DEFAULT NULL,
  `last_modified` datetime DEFAULT NULL,
  `type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `size` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `contact_email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contact_address` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `privacy_policy` varchar(1500) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `price` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `android_link` varchar(700) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `default_lang` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `count_langs` smallint(5) NOT NULL DEFAULT '0',
  `count_reviews` smallint(5) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY(`company_name`),
  KEY(`category_name`),
  KEY(`users_downloaded`),
  KEY(`rate`),
  KEY(`count_reviews`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;