Feature: Semantic Property Declaration-Storage Consistency
  As a database engineer
  I want to verify that every catalog property is present in the SMW store
  So that the database reflects the declared semantic schema

  Scenario: All catalog properties exist with matching metadata
    Given a compiled property catalog
    And a database with SMW tables loaded
    And the catalog property pages are seeded into smw_object_ids
    Then every catalog property must have an smw_object_ids entry
    And the smw_fpt_type must match the declared type
    And the smw_fpt_pval must match the declared allowed values

  Scenario: Property missing from database is a declaration violation
    Given a compiled property catalog with property "Event type"
    And a database where "Event type" has no smw_object_ids entry
    Then an audit should report a declaration violation
    And the violation diagnostic should identify the property name
