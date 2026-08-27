Feature: Derivation Totality — every declared type is fully derivable
  As a database engineer
  I want the Go derivation tables (routing, codes, titles) to be total over PHP truth
  So that no declared type is silently skipped and the wire never decides

  # Totality has two directions:
  #   • Total over truth — every KnownType (the PHP PropertyType enum) maps to
  #     a storage table and an SMW code, and codes round-trip back to their type.
  #   • No extras — Go must not route a type PHP truth never declares; every
  #     scanned table must be reachable from a known type.
  #
  # This is the negotiation side of the projection line: the maps live in Go,
  # but they are provably complete and provably free of truth-divergent entries.

  Scenario: Type-to-table routing is total over known types
    Given the catalog known types
    Then every known type maps to a storage table

  Scenario: Type-to-code is total and round-trips
    Given the catalog known types
    Then every known type maps to an SMW code
    And every SMW code decodes back to its known type

  Scenario: No routed type is outside PHP truth
    Given the catalog known types
    Then no routed type is absent from the known set

  Scenario: The audit scan set equals the routable table set
    Given the catalog known types
    Then the scanned data-item tables equal the routable table set

  Scenario: Catalog titles are collision-free under SMW normalization
    Given a compiled property catalog is loaded from the fixture
    Then no two property names share an SMWTitle

  Scenario: A type outside the known set is not routable
    Given the catalog known types
    Then a type outside the known set must not be routable
