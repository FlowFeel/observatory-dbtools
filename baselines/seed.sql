
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

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES (1,_binary 'Admin','',_binary ':pbkdf2:sha512:30000:64:2wMmqt+F0td7kGP1svkrqA==:Ibl05h8zfjVRnv6Dy9aNwn+w8DsDUILGCxJ+sA1lodgQiOAYn2uG5K+6b389xO7H9QYIEoh3aT9r8Cb0PIjKYA==','',NULL,_binary 'info@observatory.wiki',_binary '20260713204502',_binary '011a9c92c7d4bdcc1376be7d42cb96a6',_binary '20250711112440',NULL,NULL,_binary '20241219091034',11204,NULL,0),(2,_binary 'MediaWiki default','','','',NULL,'',_binary '20241219091034',_binary '*** INVALID ***\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',NULL,NULL,NULL,_binary '20241219091034',1,NULL,0),(3,_binary 'Maintenance script','','','',NULL,'',_binary '20241219122240',_binary '*** INVALID ***\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',NULL,NULL,NULL,_binary '20241219122240',806,NULL,0),(4,_binary 'Info',_binary 'The Observatory',_binary ':pbkdf2:sha512:30000:64:YPzjEb92mW2lbgN2RvoReQ==:JGBI0WCNtgjU3kW0OJXxhDrGzifgF29UIK9h9EPl6vo3/ESj5FhN/ehwEvnBooPke1OC3dnFgwhgCzLaUMckEw==','',NULL,_binary 'info@observatory.wiki',_binary '20250315161056',_binary '1d6781af2309ff9db959673e5c80894c',_binary '20250315161005',NULL,NULL,_binary '20250315160914',0,NULL,0),(5,_binary 'Jenny',_binary 'Jenny Pierson',_binary ':pbkdf2:sha512:30000:64:Bvj/0TQiuvkQoeVarzBdrA==:2dI075wE5hxRQaJsoldSK/gkxuBwYCusx7QBG0WF/aKi4Tzzh1Vdyd1OdO1O4curKLCWx8rp43+XAGKomDDFkQ==','',NULL,_binary 'jenny@ind.media',_binary '20251208174119',_binary '125ffbdec9dfec1b6d643dbf45930d3b',_binary '20250415150718',_binary 'ef2774d40f0403cc5dbe88aae1040b52',_binary '20250422150224',_binary '20250408120145',10,NULL,0),(6,_binary 'Jan',_binary 'Jan Ritch-Frel',_binary ':pbkdf2:sha512:30000:64:V579Q3lRJNrV0/pRV095dQ==:yQRrFttovX1A846cOccgTLEoxFvtItOuHRvspk0pEam8q35IGrk0P1iFG/wQCXiGzAM2M88YIxHzl1NtLTvYzA==','',NULL,_binary 'jan@ind.media',_binary '20260502120519',_binary 'caefb7edf533b0346423b943ac5f7f5e',_binary '20250711112847',NULL,NULL,_binary '20250408120152',2,NULL,0),(7,_binary 'Rey',_binary 'Reynard Loki',_binary ':pbkdf2:sha512:30000:64:uJMJD9am5buHsKQfowuXvA==:R+7xi6zS5FPsl7m0nFdLN15sZ53plX7uT1W3t+Z4MG/JE1paxkbiMCRbbBQINsZbmwI0ieY/31yV2VrW+gqClw==',_binary ':pbkdf2:sha512:30000:64:b2J2/QgvfCOfT9jVcbAVLw==:w6/SmB/tgXucLhbzz9kpdkdhnHN5Pk5wjdEI1cXiV+dKSxgCd6OIoOh5JjADHGBa+ftS/YkbMZu7btrwPKgBTQ==',_binary '20250729132242',_binary 'reynard@ind.media',_binary '20260508042309',_binary 'ff88b77b17d3c89f571a021e403b7c7e',_binary '20250711113032',NULL,NULL,_binary '20250408120200',16,NULL,0),(8,_binary 'Brittani',_binary 'Brittani Banks',_binary ':pbkdf2:sha512:30000:64:bUywXRf2qno5G1ymv8tRxg==:JZBLhESkA6LDy3ETye4OUPsTK7mGPU7ko1pf+Djj1S5FXxhWSDyG8TamtwKSw08/e8EXdWwYIJhRkcm0r3PGQA==','',NULL,_binary 'brittani@ind.media',_binary '20250711112601',_binary 'f2a3dbb01395879dac3e70cdb4eca501',_binary '20250417162902',_binary '33e778b23dea437ea601d47e92b94e78',_binary '20250424162843',_binary '20250408120437',0,NULL,0),(9,_binary 'Katherine',_binary 'Katherine Dolan',_binary ':pbkdf2:sha512:30000:64:gqGIqK/0QH8ZU/DtM1cYeA==:LmauelIjsK34GcNum47DVLaEXHSt/iHhfKFHGskv+V0vu3YK3pVVEQrkq8Ut3dU3Js7ZRHCb+6s+uRV+qLQz8g==','',NULL,_binary 'katherineliddy@gmail.com',_binary '20250711112818',_binary '0530e152a17835f4c6032bfb2078fe91',_binary '20250711112818',NULL,NULL,_binary '20250408120452',0,NULL,0),(10,_binary 'Anton',_binary 'Anton Krom',_binary ':pbkdf2:sha512:30000:64:JkNFFuYOBhkWoKBpvQ9jFw==:2SEcpuZFtdqqcGonBOHgb15jkuf6maSr89kW3FZYmxySbTy04SILAALdiT0XLqbsl7BNWhDFH/skDd+knZqxzg==','',NULL,_binary 'wikivisor@wikivisor.com',_binary '20260701070854',_binary '63a8777ed1258bb4c096ceb878902b6e',_binary '20260216173618',_binary '32fd46d0bcc509425c646203a51119f8',_binary '20260223173559',_binary '20250408120517',3413,NULL,0),(11,_binary 'Kelly',_binary 'Kelly DeLay',_binary ':pbkdf2:sha512:30000:64:LH6h6By6KX7vF+QKJNZbsw==:FfsYN/+OMPC4ozx6kwKYJSzbr+BmuOOlyv7jo1zGk/UKRGjRjVdskXLq0kX7/wZ0dzIE/i3UtRs/E1qSooNfQQ==','',NULL,_binary 'kelly@kellydelay.com',_binary '20251103191522',_binary '34d33fd0ba405bf9f6d06f04a407f226',NULL,_binary '99a21187053da3ddd22f9c8c9f94482c',_binary '20250422170750',_binary '20250408121059',6,NULL,0),(12,_binary 'Irina',_binary 'Irina Matuzava',_binary ':pbkdf2:sha512:30000:64:VGO39TBrY22oHkZ7mMPOHg==:7lCc7pqPgCTsj53UhXVw+QfK/DzII4DhwAj38X1GhecwFCMlHMHAsgZn7+vzoUbjCjCMNIswTcloiUWhXWCmUQ==','',NULL,_binary 'irina@ind.media',_binary '20260706204511',_binary 'c3ddee4f7761b181b0b41d5cc2cc7411',_binary '20250711063232',NULL,NULL,_binary '20250408121107',1187,NULL,0),(13,_binary 'Delete page script','','','',NULL,'',_binary '20250409083806',_binary '*** INVALID ***\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',NULL,NULL,NULL,_binary '20250409083806',0,NULL,0),(14,_binary 'Redirect fixer','','','',NULL,'',_binary '20250520144840',_binary '0445e16f89b3237dd7d2ff83c634b9a4',NULL,NULL,NULL,_binary '20250520144840',2,NULL,0),(16,_binary 'Lorenzo Hostetter',_binary 'Lorenzo Hostetter',_binary ':pbkdf2:sha512:30000:64:6Q1kmOGIQu3CZE0yFx12Zw==:FdgyGQv42vvZAuaa7ZZFaAJwc0sSOnZIF/d210/yCUdil14XfsFVtfzbetrvPzSRb6aqIbZ0yR0tLLtBko6h2w==','',NULL,_binary 'lorenzo.hofstetter@phersu-atlas.com',_binary '20251125163711',_binary 'ebe3e45aa33a55a4d533139064073bb0',NULL,_binary 'aa26bdd3f8a8501fbeac6296afee5b5c',_binary '20251202163708',_binary '20251125163708',0,NULL,0),(17,_binary 'DownloadBookStash','','','',NULL,'',_binary '20260525122108',_binary '*** INVALID ***\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',NULL,NULL,NULL,_binary '20260525122108',0,NULL,0),(18,_binary 'FlowBot','',_binary ':pbkdf2:sha512:30000:64:Bsn4tDgUiALZkUHwdiocrA==:qut6ioNrdVUSouRsj44834gViBYGwR+/X3oe4/+C6Ibz5viBoXKC5OzUNGBslfeaI1GC3TGuvWmE/a7cOVaRyA==','',NULL,'',_binary '20260619032108',_binary '4da55c44fdd3ed6eb37b242a36a573d7',NULL,_binary '\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',NULL,_binary '20260619032107',330,NULL,0),(19,_binary 'RebuildDataLogger','','','',NULL,'',_binary '20260626221031',_binary '*** INVALID ***\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',NULL,NULL,NULL,_binary '20260626221031',0,NULL,0),(20,_binary 'FeelingFlowing',_binary 'Flow-Bot',_binary ':pbkdf2:sha512:30000:64:KjSJBBrkieEP8xm9x+UCCA==:JoB9LQY338ADFngmsLhFbouASdTCGEicCziDfUF4d0b36ByAVs5iNKauetbUSQyNhaDgcE8rFv2uMQXr7fnRoA==','',NULL,'',_binary '20260629214441',_binary 'b0c9e5d5866198603d149b3dc4c91331',NULL,_binary '\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',NULL,_binary '20260629214440',223,NULL,0),(21,_binary 'Move page script','','','',NULL,'',_binary '20260709115200',_binary '*** INVALID ***\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0',NULL,NULL,NULL,_binary '20260709115200',1,NULL,0);
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `user_groups` WRITE;
/*!40000 ALTER TABLE `user_groups` DISABLE KEYS */;
INSERT INTO `user_groups` VALUES (1,_binary 'bureaucrat',NULL),(1,_binary 'interface-admin',NULL),(1,_binary 'staff',NULL),(1,_binary 'sysop',NULL),(5,_binary 'bureaucrat',NULL),(5,_binary 'interface-admin',NULL),(5,_binary 'sysop',NULL),(6,_binary 'bureaucrat',NULL),(6,_binary 'interface-admin',NULL),(6,_binary 'sysop',NULL),(7,_binary 'bureaucrat',NULL),(7,_binary 'interface-admin',NULL),(7,_binary 'sysop',NULL),(8,_binary 'interface-admin',NULL),(8,_binary 'sysop',NULL),(9,_binary 'interface-admin',NULL),(9,_binary 'sysop',NULL),(10,_binary 'bot',NULL),(10,_binary 'bureaucrat',NULL),(10,_binary 'interface-admin',NULL),(10,_binary 'sysop',NULL),(11,_binary 'seomaster',NULL),(11,_binary 'sysop',NULL),(12,_binary 'interface-admin',NULL),(12,_binary 'seomaster',NULL),(12,_binary 'staff',NULL),(12,_binary 'sysop',NULL),(12,_binary 'widgeteditor',NULL),(18,_binary 'bot',NULL),(18,_binary 'bureaucrat',NULL),(18,_binary 'seomaster',NULL),(18,_binary 'sysop',NULL),(20,_binary 'bot',NULL),(20,_binary 'smw-editor',NULL),(20,_binary 'sysop',NULL);
/*!40000 ALTER TABLE `user_groups` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `site_stats` WRITE;
/*!40000 ALTER TABLE `site_stats` DISABLE KEYS */;
INSERT INTO `site_stats` VALUES (1,19464,384,9527,22,1,926);
/*!40000 ALTER TABLE `site_stats` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `interwiki` WRITE;
/*!40000 ALTER TABLE `interwiki` DISABLE KEYS */;
INSERT INTO `interwiki` VALUES (_binary 'acronym',_binary 'https://www.acronymfinder.com/~/search/af.aspx?string=exact&Acronym=$1','','',0,0),(_binary 'advogato',_binary 'http://www.advogato.org/$1','','',0,0),(_binary 'arxiv',_binary 'https://www.arxiv.org/abs/$1','','',0,0),(_binary 'c2find',_binary 'http://c2.com/cgi/wiki?FindPage&value=$1','','',0,0),(_binary 'cache',_binary 'https://www.google.com/search?q=cache:$1','','',0,0),(_binary 'commons',_binary 'https://commons.wikimedia.org/wiki/$1',_binary 'https://commons.wikimedia.org/w/api.php','',0,0),(_binary 'dictionary',_binary 'http://www.dict.org/bin/Dict?Database=*&Form=Dict1&Strategy=*&Query=$1','','',0,0),(_binary 'doi',_binary 'https://dx.doi.org/$1','','',0,0),(_binary 'drumcorpswiki',_binary 'http://www.drumcorpswiki.com/$1',_binary 'http://drumcorpswiki.com/api.php','',0,0),(_binary 'dwjwiki',_binary 'http://www.suberic.net/cgi-bin/dwj/wiki.cgi?$1','','',0,0),(_binary 'elibre',_binary 'http://enciclopedia.us.es/index.php/$1',_binary 'http://enciclopedia.us.es/api.php','',0,0),(_binary 'emacswiki',_binary 'https://www.emacswiki.org/emacs/$1','','',0,0),(_binary 'foldoc',_binary 'https://foldoc.org/?$1','','',0,0),(_binary 'foxwiki',_binary 'https://fox.wikis.com/wc.dll?Wiki~$1','','',0,0),(_binary 'freebsdman',_binary 'https://www.FreeBSD.org/cgi/man.cgi?apropos=1&query=$1','','',0,0),(_binary 'gentoo-wiki',_binary 'http://gentoo-wiki.com/$1','','',0,0),(_binary 'google',_binary 'https://www.google.com/search?q=$1','','',0,0),(_binary 'googlegroups',_binary 'https://groups.google.com/groups?q=$1','','',0,0),(_binary 'hammondwiki',_binary 'http://www.dairiki.org/HammondWiki/$1','','',0,0),(_binary 'hrwiki',_binary 'http://www.hrwiki.org/wiki/$1',_binary 'http://www.hrwiki.org/w/api.php','',0,0),(_binary 'imdb',_binary 'http://www.imdb.com/find?q=$1&tt=on','','',0,0),(_binary 'kmwiki',_binary 'https://kmwiki.wikispaces.com/$1','','',0,0),(_binary 'linuxwiki',_binary 'http://linuxwiki.de/$1','','',0,0),(_binary 'lojban',_binary 'https://mw.lojban.org/papri/$1','','',0,0),(_binary 'lqwiki',_binary 'http://wiki.linuxquestions.org/wiki/$1','','',0,0),(_binary 'meatball',_binary 'http://meatballwiki.org/wiki/$1','','',0,0),(_binary 'mediawikiwiki',_binary 'https://www.mediawiki.org/wiki/$1',_binary 'https://www.mediawiki.org/w/api.php','',0,0),(_binary 'memoryalpha',_binary 'http://en.memory-alpha.org/wiki/$1',_binary 'http://en.memory-alpha.org/api.php','',0,0),(_binary 'metawiki',_binary 'http://sunir.org/apps/meta.pl?$1','','',0,0),(_binary 'metawikimedia',_binary 'https://meta.wikimedia.org/wiki/$1',_binary 'https://meta.wikimedia.org/w/api.php','',0,0),(_binary 'mozillawiki',_binary 'https://wiki.mozilla.org/$1',_binary 'https://wiki.mozilla.org/api.php','',0,0),(_binary 'mw',_binary 'https://www.mediawiki.org/wiki/$1',_binary 'https://www.mediawiki.org/w/api.php','',0,0),(_binary 'oeis',_binary 'https://oeis.org/$1','','',0,0),(_binary 'openwiki',_binary 'http://openwiki.com/ow.asp?$1','','',0,0),(_binary 'pmid',_binary 'https://www.ncbi.nlm.nih.gov/pubmed/$1?dopt=Abstract','','',0,0),(_binary 'pythoninfo',_binary 'https://wiki.python.org/moin/$1','','',0,0),(_binary 'rfc',_binary 'https://tools.ietf.org/html/rfc$1','','',0,0),(_binary 's23wiki',_binary 'http://s23.org/wiki/$1',_binary 'http://s23.org/w/api.php','',0,0),(_binary 'seattlewireless',_binary 'http://seattlewireless.net/$1','','',0,0),(_binary 'senseislibrary',_binary 'https://senseis.xmp.net/?$1','','',0,0),(_binary 'shoutwiki',_binary 'http://www.shoutwiki.com/wiki/$1',_binary 'http://www.shoutwiki.com/w/api.php','',0,0),(_binary 'squeak',_binary 'http://wiki.squeak.org/squeak/$1','','',0,0),(_binary 'theopedia',_binary 'https://www.theopedia.com/$1','','',0,0),(_binary 'tmbw',_binary 'http://www.tmbw.net/wiki/$1',_binary 'http://tmbw.net/wiki/api.php','',0,0),(_binary 'tmnet',_binary 'http://www.technomanifestos.net/?$1','','',0,0),(_binary 'twiki',_binary 'http://twiki.org/cgi-bin/view/$1','','',0,0),(_binary 'uncyclopedia',_binary 'https://en.uncyclopedia.co/wiki/$1',_binary 'https://en.uncyclopedia.co/w/api.php','',0,0),(_binary 'unreal',_binary 'https://wiki.beyondunreal.com/$1',_binary 'https://wiki.beyondunreal.com/w/api.php','',0,0),(_binary 'usemod',_binary 'http://www.usemod.com/cgi-bin/wiki.pl?$1','','',0,0),(_binary 'wiki',_binary 'http://c2.com/cgi/wiki?$1','','',0,0),(_binary 'wikia',_binary 'http://www.wikia.com/wiki/$1','','',0,0),(_binary 'wikibooks',_binary 'https://en.wikibooks.org/wiki/$1',_binary 'https://en.wikibooks.org/w/api.php','',0,0),(_binary 'wikidata',_binary 'https://www.wikidata.org/wiki/$1',_binary 'https://www.wikidata.org/w/api.php','',0,0),(_binary 'wikif1',_binary 'http://www.wikif1.org/$1','','',0,0),(_binary 'wikihow',_binary 'https://www.wikihow.com/$1',_binary 'https://www.wikihow.com/api.php','',0,0),(_binary 'wikimedia',_binary 'https://foundation.wikimedia.org/wiki/$1',_binary 'https://foundation.wikimedia.org/w/api.php','',0,0),(_binary 'wikinews',_binary 'https://en.wikinews.org/wiki/$1',_binary 'https://en.wikinews.org/w/api.php','',0,0),(_binary 'wikinfo',_binary 'http://wikinfo.co/English/index.php/$1','','',0,0),(_binary 'wikipedia',_binary 'https://en.wikipedia.org/wiki/$1',_binary 'https://en.wikipedia.org/w/api.php','',0,0),(_binary 'wikiquote',_binary 'https://en.wikiquote.org/wiki/$1',_binary 'https://en.wikiquote.org/w/api.php','',0,0),(_binary 'wikisource',_binary 'https://wikisource.org/wiki/$1',_binary 'https://wikisource.org/w/api.php','',0,0),(_binary 'wikispecies',_binary 'https://species.wikimedia.org/wiki/$1',_binary 'https://species.wikimedia.org/w/api.php','',0,0),(_binary 'wikiversity',_binary 'https://en.wikiversity.org/wiki/$1',_binary 'https://en.wikiversity.org/w/api.php','',0,0),(_binary 'wikivoyage',_binary 'https://en.wikivoyage.org/wiki/$1',_binary 'https://en.wikivoyage.org/w/api.php','',0,0),(_binary 'wikt',_binary 'https://en.wiktionary.org/wiki/$1',_binary 'https://en.wiktionary.org/w/api.php','',0,0),(_binary 'wiktionary',_binary 'https://en.wiktionary.org/wiki/$1',_binary 'https://en.wiktionary.org/w/api.php','',0,0);
/*!40000 ALTER TABLE `interwiki` ENABLE KEYS */;
UNLOCK TABLES;
DROP TABLE IF EXISTS `updatelog`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `updatelog` (
  `ul_key` varbinary(255) NOT NULL,
  `ul_value` blob,
  PRIMARY KEY (`ul_key`)
) ENGINE=InnoDB DEFAULT CHARSET=binary;
/*!40101 SET character_set_client = @saved_cs_client */;

