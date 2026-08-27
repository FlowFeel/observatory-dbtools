Feature: Projection Discipline — catalog.json carries declared facts only
  As a consumer of the catalog projection
  I want catalog.json to stay a minimal, flat, facts-only serialization of PHP truth
  So that derivation and negotiation live in consumers and the artifact never accretes decision state

  # The projection line:
  #   • Truth lives in the PHP semantic domain (manifests, PropertyType, the
  #     compiler). catalog.json is a projection surface, never a record of truth.
  #   • The artifact carries declared facts only: identity + type + allowed
  #     values + entity composition. Anything a consumer can derive locally is
  #     excluded — derivation is negotiation and lives in Go.
  #   • Provenance rides on the wire so a projection can be retranslated back
  #     to the truth-state it was projected from (Phase 2).
  #   • Unknown non-derived fields are tolerated per OWA — warnings, not failures.

  Scenario: The v1 artifact is flat and facts-only
    Given a catalog artifact with version 1
    Then the discipline inspection must pass
    And the artifact must not carry any derived field
    And the semantic nesting depth must not exceed the entry level

  Scenario: A derived field smuggled into the artifact is flagged
    Given a catalog artifact carrying the derived field "table" on a property entry
    Then the discipline inspection must report a derived-field violation

  Scenario: Unknown non-derived fields are tolerated under OWA
    Given a catalog artifact carrying the unknown top-level field "extensions"
    Then the artifact must still load successfully
    And the inspection must record an open-world warning, not a violation

  Scenario: Nested structures below the entry level are rejected
    Given a catalog artifact whose allowed values are nested objects
    Then the discipline inspection must report a nesting-depth violation

  Scenario: Routing is derived from declared type, never from the wire
    Given a catalog artifact declaring a Date property
    Then the expected storage table is derived from the type as smw_di_time
