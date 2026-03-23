Feature: Number Entity
  # Source ref: contracts/number.md

  Scenario: Create with default state
    Given a number entity "plugin-virtual.dev1.num001" named "Volume" with value 50.0 min 0.0 max 100.0 step 1.0 unit ""
    When I retrieve "plugin-virtual.dev1.num001"
    Then the entity type is "number"
    And the number value is 50.0

  Scenario: State fields hydrate correctly
    Given a number entity "plugin-virtual.dev1.num002" named "Brightness Slider" with value 75.5 min 0.0 max 100.0 step 0.5 unit "%"
    When I retrieve "plugin-virtual.dev1.num002"
    Then the number value is 75.5

  Scenario: Query by type
    Given a number entity "plugin-virtual.dev1.num003" named "Speed" with value 30.0 min 0.0 max 100.0 step 5.0 unit "rpm"
    And a sensor entity "plugin-virtual.dev1.temp001" named "Temp" with value "20" and unit "°C"
    When I query where "type" equals "number"
    Then the results include "plugin-virtual.dev1.num003"
    And the results do not include "plugin-virtual.dev1.temp001"

  Scenario: Query by value
    Given a number entity "plugin-virtual.dev1.numHigh" named "High" with value 80.0 min 0.0 max 100.0 step 1.0 unit ""
    And a number entity "plugin-virtual.dev1.numLow" named "Low" with value 20.0 min 0.0 max 100.0 step 1.0 unit ""
    When I query where "type" equals "number" and "state.value" greater than 50
    Then I get 1 result

  Scenario: Update is reflected on retrieval
    Given a number entity "plugin-virtual.dev1.numUpd" named "Number" with value 10.0 min 0.0 max 100.0 step 1.0 unit ""
    And I update number "plugin-virtual.dev1.numUpd" to value 42.0
    When I retrieve "plugin-virtual.dev1.numUpd"
    Then the number value is 42.0

  Scenario: Delete removes entity
    Given a number entity "plugin-virtual.dev1.numDel" named "Number" with value 0.0 min 0.0 max 100.0 step 1.0 unit ""
    When I delete "plugin-virtual.dev1.numDel"
    Then retrieving "plugin-virtual.dev1.numDel" should fail

  Scenario: number_set_value command updates state
    Given a number entity "plugin-virtual.dev1.num001" named "Volume" with value 50.0 min 0.0 max 100.0 step 1.0 unit ""
    When I send "number_set_value" with value 77.0 to "plugin-virtual.dev1.num001"
    And I retrieve "plugin-virtual.dev1.num001"
    Then the number value is 77.0

  Scenario: Raw payload decodes to canonical state
    When I decode a "number" payload '{"value":42.0}'
    Then the number value is 42.0

  Scenario: number_set_value encodes to wire format
    When I encode "number_set_value" command with '{"value":42.0}'
    Then the wire payload field "value" equals 42.0

  Scenario: Raw discovery data is stored internally and hidden from queries
    Given a number entity "plugin-virtual.dev1.num001" named "Volume" with value 50.0 min 0.0 max 100.0 step 1.0 unit ""
    And I write internal data for "plugin-virtual.dev1.num001" with payload '{"commandTopic":"zigbee2mqtt/number/set","min":0,"max":100,"step":1}'
    When I read internal data for "plugin-virtual.dev1.num001"
    Then the internal data matches '{"commandTopic":"zigbee2mqtt/number/set","min":0,"max":100,"step":1}'
    And querying type "number" returns only state entities
