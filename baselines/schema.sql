
/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
DROP TABLE IF EXISTS `Vote`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Vote` (
  `vote_id` int NOT NULL AUTO_INCREMENT,
  `vote_actor` bigint unsigned NOT NULL,
  `username` varbinary(255) NOT NULL DEFAULT '0',
  `vote_user_id` int NOT NULL DEFAULT '0',
  `vote_page_id` int NOT NULL DEFAULT '0',
  `vote_value` char(1) CHARACTER SET latin1 COLLATE latin1_bin NOT NULL DEFAULT '',
  `vote_date` datetime NOT NULL DEFAULT '1970-01-01 00:00:01',
  `vote_ip` varbinary(45) NOT NULL DEFAULT '',
  PRIMARY KEY (`vote_id`),
  KEY `vote_page_id_index` (`vote_page_id`),
  KEY `valueidx` (`vote_value`),
  KEY `actoridx` (`vote_actor`),
  KEY `usernameidx` (`username`),
  KEY `vote_date` (`vote_date`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `actor`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `actor` (
  `actor_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `actor_user` int unsigned DEFAULT NULL,
  `actor_name` varbinary(255) NOT NULL,
  PRIMARY KEY (`actor_id`),
  UNIQUE KEY `actor_name` (`actor_name`),
  UNIQUE KEY `actor_user` (`actor_user`)
) ENGINE=InnoDB AUTO_INCREMENT=44 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `ajaxpoll_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ajaxpoll_info` (
  `poll_id` varbinary(32) NOT NULL DEFAULT '',
  `poll_txt` blob,
  `poll_show_results_before_voting` tinyint(1) DEFAULT NULL,
  `poll_date` datetime DEFAULT NULL,
  PRIMARY KEY (`poll_id`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `ajaxpoll_vote`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ajaxpoll_vote` (
  `poll_vote_id` int unsigned NOT NULL AUTO_INCREMENT,
  `poll_id` varbinary(32) NOT NULL DEFAULT '',
  `poll_actor` bigint unsigned NOT NULL,
  `poll_ip` varbinary(255) DEFAULT NULL,
  `poll_answer` int DEFAULT NULL,
  `poll_date` datetime DEFAULT NULL,
  PRIMARY KEY (`poll_vote_id`),
  UNIQUE KEY `poll_id_actor` (`poll_id`,`poll_actor`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `archive`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `archive` (
  `ar_id` int unsigned NOT NULL AUTO_INCREMENT,
  `ar_namespace` int NOT NULL DEFAULT '0',
  `ar_title` varbinary(255) NOT NULL DEFAULT '',
  `ar_comment_id` bigint unsigned NOT NULL,
  `ar_actor` bigint unsigned NOT NULL,
  `ar_timestamp` binary(14) NOT NULL,
  `ar_minor_edit` tinyint NOT NULL DEFAULT '0',
  `ar_rev_id` int unsigned NOT NULL,
  `ar_deleted` tinyint unsigned NOT NULL DEFAULT '0',
  `ar_len` int unsigned DEFAULT NULL,
  `ar_page_id` int unsigned DEFAULT NULL,
  `ar_parent_id` int unsigned DEFAULT NULL,
  `ar_sha1` varbinary(32) NOT NULL DEFAULT '',
  PRIMARY KEY (`ar_id`),
  UNIQUE KEY `ar_revid_uniq` (`ar_rev_id`),
  KEY `ar_name_title_timestamp` (`ar_namespace`,`ar_title`,`ar_timestamp`),
  KEY `ar_actor_timestamp` (`ar_actor`,`ar_timestamp`)
) ENGINE=InnoDB AUTO_INCREMENT=853 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `block`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `block` (
  `bl_id` int unsigned NOT NULL AUTO_INCREMENT,
  `bl_target` int unsigned NOT NULL,
  `bl_by_actor` bigint unsigned NOT NULL,
  `bl_reason_id` bigint unsigned NOT NULL,
  `bl_timestamp` binary(14) NOT NULL,
  `bl_anon_only` tinyint(1) NOT NULL DEFAULT '0',
  `bl_create_account` tinyint(1) NOT NULL DEFAULT '1',
  `bl_enable_autoblock` tinyint(1) NOT NULL DEFAULT '1',
  `bl_expiry` varbinary(14) NOT NULL,
  `bl_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `bl_block_email` tinyint(1) NOT NULL DEFAULT '0',
  `bl_allow_usertalk` tinyint(1) NOT NULL DEFAULT '0',
  `bl_parent_block_id` int unsigned DEFAULT NULL,
  `bl_sitewide` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`bl_id`),
  KEY `bl_timestamp` (`bl_timestamp`),
  KEY `bl_target` (`bl_target`),
  KEY `bl_expiry` (`bl_expiry`),
  KEY `bl_parent_block_id` (`bl_parent_block_id`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `block_target`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `block_target` (
  `bt_id` int unsigned NOT NULL AUTO_INCREMENT,
  `bt_address` tinyblob,
  `bt_user` int unsigned DEFAULT NULL,
  `bt_user_text` varbinary(255) DEFAULT NULL,
  `bt_auto` tinyint(1) NOT NULL DEFAULT '0',
  `bt_range_start` tinyblob,
  `bt_range_end` tinyblob,
  `bt_ip_hex` tinyblob,
  `bt_count` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`bt_id`),
  KEY `bt_address` (`bt_address`(42)),
  KEY `bt_ip_user_text` (`bt_ip_hex`(35),`bt_user_text`),
  KEY `bt_range` (`bt_range_start`(35),`bt_range_end`(35)),
  KEY `bt_user` (`bt_user`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `bookrenderingtask`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `bookrenderingtask` (
  `brt_id` int unsigned NOT NULL AUTO_INCREMENT,
  `brt_timestamp` varbinary(14) NOT NULL,
  `brt_state` varbinary(32) NOT NULL,
  `brt_disposition` varbinary(255) DEFAULT NULL,
  `brt_stash_key` varbinary(255) DEFAULT NULL,
  PRIMARY KEY (`brt_id`)
) ENGINE=InnoDB AUTO_INCREMENT=518 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `bot_passwords`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `bot_passwords` (
  `bp_user` int unsigned NOT NULL,
  `bp_app_id` varbinary(32) NOT NULL,
  `bp_password` tinyblob NOT NULL,
  `bp_token` binary(32) NOT NULL DEFAULT '\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',
  `bp_restrictions` blob NOT NULL,
  `bp_grants` blob NOT NULL,
  PRIMARY KEY (`bp_user`,`bp_app_id`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `category` (
  `cat_id` int unsigned NOT NULL AUTO_INCREMENT,
  `cat_title` varbinary(255) NOT NULL,
  `cat_pages` int NOT NULL DEFAULT '0',
  `cat_subcats` int NOT NULL DEFAULT '0',
  `cat_files` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`cat_id`),
  UNIQUE KEY `cat_title` (`cat_title`),
  KEY `cat_pages` (`cat_pages`)
) ENGINE=InnoDB AUTO_INCREMENT=315 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `categorylinks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `categorylinks` (
  `cl_from` int unsigned NOT NULL DEFAULT '0',
  `cl_to` varbinary(255) NOT NULL DEFAULT '',
  `cl_sortkey` varbinary(230) NOT NULL DEFAULT '',
  `cl_sortkey_prefix` varbinary(255) NOT NULL DEFAULT '',
  `cl_timestamp` timestamp NOT NULL,
  `cl_collation` varbinary(32) NOT NULL DEFAULT '',
  `cl_type` enum('page','subcat','file') NOT NULL DEFAULT 'page',
  PRIMARY KEY (`cl_from`,`cl_to`),
  KEY `cl_sortkey` (`cl_to`,`cl_type`,`cl_sortkey`,`cl_from`),
  KEY `cl_timestamp` (`cl_to`,`cl_timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `change_tag`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `change_tag` (
  `ct_id` int unsigned NOT NULL AUTO_INCREMENT,
  `ct_rc_id` bigint unsigned DEFAULT NULL,
  `ct_log_id` int unsigned DEFAULT NULL,
  `ct_rev_id` int unsigned DEFAULT NULL,
  `ct_params` blob,
  `ct_tag_id` int unsigned NOT NULL,
  PRIMARY KEY (`ct_id`),
  UNIQUE KEY `ct_rc_tag_id` (`ct_rc_id`,`ct_tag_id`),
  UNIQUE KEY `ct_log_tag_id` (`ct_log_id`,`ct_tag_id`),
  UNIQUE KEY `ct_rev_tag_id` (`ct_rev_id`,`ct_tag_id`),
  KEY `ct_tag_id_id` (`ct_tag_id`,`ct_rc_id`,`ct_rev_id`,`ct_log_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6068 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `change_tag_def`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `change_tag_def` (
  `ctd_id` int unsigned NOT NULL AUTO_INCREMENT,
  `ctd_name` varbinary(255) NOT NULL,
  `ctd_user_defined` tinyint(1) NOT NULL,
  `ctd_count` bigint unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`ctd_id`),
  UNIQUE KEY `ctd_name` (`ctd_name`),
  KEY `ctd_count` (`ctd_count`),
  KEY `ctd_user_defined` (`ctd_user_defined`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `comment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `comment` (
  `comment_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `comment_hash` int NOT NULL,
  `comment_text` blob NOT NULL,
  `comment_data` blob,
  PRIMARY KEY (`comment_id`),
  KEY `comment_hash` (`comment_hash`)
) ENGINE=InnoDB AUTO_INCREMENT=3758 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `content`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `content` (
  `content_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `content_size` int unsigned NOT NULL,
  `content_sha1` varbinary(32) NOT NULL,
  `content_model` smallint unsigned NOT NULL,
  `content_address` varbinary(255) NOT NULL,
  PRIMARY KEY (`content_id`)
) ENGINE=InnoDB AUTO_INCREMENT=20013 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `content_models`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `content_models` (
  `model_id` int NOT NULL AUTO_INCREMENT,
  `model_name` varbinary(64) NOT NULL,
  PRIMARY KEY (`model_id`),
  UNIQUE KEY `model_name` (`model_name`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `discussiontools_item_ids`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `discussiontools_item_ids` (
  `itid_id` int unsigned NOT NULL AUTO_INCREMENT,
  `itid_itemid` varbinary(255) NOT NULL,
  PRIMARY KEY (`itid_id`),
  UNIQUE KEY `itid_itemid` (`itid_itemid`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `discussiontools_item_pages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `discussiontools_item_pages` (
  `itp_id` int unsigned NOT NULL AUTO_INCREMENT,
  `itp_items_id` int unsigned NOT NULL,
  `itp_page_id` bigint unsigned NOT NULL,
  `itp_oldest_revision_id` bigint unsigned NOT NULL,
  `itp_newest_revision_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`itp_id`),
  UNIQUE KEY `itp_items_id_page_id` (`itp_items_id`,`itp_page_id`),
  UNIQUE KEY `itp_items_id_newest_revision_id` (`itp_items_id`,`itp_newest_revision_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `discussiontools_item_revisions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `discussiontools_item_revisions` (
  `itr_id` int unsigned NOT NULL AUTO_INCREMENT,
  `itr_itemid_id` int NOT NULL,
  `itr_revision_id` bigint unsigned NOT NULL,
  `itr_items_id` int NOT NULL,
  `itr_parent_id` int DEFAULT NULL,
  `itr_transcludedfrom` bigint unsigned DEFAULT NULL,
  `itr_level` tinyint NOT NULL,
  `itr_headinglevel` tinyint DEFAULT NULL,
  PRIMARY KEY (`itr_id`),
  UNIQUE KEY `itr_itemid_id_revision_id` (`itr_itemid_id`,`itr_revision_id`),
  KEY `itr_revision_id` (`itr_revision_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `discussiontools_items`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `discussiontools_items` (
  `it_id` int unsigned NOT NULL AUTO_INCREMENT,
  `it_itemname` varbinary(255) NOT NULL,
  `it_timestamp` binary(14) DEFAULT NULL,
  `it_actor` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`it_id`),
  UNIQUE KEY `it_itemname` (`it_itemname`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `discussiontools_subscription`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `discussiontools_subscription` (
  `sub_id` int unsigned NOT NULL AUTO_INCREMENT,
  `sub_item` varbinary(255) NOT NULL,
  `sub_namespace` int NOT NULL DEFAULT '0',
  `sub_title` varbinary(255) NOT NULL,
  `sub_section` varbinary(255) NOT NULL,
  `sub_state` int NOT NULL DEFAULT '1',
  `sub_user` int unsigned NOT NULL,
  `sub_created` binary(14) NOT NULL,
  `sub_notified` binary(14) DEFAULT NULL,
  PRIMARY KEY (`sub_id`),
  UNIQUE KEY `discussiontools_subscription_itemuser` (`sub_item`,`sub_user`),
  KEY `discussiontools_subscription_user` (`sub_user`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `dpl_clview`;
/*!50001 DROP VIEW IF EXISTS `dpl_clview`*/;
SET @saved_cs_client     = @@character_set_client;
/*!50503 SET character_set_client = utf8mb4 */;
/*!50001 CREATE VIEW `dpl_clview` AS SELECT 
 1 AS `cl_to`,
 1 AS `cl_from`,
 1 AS `cl_sortkey`*/;
SET character_set_client = @saved_cs_client;
DROP TABLE IF EXISTS `echo_email_batch`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `echo_email_batch` (
  `eeb_id` int unsigned NOT NULL AUTO_INCREMENT,
  `eeb_user_id` int unsigned NOT NULL,
  `eeb_event_priority` tinyint unsigned NOT NULL DEFAULT '10',
  `eeb_event_id` int unsigned NOT NULL,
  `eeb_event_hash` varbinary(32) NOT NULL,
  PRIMARY KEY (`eeb_id`),
  UNIQUE KEY `echo_email_batch_user_event` (`eeb_user_id`,`eeb_event_id`),
  KEY `echo_email_batch_user_hash_priority` (`eeb_user_id`,`eeb_event_hash`,`eeb_event_priority`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `echo_event`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `echo_event` (
  `event_id` int unsigned NOT NULL AUTO_INCREMENT,
  `event_type` varbinary(64) NOT NULL,
  `event_variant` varbinary(64) DEFAULT NULL,
  `event_agent_id` int unsigned DEFAULT NULL,
  `event_agent_ip` varbinary(39) DEFAULT NULL,
  `event_extra` blob,
  `event_page_id` int unsigned DEFAULT NULL,
  `event_deleted` tinyint unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`event_id`),
  KEY `echo_event_type` (`event_type`),
  KEY `echo_event_page_id` (`event_page_id`)
) ENGINE=InnoDB AUTO_INCREMENT=11696 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `echo_notification`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `echo_notification` (
  `notification_event` int unsigned NOT NULL,
  `notification_user` int unsigned NOT NULL,
  `notification_timestamp` binary(14) NOT NULL,
  `notification_read_timestamp` binary(14) DEFAULT NULL,
  `notification_bundle_hash` varbinary(32) NOT NULL,
  PRIMARY KEY (`notification_user`,`notification_event`),
  KEY `echo_user_timestamp` (`notification_user`,`notification_timestamp`),
  KEY `echo_notification_event` (`notification_event`),
  KEY `echo_notification_user_read_timestamp` (`notification_user`,`notification_read_timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `echo_push_provider`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `echo_push_provider` (
  `epp_id` tinyint unsigned NOT NULL AUTO_INCREMENT,
  `epp_name` tinyblob NOT NULL,
  PRIMARY KEY (`epp_id`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `echo_push_subscription`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `echo_push_subscription` (
  `eps_id` int unsigned NOT NULL AUTO_INCREMENT,
  `eps_user` int unsigned NOT NULL,
  `eps_token` blob NOT NULL,
  `eps_token_sha256` binary(64) NOT NULL,
  `eps_provider` tinyint unsigned NOT NULL,
  `eps_updated` timestamp NOT NULL,
  `eps_data` blob,
  `eps_topic` tinyint unsigned DEFAULT NULL,
  PRIMARY KEY (`eps_id`),
  UNIQUE KEY `eps_token_sha256` (`eps_token_sha256`),
  KEY `eps_provider` (`eps_provider`),
  KEY `eps_topic` (`eps_topic`),
  KEY `eps_user` (`eps_user`),
  KEY `eps_token` (`eps_token`(10))
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `echo_push_topic`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `echo_push_topic` (
  `ept_id` tinyint unsigned NOT NULL AUTO_INCREMENT,
  `ept_text` tinyblob NOT NULL,
  PRIMARY KEY (`ept_id`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `echo_target_page`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `echo_target_page` (
  `etp_id` int unsigned NOT NULL AUTO_INCREMENT,
  `etp_page` int unsigned NOT NULL DEFAULT '0',
  `etp_event` int unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`etp_id`),
  KEY `echo_target_page_event` (`etp_event`),
  KEY `echo_target_page_page_event` (`etp_page`,`etp_event`)
) ENGINE=InnoDB AUTO_INCREMENT=7087 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `echo_unread_wikis`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `echo_unread_wikis` (
  `euw_id` int unsigned NOT NULL AUTO_INCREMENT,
  `euw_user` int unsigned NOT NULL,
  `euw_wiki` varbinary(64) NOT NULL,
  `euw_alerts` int unsigned NOT NULL,
  `euw_alerts_ts` binary(14) NOT NULL,
  `euw_messages` int unsigned NOT NULL,
  `euw_messages_ts` binary(14) NOT NULL,
  PRIMARY KEY (`euw_id`),
  UNIQUE KEY `echo_unread_wikis_user_wiki` (`euw_user`,`euw_wiki`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `ed_url_cache`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ed_url_cache` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `url` varbinary(255) NOT NULL,
  `req_time` int NOT NULL,
  `result` longblob,
  UNIQUE KEY `id` (`id`),
  KEY `url` (`url`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `ext_meta`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ext_meta` (
  `pageid` int NOT NULL,
  `rindex` tinyint(1) NOT NULL,
  `rfollow` tinyint(1) NOT NULL,
  `titlealias` varbinary(255) DEFAULT NULL,
  `keywords` blob,
  `description` blob,
  PRIMARY KEY (`pageid`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `externallinks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `externallinks` (
  `el_id` int unsigned NOT NULL AUTO_INCREMENT,
  `el_from` int unsigned NOT NULL DEFAULT '0',
  `el_to_domain_index` varbinary(255) NOT NULL DEFAULT '',
  `el_to_path` blob,
  PRIMARY KEY (`el_id`),
  KEY `el_to_domain_index_to_path` (`el_to_domain_index`,`el_to_path`(60)),
  KEY `el_from` (`el_from`)
) ENGINE=InnoDB AUTO_INCREMENT=46497 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `filearchive`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `filearchive` (
  `fa_id` int unsigned NOT NULL AUTO_INCREMENT,
  `fa_name` varbinary(255) NOT NULL DEFAULT '',
  `fa_archive_name` varbinary(255) DEFAULT '',
  `fa_storage_group` varbinary(16) DEFAULT NULL,
  `fa_storage_key` varbinary(64) DEFAULT '',
  `fa_deleted_user` int DEFAULT NULL,
  `fa_deleted_timestamp` binary(14) DEFAULT NULL,
  `fa_deleted_reason_id` bigint unsigned NOT NULL,
  `fa_size` bigint unsigned DEFAULT '0',
  `fa_width` int DEFAULT '0',
  `fa_height` int DEFAULT '0',
  `fa_metadata` mediumblob,
  `fa_bits` int DEFAULT '0',
  `fa_media_type` enum('UNKNOWN','BITMAP','DRAWING','AUDIO','VIDEO','MULTIMEDIA','OFFICE','TEXT','EXECUTABLE','ARCHIVE','3D') DEFAULT NULL,
  `fa_major_mime` enum('unknown','application','audio','image','text','video','message','model','multipart','chemical') DEFAULT 'unknown',
  `fa_minor_mime` varbinary(100) DEFAULT 'unknown',
  `fa_description_id` bigint unsigned NOT NULL,
  `fa_actor` bigint unsigned NOT NULL,
  `fa_timestamp` binary(14) DEFAULT NULL,
  `fa_deleted` tinyint unsigned NOT NULL DEFAULT '0',
  `fa_sha1` varbinary(32) NOT NULL DEFAULT '',
  PRIMARY KEY (`fa_id`),
  KEY `fa_name` (`fa_name`,`fa_timestamp`),
  KEY `fa_storage_group` (`fa_storage_group`,`fa_storage_key`),
  KEY `fa_deleted_timestamp` (`fa_deleted_timestamp`),
  KEY `fa_actor_timestamp` (`fa_actor`,`fa_timestamp`),
  KEY `fa_sha1` (`fa_sha1`(10))
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `image`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `image` (
  `img_name` varbinary(255) NOT NULL DEFAULT '',
  `img_size` bigint unsigned NOT NULL DEFAULT '0',
  `img_width` int NOT NULL DEFAULT '0',
  `img_height` int NOT NULL DEFAULT '0',
  `img_metadata` mediumblob NOT NULL,
  `img_bits` int NOT NULL DEFAULT '0',
  `img_media_type` enum('UNKNOWN','BITMAP','DRAWING','AUDIO','VIDEO','MULTIMEDIA','OFFICE','TEXT','EXECUTABLE','ARCHIVE','3D') DEFAULT NULL,
  `img_major_mime` enum('unknown','application','audio','image','text','video','message','model','multipart','chemical') NOT NULL DEFAULT 'unknown',
  `img_minor_mime` varbinary(100) NOT NULL DEFAULT 'unknown',
  `img_description_id` bigint unsigned NOT NULL,
  `img_actor` bigint unsigned NOT NULL,
  `img_timestamp` binary(14) NOT NULL,
  `img_sha1` varbinary(32) NOT NULL DEFAULT '',
  PRIMARY KEY (`img_name`),
  KEY `img_actor_timestamp` (`img_actor`,`img_timestamp`),
  KEY `img_size` (`img_size`),
  KEY `img_timestamp` (`img_timestamp`),
  KEY `img_sha1` (`img_sha1`(10)),
  KEY `img_media_mime` (`img_media_type`,`img_major_mime`,`img_minor_mime`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `imagelinks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `imagelinks` (
  `il_from` int unsigned NOT NULL DEFAULT '0',
  `il_to` varbinary(255) NOT NULL DEFAULT '',
  `il_from_namespace` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`il_from`,`il_to`),
  KEY `il_to` (`il_to`,`il_from`),
  KEY `il_backlinks_namespace` (`il_from_namespace`,`il_to`,`il_from`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `interwiki`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `interwiki` (
  `iw_prefix` varbinary(32) NOT NULL,
  `iw_url` blob NOT NULL,
  `iw_api` blob NOT NULL,
  `iw_wikiid` varbinary(64) NOT NULL,
  `iw_local` tinyint(1) NOT NULL,
  `iw_trans` tinyint NOT NULL DEFAULT '0',
  PRIMARY KEY (`iw_prefix`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `ip_changes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ip_changes` (
  `ipc_rev_id` int unsigned NOT NULL DEFAULT '0',
  `ipc_rev_timestamp` binary(14) NOT NULL,
  `ipc_hex` varbinary(35) NOT NULL DEFAULT '',
  PRIMARY KEY (`ipc_rev_id`),
  KEY `ipc_rev_timestamp` (`ipc_rev_timestamp`),
  KEY `ipc_hex_time` (`ipc_hex`,`ipc_rev_timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `ipblocks_restrictions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ipblocks_restrictions` (
  `ir_ipb_id` int unsigned NOT NULL,
  `ir_type` tinyint NOT NULL,
  `ir_value` int unsigned NOT NULL,
  PRIMARY KEY (`ir_ipb_id`,`ir_type`,`ir_value`),
  KEY `ir_type_value` (`ir_type`,`ir_value`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `iwlinks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `iwlinks` (
  `iwl_from` int unsigned NOT NULL DEFAULT '0',
  `iwl_prefix` varbinary(32) NOT NULL DEFAULT '',
  `iwl_title` varbinary(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`iwl_from`,`iwl_prefix`,`iwl_title`),
  KEY `iwl_prefix_title_from` (`iwl_prefix`,`iwl_title`,`iwl_from`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `job`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `job` (
  `job_id` int unsigned NOT NULL AUTO_INCREMENT,
  `job_cmd` varbinary(60) NOT NULL DEFAULT '',
  `job_namespace` int NOT NULL,
  `job_title` varbinary(255) NOT NULL,
  `job_timestamp` binary(14) DEFAULT NULL,
  `job_params` mediumblob NOT NULL,
  `job_random` int unsigned NOT NULL DEFAULT '0',
  `job_attempts` int unsigned NOT NULL DEFAULT '0',
  `job_token` varbinary(32) NOT NULL DEFAULT '',
  `job_token_timestamp` binary(14) DEFAULT NULL,
  `job_sha1` varbinary(32) NOT NULL DEFAULT '',
  PRIMARY KEY (`job_id`),
  KEY `job_sha1` (`job_sha1`),
  KEY `job_cmd_token` (`job_cmd`,`job_token`,`job_random`),
  KEY `job_cmd_token_id` (`job_cmd`,`job_token`,`job_id`),
  KEY `job_cmd` (`job_cmd`,`job_namespace`,`job_title`,`job_params`(128)),
  KEY `job_timestamp` (`job_timestamp`)
) ENGINE=InnoDB AUTO_INCREMENT=315765 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `l10n_cache`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `l10n_cache` (
  `lc_lang` varbinary(35) NOT NULL,
  `lc_key` varbinary(255) NOT NULL,
  `lc_value` mediumblob NOT NULL,
  PRIMARY KEY (`lc_lang`,`lc_key`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `langlinks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `langlinks` (
  `ll_from` int unsigned NOT NULL DEFAULT '0',
  `ll_lang` varbinary(35) NOT NULL DEFAULT '',
  `ll_title` varbinary(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`ll_from`,`ll_lang`),
  KEY `ll_lang` (`ll_lang`,`ll_title`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `linktarget`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `linktarget` (
  `lt_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `lt_namespace` int NOT NULL,
  `lt_title` varbinary(255) NOT NULL,
  PRIMARY KEY (`lt_id`),
  UNIQUE KEY `lt_namespace_title` (`lt_namespace`,`lt_title`)
) ENGINE=InnoDB AUTO_INCREMENT=3850 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `linter`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `linter` (
  `linter_id` int unsigned NOT NULL AUTO_INCREMENT,
  `linter_page` int unsigned NOT NULL,
  `linter_namespace` int DEFAULT NULL,
  `linter_cat` int unsigned NOT NULL,
  `linter_start` int unsigned NOT NULL,
  `linter_end` int unsigned NOT NULL,
  `linter_params` blob NOT NULL,
  `linter_template` varbinary(255) NOT NULL DEFAULT '',
  `linter_tag` varbinary(32) NOT NULL DEFAULT '',
  PRIMARY KEY (`linter_id`),
  UNIQUE KEY `linter_cat_page_position` (`linter_cat`,`linter_page`,`linter_start`,`linter_end`),
  KEY `linter_page` (`linter_page`),
  KEY `linter_cat_namespace` (`linter_cat`,`linter_namespace`),
  KEY `linter_cat_template` (`linter_cat`,`linter_template`),
  KEY `linter_cat_tag` (`linter_cat`,`linter_tag`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `log_search`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `log_search` (
  `ls_field` varbinary(32) NOT NULL,
  `ls_value` varbinary(255) NOT NULL,
  `ls_log_id` int unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`ls_field`,`ls_value`,`ls_log_id`),
  KEY `ls_log_id` (`ls_log_id`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `logging`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `logging` (
  `log_id` int unsigned NOT NULL AUTO_INCREMENT,
  `log_type` varbinary(32) NOT NULL DEFAULT '',
  `log_action` varbinary(32) NOT NULL DEFAULT '',
  `log_timestamp` binary(14) NOT NULL DEFAULT '19700101000000',
  `log_actor` bigint unsigned NOT NULL,
  `log_namespace` int NOT NULL DEFAULT '0',
  `log_title` varbinary(255) NOT NULL DEFAULT '',
  `log_page` int unsigned DEFAULT NULL,
  `log_comment_id` bigint unsigned NOT NULL,
  `log_params` blob NOT NULL,
  `log_deleted` tinyint unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`log_id`),
  KEY `log_type_time` (`log_type`,`log_timestamp`),
  KEY `log_actor_time` (`log_actor`,`log_timestamp`),
  KEY `log_page_time` (`log_namespace`,`log_title`,`log_timestamp`),
  KEY `log_times` (`log_timestamp`),
  KEY `log_actor_type_time` (`log_actor`,`log_type`,`log_timestamp`),
  KEY `log_page_id_time` (`log_page`,`log_timestamp`),
  KEY `log_type_action` (`log_type`,`log_action`,`log_timestamp`)
) ENGINE=InnoDB AUTO_INCREMENT=11601 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `module_deps`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `module_deps` (
  `md_module` varbinary(255) NOT NULL,
  `md_skin` varbinary(32) NOT NULL,
  `md_deps` mediumblob NOT NULL,
  PRIMARY KEY (`md_module`,`md_skin`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `objectcache`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `objectcache` (
  `keyname` varbinary(255) NOT NULL DEFAULT '',
  `value` mediumblob,
  `exptime` binary(14) NOT NULL,
  `modtoken` varbinary(17) NOT NULL DEFAULT '00000000000000000',
  `flags` int unsigned DEFAULT NULL,
  PRIMARY KEY (`keyname`),
  KEY `exptime` (`exptime`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `oldimage`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `oldimage` (
  `oi_name` varbinary(255) NOT NULL DEFAULT '',
  `oi_archive_name` varbinary(255) NOT NULL DEFAULT '',
  `oi_size` bigint unsigned NOT NULL DEFAULT '0',
  `oi_width` int NOT NULL DEFAULT '0',
  `oi_height` int NOT NULL DEFAULT '0',
  `oi_bits` int NOT NULL DEFAULT '0',
  `oi_description_id` bigint unsigned NOT NULL,
  `oi_actor` bigint unsigned NOT NULL,
  `oi_timestamp` binary(14) NOT NULL,
  `oi_metadata` mediumblob NOT NULL,
  `oi_media_type` enum('UNKNOWN','BITMAP','DRAWING','AUDIO','VIDEO','MULTIMEDIA','OFFICE','TEXT','EXECUTABLE','ARCHIVE','3D') DEFAULT NULL,
  `oi_major_mime` enum('unknown','application','audio','image','text','video','message','model','multipart','chemical') NOT NULL DEFAULT 'unknown',
  `oi_minor_mime` varbinary(100) NOT NULL DEFAULT 'unknown',
  `oi_deleted` tinyint unsigned NOT NULL DEFAULT '0',
  `oi_sha1` varbinary(32) NOT NULL DEFAULT '',
  KEY `oi_actor_timestamp` (`oi_actor`,`oi_timestamp`),
  KEY `oi_name_timestamp` (`oi_name`,`oi_timestamp`),
  KEY `oi_name_archive_name` (`oi_name`,`oi_archive_name`(14)),
  KEY `oi_sha1` (`oi_sha1`(10)),
  KEY `oi_timestamp` (`oi_timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `page`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `page` (
  `page_id` int unsigned NOT NULL AUTO_INCREMENT,
  `page_namespace` int NOT NULL,
  `page_title` varbinary(255) NOT NULL,
  `page_is_redirect` tinyint unsigned NOT NULL DEFAULT '0',
  `page_is_new` tinyint unsigned NOT NULL DEFAULT '0',
  `page_random` double unsigned NOT NULL,
  `page_touched` binary(14) NOT NULL,
  `page_links_updated` binary(14) DEFAULT NULL,
  `page_latest` int unsigned NOT NULL,
  `page_len` int unsigned NOT NULL,
  `page_content_model` varbinary(32) DEFAULT NULL,
  `page_lang` varbinary(35) DEFAULT NULL,
  PRIMARY KEY (`page_id`),
  UNIQUE KEY `page_name_title` (`page_namespace`,`page_title`),
  KEY `page_random` (`page_random`),
  KEY `page_len` (`page_len`),
  KEY `page_redirect_namespace_len` (`page_is_redirect`,`page_namespace`,`page_len`)
) ENGINE=InnoDB AUTO_INCREMENT=9793 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `page_props`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `page_props` (
  `pp_page` int unsigned NOT NULL,
  `pp_propname` varbinary(60) NOT NULL,
  `pp_value` blob NOT NULL,
  `pp_sortkey` float DEFAULT NULL,
  PRIMARY KEY (`pp_page`,`pp_propname`),
  UNIQUE KEY `pp_propname_page` (`pp_propname`,`pp_page`),
  UNIQUE KEY `pp_propname_sortkey_page` (`pp_propname`,`pp_sortkey`,`pp_page`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `page_restrictions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `page_restrictions` (
  `pr_id` int unsigned NOT NULL AUTO_INCREMENT,
  `pr_page` int unsigned NOT NULL,
  `pr_type` varbinary(60) NOT NULL,
  `pr_level` varbinary(60) NOT NULL,
  `pr_cascade` tinyint NOT NULL,
  `pr_expiry` varbinary(14) DEFAULT NULL,
  PRIMARY KEY (`pr_id`),
  UNIQUE KEY `pr_pagetype` (`pr_page`,`pr_type`),
  KEY `pr_typelevel` (`pr_type`,`pr_level`),
  KEY `pr_level` (`pr_level`),
  KEY `pr_cascade` (`pr_cascade`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `pagelinks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pagelinks` (
  `pl_from` int unsigned NOT NULL DEFAULT '0',
  `pl_from_namespace` int NOT NULL DEFAULT '0',
  `pl_target_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`pl_from`,`pl_target_id`),
  KEY `pl_target_id` (`pl_target_id`,`pl_from`),
  KEY `pl_backlinks_namespace_target_id` (`pl_from_namespace`,`pl_target_id`,`pl_from`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `protected_titles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `protected_titles` (
  `pt_namespace` int NOT NULL,
  `pt_title` varbinary(255) NOT NULL,
  `pt_user` int unsigned NOT NULL,
  `pt_reason_id` bigint unsigned NOT NULL,
  `pt_timestamp` binary(14) NOT NULL,
  `pt_expiry` varbinary(14) NOT NULL,
  `pt_create_perm` varbinary(60) NOT NULL,
  PRIMARY KEY (`pt_namespace`,`pt_title`),
  KEY `pt_timestamp` (`pt_timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `querycache`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `querycache` (
  `qc_type` varbinary(32) NOT NULL,
  `qc_value` int unsigned NOT NULL DEFAULT '0',
  `qc_namespace` int NOT NULL DEFAULT '0',
  `qc_title` varbinary(255) NOT NULL DEFAULT '',
  KEY `qc_type` (`qc_type`,`qc_value`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `querycache_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `querycache_info` (
  `qci_type` varbinary(32) NOT NULL DEFAULT '',
  `qci_timestamp` binary(14) NOT NULL DEFAULT '19700101000000',
  PRIMARY KEY (`qci_type`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `querycachetwo`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `querycachetwo` (
  `qcc_type` varbinary(32) NOT NULL,
  `qcc_value` int unsigned NOT NULL DEFAULT '0',
  `qcc_namespace` int NOT NULL DEFAULT '0',
  `qcc_title` varbinary(255) NOT NULL DEFAULT '',
  `qcc_namespacetwo` int NOT NULL DEFAULT '0',
  `qcc_titletwo` varbinary(255) NOT NULL DEFAULT '',
  KEY `qcc_type` (`qcc_type`,`qcc_value`),
  KEY `qcc_title` (`qcc_type`,`qcc_namespace`,`qcc_title`),
  KEY `qcc_titletwo` (`qcc_type`,`qcc_namespacetwo`,`qcc_titletwo`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `recentchanges`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `recentchanges` (
  `rc_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `rc_timestamp` binary(14) NOT NULL,
  `rc_actor` bigint unsigned NOT NULL,
  `rc_namespace` int NOT NULL DEFAULT '0',
  `rc_title` varbinary(255) NOT NULL DEFAULT '',
  `rc_comment_id` bigint unsigned NOT NULL,
  `rc_minor` tinyint unsigned NOT NULL DEFAULT '0',
  `rc_bot` tinyint unsigned NOT NULL DEFAULT '0',
  `rc_new` tinyint unsigned NOT NULL DEFAULT '0',
  `rc_cur_id` int unsigned NOT NULL DEFAULT '0',
  `rc_this_oldid` int unsigned NOT NULL DEFAULT '0',
  `rc_last_oldid` int unsigned NOT NULL DEFAULT '0',
  `rc_type` tinyint unsigned NOT NULL DEFAULT '0',
  `rc_source` varbinary(16) NOT NULL DEFAULT '',
  `rc_patrolled` tinyint unsigned NOT NULL DEFAULT '0',
  `rc_ip` varbinary(40) NOT NULL DEFAULT '',
  `rc_old_len` int DEFAULT NULL,
  `rc_new_len` int DEFAULT NULL,
  `rc_deleted` tinyint unsigned NOT NULL DEFAULT '0',
  `rc_logid` int unsigned NOT NULL DEFAULT '0',
  `rc_log_type` varbinary(255) DEFAULT NULL,
  `rc_log_action` varbinary(255) DEFAULT NULL,
  `rc_params` blob,
  PRIMARY KEY (`rc_id`),
  KEY `rc_timestamp` (`rc_timestamp`),
  KEY `rc_namespace_title_timestamp` (`rc_namespace`,`rc_title`,`rc_timestamp`),
  KEY `rc_cur_id` (`rc_cur_id`),
  KEY `rc_new_name_timestamp` (`rc_new`,`rc_namespace`,`rc_timestamp`),
  KEY `rc_ip` (`rc_ip`),
  KEY `rc_ns_actor` (`rc_namespace`,`rc_actor`),
  KEY `rc_actor` (`rc_actor`,`rc_timestamp`),
  KEY `rc_name_type_patrolled_timestamp` (`rc_namespace`,`rc_type`,`rc_patrolled`,`rc_timestamp`),
  KEY `rc_this_oldid` (`rc_this_oldid`)
) ENGINE=InnoDB AUTO_INCREMENT=28566 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `redirect`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `redirect` (
  `rd_from` int unsigned NOT NULL DEFAULT '0',
  `rd_namespace` int NOT NULL DEFAULT '0',
  `rd_title` varbinary(255) NOT NULL DEFAULT '',
  `rd_interwiki` varbinary(32) DEFAULT NULL,
  `rd_fragment` varbinary(255) DEFAULT NULL,
  PRIMARY KEY (`rd_from`),
  KEY `rd_ns_title` (`rd_namespace`,`rd_title`,`rd_from`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `revision`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `revision` (
  `rev_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `rev_page` int unsigned NOT NULL,
  `rev_comment_id` bigint unsigned NOT NULL,
  `rev_actor` bigint unsigned NOT NULL,
  `rev_timestamp` binary(14) NOT NULL,
  `rev_minor_edit` tinyint unsigned NOT NULL DEFAULT '0',
  `rev_deleted` tinyint unsigned NOT NULL DEFAULT '0',
  `rev_len` int unsigned DEFAULT NULL,
  `rev_parent_id` bigint unsigned DEFAULT NULL,
  `rev_sha1` varbinary(32) NOT NULL DEFAULT '',
  PRIMARY KEY (`rev_id`),
  KEY `rev_timestamp` (`rev_timestamp`),
  KEY `rev_page_timestamp` (`rev_page`,`rev_timestamp`),
  KEY `rev_actor_timestamp` (`rev_actor`,`rev_timestamp`,`rev_id`),
  KEY `rev_page_actor_timestamp` (`rev_page`,`rev_actor`,`rev_timestamp`)
) ENGINE=InnoDB AUTO_INCREMENT=22033 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `searchindex`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `searchindex` (
  `si_page` int unsigned NOT NULL,
  `si_title` mediumtext NOT NULL,
  `si_text` longtext NOT NULL,
  PRIMARY KEY (`si_page`),
  FULLTEXT KEY `si_title` (`si_title`),
  FULLTEXT KEY `si_text` (`si_text`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `site_identifiers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `site_identifiers` (
  `si_type` varbinary(32) NOT NULL,
  `si_key` varbinary(32) NOT NULL,
  `si_site` int unsigned NOT NULL,
  PRIMARY KEY (`si_type`,`si_key`),
  KEY `si_site` (`si_site`),
  KEY `si_key` (`si_key`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `site_stats`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `site_stats` (
  `ss_row_id` int unsigned NOT NULL,
  `ss_total_edits` bigint unsigned DEFAULT NULL,
  `ss_good_articles` bigint unsigned DEFAULT NULL,
  `ss_total_pages` bigint unsigned DEFAULT NULL,
  `ss_users` bigint unsigned DEFAULT NULL,
  `ss_active_users` bigint unsigned DEFAULT NULL,
  `ss_images` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`ss_row_id`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `sites`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sites` (
  `site_id` int unsigned NOT NULL AUTO_INCREMENT,
  `site_global_key` varbinary(64) NOT NULL,
  `site_type` varbinary(32) NOT NULL,
  `site_group` varbinary(32) NOT NULL,
  `site_source` varbinary(32) NOT NULL,
  `site_language` varbinary(35) NOT NULL,
  `site_protocol` varbinary(32) NOT NULL,
  `site_domain` varbinary(255) NOT NULL,
  `site_data` blob NOT NULL,
  `site_forward` tinyint(1) NOT NULL,
  `site_config` blob NOT NULL,
  PRIMARY KEY (`site_id`),
  UNIQUE KEY `site_global_key` (`site_global_key`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `slot_roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `slot_roles` (
  `role_id` int NOT NULL AUTO_INCREMENT,
  `role_name` varbinary(64) NOT NULL,
  PRIMARY KEY (`role_id`),
  UNIQUE KEY `role_name` (`role_name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `slots`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `slots` (
  `slot_revision_id` bigint unsigned NOT NULL,
  `slot_role_id` smallint unsigned NOT NULL,
  `slot_content_id` bigint unsigned NOT NULL,
  `slot_origin` bigint unsigned NOT NULL,
  PRIMARY KEY (`slot_revision_id`,`slot_role_id`),
  KEY `slot_revision_origin_role` (`slot_revision_id`,`slot_origin`,`slot_role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_concept_cache`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_concept_cache` (
  `s_id` int unsigned NOT NULL,
  `o_id` int unsigned NOT NULL,
  KEY `o_id` (`o_id`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_di_blob`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_di_blob` (
  `s_id` int unsigned NOT NULL,
  `p_id` int unsigned NOT NULL,
  `o_blob` mediumblob,
  `o_hash` varbinary(255) DEFAULT NULL,
  KEY `s_id` (`s_id`,`p_id`),
  KEY `o_hash` (`o_hash`),
  KEY `s_id_2` (`s_id`,`o_hash`),
  KEY `p_id` (`p_id`,`o_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_di_bool`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_di_bool` (
  `s_id` int unsigned NOT NULL,
  `p_id` int unsigned NOT NULL,
  `o_value` tinyint(1) DEFAULT NULL,
  KEY `s_id` (`s_id`,`p_id`),
  KEY `o_value` (`o_value`),
  KEY `p_id` (`p_id`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_di_coords`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_di_coords` (
  `s_id` int unsigned NOT NULL,
  `p_id` int unsigned NOT NULL,
  `o_serialized` varbinary(255) DEFAULT NULL,
  `o_lat` double DEFAULT NULL,
  `o_lon` double DEFAULT NULL,
  KEY `s_id` (`s_id`,`p_id`),
  KEY `o_serialized` (`o_serialized`),
  KEY `p_id` (`p_id`,`o_serialized`),
  KEY `o_lat` (`o_lat`,`o_lon`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_di_number`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_di_number` (
  `s_id` int unsigned NOT NULL,
  `p_id` int unsigned NOT NULL,
  `o_serialized` varbinary(255) DEFAULT NULL,
  `o_sortkey` double DEFAULT NULL,
  KEY `s_id` (`s_id`,`p_id`),
  KEY `o_sortkey` (`o_sortkey`),
  KEY `p_id` (`p_id`,`o_serialized`),
  KEY `p_id_2` (`p_id`,`o_sortkey`),
  KEY `s_id_2` (`s_id`,`p_id`,`o_sortkey`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_di_time`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_di_time` (
  `s_id` int unsigned NOT NULL,
  `p_id` int unsigned NOT NULL,
  `o_serialized` varbinary(255) DEFAULT NULL,
  `o_sortkey` double DEFAULT NULL,
  KEY `s_id` (`s_id`,`p_id`),
  KEY `o_sortkey` (`o_sortkey`),
  KEY `p_id` (`p_id`,`o_serialized`),
  KEY `p_id_2` (`p_id`,`o_sortkey`),
  KEY `s_id_2` (`s_id`,`p_id`,`o_sortkey`,`o_serialized`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_di_uri`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_di_uri` (
  `s_id` int unsigned NOT NULL,
  `p_id` int unsigned NOT NULL,
  `o_blob` mediumblob,
  `o_serialized` varbinary(255) DEFAULT NULL,
  KEY `s_id` (`s_id`,`p_id`),
  KEY `o_serialized` (`o_serialized`),
  KEY `p_id` (`p_id`,`o_serialized`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_di_wikipage`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_di_wikipage` (
  `s_id` int unsigned NOT NULL,
  `p_id` int unsigned NOT NULL,
  `o_id` int unsigned DEFAULT NULL,
  KEY `s_id` (`s_id`,`p_id`),
  KEY `o_id` (`o_id`),
  KEY `p_id` (`p_id`,`s_id`),
  KEY `s_id_2` (`s_id`,`o_id`),
  KEY `s_id_3` (`s_id`,`p_id`,`o_id`),
  KEY `o_id_2` (`o_id`,`s_id`),
  KEY `o_id_3` (`o_id`,`p_id`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_ask`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_ask` (
  `s_id` int unsigned NOT NULL,
  `o_id` int unsigned DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_id` (`o_id`),
  KEY `s_id_2` (`s_id`,`o_id`),
  KEY `o_id_2` (`o_id`,`s_id`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_askde`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_askde` (
  `s_id` int unsigned NOT NULL,
  `o_serialized` varbinary(255) DEFAULT NULL,
  `o_sortkey` double DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_sortkey` (`o_sortkey`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_askdu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_askdu` (
  `s_id` int unsigned NOT NULL,
  `o_serialized` varbinary(255) DEFAULT NULL,
  `o_sortkey` double DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_sortkey` (`o_sortkey`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_askfo`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_askfo` (
  `s_id` int unsigned NOT NULL,
  `o_blob` mediumblob,
  `o_hash` varbinary(255) DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_hash` (`o_hash`),
  KEY `s_id_2` (`s_id`,`o_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_askpa`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_askpa` (
  `s_id` int unsigned NOT NULL,
  `o_blob` mediumblob,
  `o_hash` varbinary(255) DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_hash` (`o_hash`),
  KEY `s_id_2` (`s_id`,`o_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_asksi`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_asksi` (
  `s_id` int unsigned NOT NULL,
  `o_serialized` varbinary(255) DEFAULT NULL,
  `o_sortkey` double DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_sortkey` (`o_sortkey`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_askst`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_askst` (
  `s_id` int unsigned NOT NULL,
  `o_blob` mediumblob,
  `o_hash` varbinary(255) DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_hash` (`o_hash`),
  KEY `s_id_2` (`s_id`,`o_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_cdat`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_cdat` (
  `s_id` int unsigned NOT NULL,
  `o_serialized` varbinary(255) DEFAULT NULL,
  `o_sortkey` double DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_sortkey` (`o_sortkey`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_conc`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_conc` (
  `s_id` int unsigned NOT NULL,
  `concept_txt` mediumblob,
  `concept_docu` mediumblob,
  `concept_features` int DEFAULT NULL,
  `concept_size` int DEFAULT NULL,
  `concept_depth` int DEFAULT NULL,
  `cache_date` int unsigned DEFAULT NULL,
  `cache_count` int unsigned DEFAULT NULL,
  KEY `s_id` (`s_id`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_conv`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_conv` (
  `s_id` int unsigned NOT NULL,
  `o_blob` mediumblob,
  `o_hash` varbinary(255) DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_hash` (`o_hash`),
  KEY `s_id_2` (`s_id`,`o_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_dtitle`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_dtitle` (
  `s_id` int unsigned NOT NULL,
  `o_blob` mediumblob,
  `o_hash` varbinary(255) DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_hash` (`o_hash`),
  KEY `s_id_2` (`s_id`,`o_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_impo`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_impo` (
  `s_id` int unsigned NOT NULL,
  `o_blob` mediumblob,
  `o_hash` varbinary(255) DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_hash` (`o_hash`),
  KEY `s_id_2` (`s_id`,`o_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_inst`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_inst` (
  `s_id` int unsigned NOT NULL,
  `o_id` int unsigned DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_id` (`o_id`),
  KEY `s_id_2` (`s_id`,`o_id`),
  KEY `o_id_2` (`o_id`,`s_id`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_lcode`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_lcode` (
  `s_id` int unsigned NOT NULL,
  `o_blob` mediumblob,
  `o_hash` varbinary(255) DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_hash` (`o_hash`),
  KEY `s_id_2` (`s_id`,`o_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_ledt`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_ledt` (
  `s_id` int unsigned NOT NULL,
  `o_id` int unsigned DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_id` (`o_id`),
  KEY `s_id_2` (`s_id`,`o_id`),
  KEY `o_id_2` (`o_id`,`s_id`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_list`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_list` (
  `s_id` int unsigned NOT NULL,
  `o_blob` mediumblob,
  `o_hash` varbinary(255) DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_hash` (`o_hash`),
  KEY `s_id_2` (`s_id`,`o_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_mdat`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_mdat` (
  `s_id` int unsigned NOT NULL,
  `o_serialized` varbinary(255) DEFAULT NULL,
  `o_sortkey` double DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_sortkey` (`o_sortkey`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_pplb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_pplb` (
  `s_id` int unsigned NOT NULL,
  `o_id` int unsigned DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_id` (`o_id`),
  KEY `s_id_2` (`s_id`,`o_id`),
  KEY `o_id_2` (`o_id`,`s_id`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_prec`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_prec` (
  `s_id` int unsigned NOT NULL,
  `o_serialized` varbinary(255) DEFAULT NULL,
  `o_sortkey` double DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_sortkey` (`o_sortkey`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_pval`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_pval` (
  `s_id` int unsigned NOT NULL,
  `o_blob` mediumblob,
  `o_hash` varbinary(255) DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_hash` (`o_hash`),
  KEY `s_id_2` (`s_id`,`o_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_redi`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_redi` (
  `s_title` varbinary(255) NOT NULL,
  `s_namespace` int NOT NULL,
  `o_id` int unsigned DEFAULT NULL,
  KEY `s_title` (`s_title`,`s_namespace`,`o_id`),
  KEY `o_id` (`o_id`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_serv`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_serv` (
  `s_id` int unsigned NOT NULL,
  `o_blob` mediumblob,
  `o_hash` varbinary(255) DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_hash` (`o_hash`),
  KEY `s_id_2` (`s_id`,`o_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_sobj`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_sobj` (
  `s_id` int unsigned NOT NULL,
  `o_id` int unsigned DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_id` (`o_id`),
  KEY `s_id_2` (`s_id`,`o_id`),
  KEY `o_id_2` (`o_id`,`s_id`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_subc`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_subc` (
  `s_id` int unsigned NOT NULL,
  `o_id` int unsigned DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_id` (`o_id`),
  KEY `s_id_2` (`s_id`,`o_id`),
  KEY `o_id_2` (`o_id`,`s_id`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_subp`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_subp` (
  `s_id` int unsigned NOT NULL,
  `o_id` int unsigned DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_id` (`o_id`),
  KEY `s_id_2` (`s_id`,`o_id`),
  KEY `o_id_2` (`o_id`,`s_id`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_text`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_text` (
  `s_id` int unsigned NOT NULL,
  `o_blob` mediumblob,
  `o_hash` varbinary(255) DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_hash` (`o_hash`),
  KEY `s_id_2` (`s_id`,`o_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_type` (
  `s_id` int unsigned NOT NULL,
  `o_blob` mediumblob,
  `o_serialized` varbinary(255) DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_serialized` (`o_serialized`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_unit`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_unit` (
  `s_id` int unsigned NOT NULL,
  `o_blob` mediumblob,
  `o_hash` varbinary(255) DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_hash` (`o_hash`),
  KEY `s_id_2` (`s_id`,`o_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_fpt_uri`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_fpt_uri` (
  `s_id` int unsigned NOT NULL,
  `o_blob` mediumblob,
  `o_serialized` varbinary(255) DEFAULT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_serialized` (`o_serialized`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_object_aux`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_object_aux` (
  `smw_id` int unsigned NOT NULL,
  `smw_seqmap` mediumblob,
  `smw_countmap` mediumblob,
  PRIMARY KEY (`smw_id`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_object_ids`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_object_ids` (
  `smw_id` int unsigned NOT NULL AUTO_INCREMENT,
  `smw_namespace` int NOT NULL,
  `smw_title` varbinary(255) NOT NULL,
  `smw_iw` varbinary(32) NOT NULL,
  `smw_subobject` varbinary(255) NOT NULL,
  `smw_sortkey` varbinary(255) NOT NULL,
  `smw_sort` varbinary(255) DEFAULT NULL,
  `smw_proptable_hash` mediumblob,
  `smw_hash` varbinary(40) DEFAULT NULL,
  `smw_rev` int unsigned DEFAULT NULL,
  `smw_touched` binary(14) DEFAULT NULL,
  PRIMARY KEY (`smw_id`),
  KEY `smw_id` (`smw_id`,`smw_sortkey`),
  KEY `smw_hash` (`smw_hash`,`smw_id`),
  KEY `smw_iw` (`smw_iw`),
  KEY `smw_iw_2` (`smw_iw`,`smw_id`),
  KEY `smw_title` (`smw_title`,`smw_namespace`,`smw_iw`,`smw_subobject`),
  KEY `smw_sortkey` (`smw_sortkey`),
  KEY `smw_sort` (`smw_sort`,`smw_id`),
  KEY `smw_namespace` (`smw_namespace`,`smw_sortkey`),
  KEY `smw_rev` (`smw_rev`,`smw_id`),
  KEY `smw_touched` (`smw_touched`,`smw_id`)
) ENGINE=InnoDB AUTO_INCREMENT=861666 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_prop_stats`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_prop_stats` (
  `p_id` int unsigned DEFAULT NULL,
  `usage_count` int unsigned DEFAULT NULL,
  `null_count` int unsigned DEFAULT NULL,
  UNIQUE KEY `p_id` (`p_id`),
  KEY `usage_count` (`usage_count`),
  KEY `null_count` (`null_count`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `smw_query_links`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `smw_query_links` (
  `s_id` int unsigned NOT NULL,
  `o_id` int unsigned NOT NULL,
  KEY `s_id` (`s_id`),
  KEY `o_id` (`o_id`),
  KEY `s_id_2` (`s_id`,`o_id`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `templatelinks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `templatelinks` (
  `tl_from` int unsigned NOT NULL DEFAULT '0',
  `tl_target_id` bigint unsigned NOT NULL,
  `tl_from_namespace` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`tl_from`,`tl_target_id`),
  KEY `tl_target_id` (`tl_target_id`,`tl_from`),
  KEY `tl_backlinks_namespace_target_id` (`tl_from_namespace`,`tl_target_id`,`tl_from`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `text`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `text` (
  `old_id` int unsigned NOT NULL AUTO_INCREMENT,
  `old_text` mediumblob NOT NULL,
  `old_flags` tinyblob NOT NULL,
  PRIMARY KEY (`old_id`)
) ENGINE=InnoDB AUTO_INCREMENT=20013 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `updatelog`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `updatelog` (
  `ul_key` varbinary(255) NOT NULL,
  `ul_value` blob,
  PRIMARY KEY (`ul_key`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `uploadstash`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `uploadstash` (
  `us_id` int unsigned NOT NULL AUTO_INCREMENT,
  `us_user` int unsigned NOT NULL,
  `us_key` varbinary(255) NOT NULL,
  `us_orig_path` varbinary(255) NOT NULL,
  `us_path` varbinary(255) NOT NULL,
  `us_source_type` varbinary(50) DEFAULT NULL,
  `us_timestamp` binary(14) NOT NULL,
  `us_status` varbinary(50) NOT NULL,
  `us_chunk_inx` int unsigned DEFAULT NULL,
  `us_props` blob,
  `us_size` bigint unsigned NOT NULL,
  `us_sha1` varbinary(31) NOT NULL,
  `us_mime` varbinary(255) DEFAULT NULL,
  `us_media_type` enum('UNKNOWN','BITMAP','DRAWING','AUDIO','VIDEO','MULTIMEDIA','OFFICE','TEXT','EXECUTABLE','ARCHIVE','3D') DEFAULT NULL,
  `us_image_width` int unsigned DEFAULT NULL,
  `us_image_height` int unsigned DEFAULT NULL,
  `us_image_bits` smallint unsigned DEFAULT NULL,
  PRIMARY KEY (`us_id`),
  UNIQUE KEY `us_key` (`us_key`),
  KEY `us_user` (`us_user`),
  KEY `us_timestamp` (`us_timestamp`)
) ENGINE=InnoDB AUTO_INCREMENT=523 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user` (
  `user_id` int unsigned NOT NULL AUTO_INCREMENT,
  `user_name` varbinary(255) NOT NULL DEFAULT '',
  `user_real_name` varbinary(255) NOT NULL DEFAULT '',
  `user_password` tinyblob NOT NULL,
  `user_newpassword` tinyblob NOT NULL,
  `user_newpass_time` binary(14) DEFAULT NULL,
  `user_email` tinyblob NOT NULL,
  `user_touched` binary(14) NOT NULL,
  `user_token` binary(32) NOT NULL DEFAULT '\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',
  `user_email_authenticated` binary(14) DEFAULT NULL,
  `user_email_token` binary(32) DEFAULT NULL,
  `user_email_token_expires` binary(14) DEFAULT NULL,
  `user_registration` binary(14) DEFAULT NULL,
  `user_editcount` int unsigned DEFAULT NULL,
  `user_password_expires` varbinary(14) DEFAULT NULL,
  `user_is_temp` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `user_name` (`user_name`),
  KEY `user_email_token` (`user_email_token`),
  KEY `user_email` (`user_email`(50))
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `user_autocreate_serial`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_autocreate_serial` (
  `uas_shard` int unsigned NOT NULL,
  `uas_value` int unsigned NOT NULL,
  `uas_year` smallint unsigned NOT NULL,
  PRIMARY KEY (`uas_shard`,`uas_year`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `user_former_groups`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_former_groups` (
  `ufg_user` int unsigned NOT NULL DEFAULT '0',
  `ufg_group` varbinary(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`ufg_user`,`ufg_group`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `user_groups`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_groups` (
  `ug_user` int unsigned NOT NULL DEFAULT '0',
  `ug_group` varbinary(255) NOT NULL DEFAULT '',
  `ug_expiry` varbinary(14) DEFAULT NULL,
  PRIMARY KEY (`ug_user`,`ug_group`),
  KEY `ug_group` (`ug_group`),
  KEY `ug_expiry` (`ug_expiry`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `user_newtalk`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_newtalk` (
  `user_id` int unsigned NOT NULL DEFAULT '0',
  `user_ip` varbinary(40) NOT NULL DEFAULT '',
  `user_last_timestamp` binary(14) DEFAULT NULL,
  KEY `un_user_id` (`user_id`),
  KEY `un_user_ip` (`user_ip`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `user_properties`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_properties` (
  `up_user` int unsigned NOT NULL,
  `up_property` varbinary(255) NOT NULL,
  `up_value` blob,
  PRIMARY KEY (`up_user`,`up_property`),
  KEY `up_property` (`up_property`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `user_reference`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_reference` (
  `user_id` int unsigned NOT NULL AUTO_INCREMENT,
  `user_name` varbinary(255) NOT NULL DEFAULT '',
  `user_real_name` varbinary(255) NOT NULL DEFAULT '',
  `user_password` tinyblob NOT NULL,
  `user_newpassword` tinyblob NOT NULL,
  `user_newpass_time` binary(14) DEFAULT NULL,
  `user_email` tinyblob NOT NULL,
  `user_touched` binary(14) NOT NULL,
  `user_token` binary(32) NOT NULL DEFAULT '\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',
  `user_email_authenticated` binary(14) DEFAULT NULL,
  `user_email_token` binary(32) DEFAULT NULL,
  `user_email_token_expires` binary(14) DEFAULT NULL,
  `user_registration` binary(14) DEFAULT NULL,
  `user_editcount` int unsigned DEFAULT NULL,
  `user_password_expires` varbinary(14) DEFAULT NULL,
  `user_is_temp` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `user_name` (`user_name`),
  KEY `user_email_token` (`user_email_token`),
  KEY `user_email` (`user_email`(50))
) ENGINE=InnoDB AUTO_INCREMENT=243 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `watchlist`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `watchlist` (
  `wl_id` int unsigned NOT NULL AUTO_INCREMENT,
  `wl_user` int unsigned NOT NULL,
  `wl_namespace` int NOT NULL DEFAULT '0',
  `wl_title` varbinary(255) NOT NULL DEFAULT '',
  `wl_notificationtimestamp` binary(14) DEFAULT NULL,
  PRIMARY KEY (`wl_id`),
  UNIQUE KEY `wl_user` (`wl_user`,`wl_namespace`,`wl_title`),
  KEY `wl_namespace_title` (`wl_namespace`,`wl_title`),
  KEY `wl_user_notificationtimestamp` (`wl_user`,`wl_notificationtimestamp`)
) ENGINE=InnoDB AUTO_INCREMENT=7727 DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
DROP TABLE IF EXISTS `watchlist_expiry`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `watchlist_expiry` (
  `we_item` int unsigned NOT NULL,
  `we_expiry` binary(14) NOT NULL,
  PRIMARY KEY (`we_item`),
  KEY `we_expiry` (`we_expiry`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!50001 DROP VIEW IF EXISTS `dpl_clview`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`root`@`localhost` SQL SECURITY DEFINER */
/*!50001 VIEW `dpl_clview` AS select 1 AS `cl_to`,1 AS `cl_from`,1 AS `cl_sortkey` */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

