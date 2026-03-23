Feature: Switch Entity
  # Source ref: contracts/switch.md

  Scenario: Create with default state
    Given a switch entity "plugin-virtual.dev1.sw001" named "Outlet" with power off
    When I retrieve "plugin-virtual.dev1.sw001"
    Then the entity type is "switch"
    And the switch power is off

  Scenario: State fields hydrate correctly
    Given a switch entity "plugin-virtual.dev1.sw002" named "Wall Switch" with power on
    When I retrieve "plugin-virtual.dev1.sw002"
    Then the switch power is on

  Scenario: Query by type
    Given a switch entity "plugin-virtual.dev1.sw003" named "Garage" with power off
    And a light entity "plugin-virtual.dev1.light001" named "Light" with power off
    When I query where "type" equals "switch"
    Then the results include "plugin-virtual.dev1.sw003"
    And the results do not include "plugin-virtual.dev1.light001"

  Scenario: Query by power state
    Given a switch entity "plugin-virtual.dev1.swOn" named "On Switch" with power on
    And a switch entity "plugin-virtual.dev1.swOff" named "Off Switch" with power off
    When I query where "type" equals "switch" and "state.power" equals "true"
    Then I get 1 result

  Scenario: Update is reflected on retrieval
    Given a switch entity "plugin-virtual.dev1.swUpd" named "Switch" with power off
    And I update switch "plugin-virtual.dev1.swUpd" to power on
    When I retrieve "plugin-virtual.dev1.swUpd"
    Then the switch power is on

  Scenario: Delete removes entity
    Given a switch entity "plugin-virtual.dev1.swDel" named "Switch" with power off
    When I delete "plugin-virtual.dev1.swDel"
    Then retrieving "plugin-virtual.dev1.swDel" should fail

  Scenario: switch_turn_on command updates state
    Given a switch entity "plugin-virtual.dev1.sw001" named "Outlet" with power off
    When I send "switch_turn_on" to "plugin-virtual.dev1.sw001"
    And I retrieve "plugin-virtual.dev1.sw001"
    Then the switch power is on

  Scenario: switch_turn_off command updates state
    Given a switch entity "plugin-virtual.dev1.sw001" named "Outlet" with power on
    When I send "switch_turn_off" to "plugin-virtual.dev1.sw001"
    And I retrieve "plugin-virtual.dev1.sw001"
    Then the switch power is off

  Scenario: switch_toggle command updates state
    Given a switch entity "plugin-virtual.dev1.sw001" named "Outlet" with power off
    When I send "switch_toggle" to "plugin-virtual.dev1.sw001"
    And I retrieve "plugin-virtual.dev1.sw001"
    Then the switch power is on

  Scenario: Raw payload decodes to canonical state
    When I decode a "switch" payload '{"power":true}'
    Then the switch power is on

  Scenario: switch_turn_on encodes to wire format
    When I encode "switch_turn_on" command with '{}'
    Then the wire payload field "state" equals "ON"

  Scenario: Raw discovery data is stored internally and hidden from queries
    Given a switch entity "plugin-virtual.dev1.sw001" named "Outlet" with power off
    And I write internal data for "plugin-virtual.dev1.sw001" with payload '{"commandTopic":"zigbee2mqtt/outlet/set"}'
    When I read internal data for "plugin-virtual.dev1.sw001"
    Then the internal data matches '{"commandTopic":"zigbee2mqtt/outlet/set"}'
    And querying type "switch" returns only state entities
