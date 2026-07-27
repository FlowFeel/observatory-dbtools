Feature: Semantic MediaWiki Property Storage Drift
  As a database administrator
  I want to detect and repair _MDAT routing mismatches between smw_fpt_mdat and smw_di_time
  So that semantic property queries return accurate timestamp data

  Scenario: Detecting and fixing SMW storage layer drift
    Given an SMW database with 10 FPT entries and 7 DI entries
    When a drift check is performed
    Then 3 missing DI entries should be detected
    When a drift fix is executed
    Then 3 rows should be inserted into smw_di_time
    And a subsequent drift check should detect zero drift
