Feature: Semantic Value-Range Consistency and Historical Invariants
  As a database engineer
  I want to verify that stored property values satisfy declared constraints
  So that historical drift is detected and reported

  Scenario: Value outside the declared allowed set
    Given a compiled property catalog with property "Event type" and allowed values "In-Person, Virtual, Hybrid"
    And a database with smw_di_blob containing "Webinar" for that property
    Then an audit should report a range violation
    And the violation diagnostic should include the declared allowed values

  Scenario: Malformed date in smw_di_time
    Given a compiled property catalog with property "Event start date" declared as Date
    And a database with smw_di_time containing "August 27, 2026" for that property
    Then an audit should report a syntax violation

  Scenario: Dangling reference in smw_di_wikipage
    Given a compiled property catalog with a Page property "Event organizer"
    And a database with smw_di_wikipage o_id pointing to a non-existent smw_object_ids entry
    Then an audit should report a reference violation
