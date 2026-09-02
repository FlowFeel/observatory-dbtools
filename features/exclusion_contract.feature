Feature: Exclusion Contract — derived fields never cross the wire
  As a consumer of the catalog projection
  I want every field the Go side derives to be excluded from the artifact
  So that the wire carries facts, never decisions

  # The exclusion list is the projection's refusal surface. Every field the Go
  # auditor derives locally (routing table, SMW code, canonical title, sort
  # key, normalized form) must never be serialized. A smuggled derived field is
  # flagged as a violation, and — because consumers derive locally — it can
  # never influence negotiation.

  Scenario: A derived field on a property entry is never silently accepted
    Given a catalog artifact carrying the derived field "table" on a property entry
    Then the discipline inspection must report a derived-field violation

  Scenario: A derived field on an entity entry is flagged
    Given a catalog artifact carrying the derived field "smw_title" on an entity entry
    Then the discipline inspection must report a derived-field violation

  Scenario: A derived field at the top level is flagged
    Given a catalog artifact carrying the derived field "p_id" at the top level
    Then the discipline inspection must report a derived-field violation

  Scenario: The exclusion list covers every field the auditor derives
    Given the catalog known types
    Then the exclusion list covers the derived fields table, smw_code, smw_title, canonical_title, sortkey, p_id, normalized

  Scenario: A smuggled routing table never influences negotiation
    Given a catalog artifact carrying the derived field "table" on a property entry
    Then ExpectedTable still derives routing from the declared type