LOCK TABLES `updatelog` WRITE;
/*!40000 ALTER TABLE `updatelog` DISABLE KEYS */;
INSERT INTO `updatelog` VALUES (_binary 'AddRFCandPMIDInterwiki',NULL),(_binary 'ConvertDjvuMetadata',NULL),(_binary 'DeduplicateArchiveRevId',NULL),(_binary 'DeleteDefaultMessages',NULL),(_binary 'FixDefaultJsonContentPages',NULL),(_binary 'FixInconsistentRedirects',NULL),(_binary 'MediaWiki\\Maintenance\\FixAutoblockLogTitles',NULL),(_binary 'MigrateActors',NULL),(_binary 'MigrateBlocks',NULL),(_binary 'MigrateComments',NULL),(_binary 'MigrateExternallinks',NULL),(_binary 'MigrateLinksTablepagelinks',NULL),(_binary 'MigrateLinksTabletemplatelinks',NULL),(_binary 'MigrateOldAJAXPollUserColumnsToActor',NULL),(_binary 'MigrateRevisionActorTemp',NULL),(_binary 'MigrateRevisionCommentTemp',NULL),(_binary 'PopulateChangeTagDef',NULL),(_binary 'PopulateContentTables',NULL),(_binary 'PopulateUserIsTemp',NULL),(_binary 'RefreshExternallinksIndex v1+IDN',NULL),(_binary 'UpdateCollation::uca-default-u-kn',NULL),(_binary 'actor-actor_name-patch-actor-actor_name-varbinary.sql',NULL),(_binary 'archive-ar_timestamp-dropDefault',NULL),(_binary 'archive-ar_title-patch-archive-ar_title-varbinary.sql',NULL),(_binary 'category-cat_title-patch-category-cat_title-varbinary.sql',NULL),(_binary 'categorylinks-cl_to-patch-categorylinks-cl_to-varbinary.sql',NULL),(_binary 'change_tag-ct_rc_id-patch-change_tag-ct_rc_id.sql',NULL),(_binary 'cleanup empty categories',NULL),(_binary 'content_models-model_id-patch-content_models-model_id.sql',NULL),(_binary 'dynamic-page-list-4-create-view',NULL),(_binary 'dynamic-page-list-4-delete-template',NULL),(_binary 'externallinks-el_index_60-dropDefault',NULL),(_binary 'externallinks-el_to-patch-externallinks-el_to_default.sql',NULL),(_binary 'filearchive-fa_deleted_timestamp-dropDefault',NULL),(_binary 'filearchive-fa_id-patch-filearchive-fa_id.sql',NULL),(_binary 'filearchive-fa_major_mime-patch-fa_major_mime-chemical.sql',NULL),(_binary 'filearchive-fa_name-patch-filearchive-fa_name.sql',NULL),(_binary 'filearchive-fa_size-patch-filearchive-fa_size_to_bigint.sql',NULL),(_binary 'filearchive-fa_timestamp-dropDefault',NULL),(_binary 'fix protocol-relative URLs in externallinks',NULL),(_binary 'image-img_major_mime-patch-image-img_major_mime-default.sql',NULL),(_binary 'image-img_major_mime-patch-img_major_mime-chemical.sql',NULL),(_binary 'image-img_name-patch-image-img_name-varbinary.sql',NULL),(_binary 'image-img_size-patch-image-img_size_to_bigint.sql',NULL),(_binary 'image-img_timestamp-dropDefault',NULL),(_binary 'image-img_timestamp-patch-image-img_timestamp.sql',NULL),(_binary 'imagelinks-il_to-patch-imagelinks-il_to-varbinary.sql',NULL),(_binary 'ip_changes-ipc_rev_timestamp-dropDefault',NULL),(_binary 'ipblocks-ipb_expiry-dropDefault',NULL),(_binary 'ipblocks-ipb_id-patch-ipblocks-ipb_id.sql',NULL),(_binary 'ipblocks-ipb_timestamp-dropDefault',NULL),(_binary 'ipblocks_restrictions-ir_ipb_id-patch-ipblocks_restrictions-ir_ipb_id.sql',NULL),(_binary 'ipblocks_restrictions-ir_type-patch-ipblocks_restrictions-ir_type.sql',NULL),(_binary 'ipblocks_restrictions-ir_value-patch-ipblocks_restrictions-ir_value.sql',NULL),(_binary 'iwlinks-iwl_prefix-patch-extend-iwlinks-iwl_prefix.sql',NULL),(_binary 'iwlinks-iwl_title-patch-iwlinks-iwl_title-varbinary.sql',NULL),(_binary 'job-job_timestamp-patch-job_job_timestamp.sql',NULL),(_binary 'job-job_title-patch-job-job_title-varbinary.sql',NULL),(_binary 'job-job_token_timestamp-patch-job_job_token_timestamp.sql',NULL),(_binary 'job-patch-job-params-mediumblob.sql',NULL),(_binary 'langlinks-ll_title-patch-langlinks-ll_title-varbinary.sql',NULL),(_binary 'linter-linter_params-/var/www/mediawiki/w/extensions/Linter/sql/mysql/patch-linter-fix-params-null-definition.sql',NULL),(_binary 'logging-log_title-patch-logging-log_title-varbinary.sql',NULL),(_binary 'migrate linter table linter_params data to the linter_tag and linter_template fields',NULL),(_binary 'migrate namespace id from page to linter table',NULL),(_binary 'objectcache-exptime-patch-objectcache-exptime-notnull.sql',NULL),(_binary 'oldimage-oi_major_mime-patch-oi_major_mime-chemical.sql',NULL),(_binary 'oldimage-oi_name-patch-oldimage-oi_name-varbinary.sql',NULL),(_binary 'oldimage-oi_size-patch-oldimage-oi_size_to_bigint.sql',NULL),(_binary 'oldimage-oi_timestamp-dropDefault',NULL),(_binary 'page-page_links_updated-patch-page-page_links_updated-noinfinite.sql',NULL),(_binary 'page-page_title-patch-page-page_title-varbinary.sql',NULL),(_binary 'page-page_touched-dropDefault',NULL),(_binary 'page_props-pp_page-patch-page_props-pp_page.sql',NULL),(_binary 'page_restrictions-pr_page-patch-page_restrictions-pr_page.sql',NULL),(_binary 'pagelinks-pl_title-patch-pagelinks-pl_title-varbinary.sql',NULL),(_binary 'populate *_from_namespace',NULL),(_binary 'populate externallinks.el_index_60',NULL),(_binary 'populate fa_sha1',NULL),(_binary 'populate img_sha1',NULL),(_binary 'populate ip_changes',NULL),(_binary 'populate pp_sortkey',NULL),(_binary 'populate rev_len and ar_len',NULL),(_binary 'populate rev_sha1',NULL),(_binary 'protected_titles-pt_expiry-dropDefault',NULL),(_binary 'protected_titles-pt_title-patch-protected_titles-pt_title-varbinary.sql',NULL),(_binary 'querycache-qc_title-patch-querycache-qc_title-varbinary.sql',NULL),(_binary 'querycachetwo-qcc_title-patch-querycachetwo-qcc_title-varbinary.sql',NULL),(_binary 'recentchanges-rc_id-patch-recentchanges-rc_id-bigint.sql',NULL),(_binary 'recentchanges-rc_id-patch-recentchanges-rc_id.sql',NULL),(_binary 'recentchanges-rc_timestamp-dropDefault',NULL),(_binary 'recentchanges-rc_timestamp-patch-recentchanges-rc_timestamp.sql',NULL),(_binary 'recentchanges-rc_title-patch-recentchanges-rc_title-varbinary.sql',NULL),(_binary 'redirect-rd_title-patch-redirect-rd_title-varbinary.sql',NULL),(_binary 'revision-rev_id-patch-revision-cleanup.sql',NULL),(_binary 'revision-rev_timestamp-dropDefault',NULL),(_binary 'searchindex-pk-titlelength',NULL),(_binary 'searchindex-tableoption-utf8mb4',NULL),(_binary 'site_stats-patch-site_stats-modify.sql',NULL),(_binary 'sites-site_global_key-patch-sites-site_global_key.sql',NULL),(_binary 'slot_roles-role_id-patch-slot_roles-role_id.sql',NULL),(_binary 'uploadstash-us_size-patch-uploadstash-us_size_to_bigint.sql',NULL),(_binary 'uploadstash-us_timestamp-patch-uploadstash-us_timestamp.sql',NULL),(_binary 'user-user_editcount-patch-user-user_editcount.sql',NULL),(_binary 'user-user_name-patch-user_table-updates.sql',NULL),(_binary 'user_former_groups-ufg_group-patch-ufg_group-length-increase-255.sql',NULL),(_binary 'user_groups-ug_group-patch-ug_group-length-increase-255.sql',NULL),(_binary 'user_newtalk-user_last_timestamp-patch-user_newtalk-user_last_timestamp-binary.sql',NULL),(_binary 'user_properties-up_property-patch-up_property.sql',NULL),(_binary 'watchlist-wl_notificationtimestamp-patch-watchlist-wl_notificationtimestamp.sql',NULL),(_binary 'watchlist-wl_title-patch-watchlist-wl_title-varbinary.sql',NULL);
/*!40000 ALTER TABLE `updatelog` ENABLE KEYS */;
UNLOCK TABLES;
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

