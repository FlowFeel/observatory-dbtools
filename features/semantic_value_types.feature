Feature: Semantic Value-Type Consistency and Table Routing
  As a database engineer
  I want to verify that stored values reside in the correct data-item table
  So that SMW routing integrity is maintained

  Scenario: Values in the correct table pass the audit
    Given a compiled property catalog
    And a database with a Date property stored in smw_di_time
    Then an audit should report zero routing violations

  Scenario: Historical type drift — Date values in smw_di_blob
    Given a compiled property catalog with property "Event start date" declared as Date
    And a database with smw_di_blob rows for that property's p_id
    Then an audit should report a routing violation
    And the diagnostic should include the expected table smw_di_time

  Scenario: Orphaned predicate — unknown p_id is not silently discarded
    Given a compiled property catalog
    And a database with smw_di_blob rows referencing an unknown p_id
    Then an audit should report an orphaned predicate
