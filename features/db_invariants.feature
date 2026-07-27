Feature: MediaWiki Core and Extension Database Invariants
  As a database engineer for Observatory
  I want to ensure core MediaWiki tables and extension structures are installed
  So that the database layer remains stable and uncorrupted

  Scenario: Core MediaWiki and SMW schema invariant verification
    Given a clean baseline database is loaded
    Then the following core tables must exist:
      | user           |
      | page           |
      | revision       |
      | text           |
      | site_stats     |
      | smw_object_ids |
      | smw_fpt_mdat   |
      | smw_di_time    |
    And the total table count must be at least 100