LOCK TABLES `change_tag_def` WRITE;
/*!40000 ALTER TABLE `change_tag_def` DISABLE KEYS */;
INSERT INTO `change_tag_def` VALUES (1,_binary 'visualeditor-wikitext',0,166),(2,_binary 'mw-manual-revert',0,867),(3,_binary 'mw-reverted',0,470),(4,_binary 'wikieditor',0,4132),(5,_binary 'mw-undo',0,30),(6,_binary 'mw-new-redirect',0,21),(7,_binary 'mw-replace',0,14),(8,_binary 'mw-blank',0,15),(9,_binary 'mw-server-side-upload',0,332),(10,_binary 'discussiontools',0,2),(11,_binary 'discussiontools-visual',0,1),(12,_binary 'discussiontools-newtopic',0,2),(13,_binary 'mw-rollback',0,3),(14,_binary 'mw-changed-redirect-target',0,2),(15,_binary 'mw-removed-redirect',0,2),(16,_binary 'visualeditor',0,3),(17,_binary 'mw-contentmodelchange',0,2),(18,_binary 'discussiontools-source',0,1),(19,_binary 'discussiontools-source-enhanced',0,1),(20,_binary 'discussiontools-added-comment',0,1);
/*!40000 ALTER TABLE `change_tag_def` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

