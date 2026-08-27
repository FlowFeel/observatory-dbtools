Feature: Load-Time Cross-Check — unknown types fail loudly
  As a consumer of the catalog projection
  I want the loader to reject any property whose type PHP truth does not declare
  So that silent divergence is impossible at the moment of load

  # OWA applies to unknown fields, never to unknown types. An undecodable type
  # would make the auditor silently skip rows, so it is a hard failure — the
  # consumer refuses to proceed rather than guess.

  Scenario: A property with an unknown type is rejected at load
    Given a catalog artifact with a property of unknown type "Geo"
    Then loading the artifact must fail
    And the error must identify the property name and the unknown type

  Scenario: The inspection diagnoses the unknown type as a violation
    Given a catalog artifact with a property of unknown type "Geo"
    Then the discipline inspection must report an unknown-type violation

  Scenario: An unknown type is a hard failure, not an OWA warning
    Given a catalog artifact with a property of unknown type "Geo"
    Then the inspection must not classify the unknown type as an open-world warning

  Scenario: A clean artifact loads without a cross-check error
    Given a catalog artifact with version 1
    Then loading the artifact must succeed
