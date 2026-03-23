package app

import (
	"fmt"

	domain "github.com/slidebolt/sb-domain"
)

func (a *App) seedDemo() error {
	entities := []domain.Entity{
		{
			ID: "demo_light", Plugin: PluginID, DeviceID: "demo_device",
			Type: "light", Name: "Demo Light",
			Commands: []string{"light_turn_on", "light_turn_off", "light_set_brightness", "light_set_color_temp"},
			State:    domain.Light{Power: false, Brightness: 128},
		},
		{
			ID: "demo_switch", Plugin: PluginID, DeviceID: "demo_device",
			Type: "switch", Name: "Demo Switch",
			Commands: []string{"switch_turn_on", "switch_turn_off", "switch_toggle"},
			State:    domain.Switch{Power: false},
		},
		{
			ID: "demo_sensor", Plugin: PluginID, DeviceID: "demo_device",
			Type: "sensor", Name: "Demo Temperature",
			State: domain.Sensor{Value: 21.0, Unit: "°C"},
		},
	}
	for _, entity := range entities {
		if err := a.store.Save(entity); err != nil {
			return fmt.Errorf("save %s: %w", entity.ID, err)
		}
	}
	return nil
}
