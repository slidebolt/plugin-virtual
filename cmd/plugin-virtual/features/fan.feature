Feature: Fan Entity
  # Source ref: contracts/fan.md

  Scenario: Create with default state
    Given a fan entity "plugin-virtual.dev1.fan001" named "Ceiling Fan" with power off and percentage 0
    When I retrieve "plugin-virtual.dev1.fan001"
    Then the entity type is "fan"
    And the fan power is off

  Scenario: State fields hydrate correctly
    Given a fan entity "plugin-virtual.dev1.fan002" named "Tower Fan" with power on and percentage 75
    When I retrieve "plugin-virtual.dev1.fan002"
    Then the fan power is on
    And the fan percentage is 75

  Scenario: Query by type
    Given a fan entity "plugin-virtual.dev1.fan003" named "Box Fan" with power off and percentage 0
    And a switch entity "plugin-virtual.dev1.sw001" named "Switch" with power off
    When I query where "type" equals "fan"
    Then the results include "plugin-virtual.dev1.fan003"
    And the results do not include "plugin-virtual.dev1.sw001"

  Scenario: Query by power state
    Given a fan entity "plugin-virtual.dev1.fanOn" named "Fan On" with power on and percentage 50
    And a fan entity "plugin-virtual.dev1.fanOff" named "Fan Off" with power off and percentage 0
    When I query where "type" equals "fan" and "state.power" equals "true"
    Then I get 1 result

  Scenario: Query by percentage
    Given a fan entity "plugin-virtual.dev1.fanFast" named "Fast" with power on and percentage 80
    And a fan entity "plugin-virtual.dev1.fanSlow" named "Slow" with power on and percentage 30
    When I query where "type" equals "fan" and "state.percentage" greater than 50
    Then I get 1 result

  Scenario: Update is reflected on retrieval
    Given a fan entity "plugin-virtual.dev1.fanUpd" named "Fan" with power off and percentage 0
    And I update fan "plugin-virtual.dev1.fanUpd" to power on and percentage 60
    When I retrieve "plugin-virtual.dev1.fanUpd"
    Then the fan power is on
    And the fan percentage is 60

  Scenario: Delete removes entity
    Given a fan entity "plugin-virtual.dev1.fanDel" named "Fan" with power off and percentage 0
    When I delete "plugin-virtual.dev1.fanDel"
    Then retrieving "plugin-virtual.dev1.fanDel" should fail

  Scenario: fan_turn_on command updates state
    Given a fan entity "plugin-virtual.dev1.fan001" named "Ceiling Fan" with power off and percentage 0
    When I send "fan_turn_on" to "plugin-virtual.dev1.fan001"
    And I retrieve "plugin-virtual.dev1.fan001"
    Then the fan power is on

  Scenario: fan_turn_off command updates state
    Given a fan entity "plugin-virtual.dev1.fan001" named "Ceiling Fan" with power on and percentage 75
    When I send "fan_turn_off" to "plugin-virtual.dev1.fan001"
    And I retrieve "plugin-virtual.dev1.fan001"
    Then the fan power is off

  Scenario: fan_set_speed command updates state
    Given a fan entity "plugin-virtual.dev1.fan001" named "Ceiling Fan" with power off and percentage 0
    When I send "fan_set_speed" with percentage 50 to "plugin-virtual.dev1.fan001"
    And I retrieve "plugin-virtual.dev1.fan001"
    Then the fan power is on
    And the fan percentage is 50

  Scenario: Raw payload decodes to canonical state
    When I decode a "fan" payload '{"power":true,"percentage":75}'
    Then the fan power is on
    And the fan percentage is 75

  Scenario: fan_set_speed encodes to wire format
    When I encode "fan_set_speed" command with '{"percentage":50}'
    Then the wire payload field "percentage" equals 50

  Scenario: Raw discovery data is stored internally and hidden from queries
    Given a fan entity "plugin-virtual.dev1.fan001" named "Ceiling Fan" with power off and percentage 0
    And I write internal data for "plugin-virtual.dev1.fan001" with payload '{"commandTopic":"zigbee2mqtt/fan/set","speedRange":{"min":0,"max":100}}'
    When I read internal data for "plugin-virtual.dev1.fan001"
    Then the internal data matches '{"commandTopic":"zigbee2mqtt/fan/set","speedRange":{"min":0,"max":100}}'
    And querying type "fan" returns only state entities
