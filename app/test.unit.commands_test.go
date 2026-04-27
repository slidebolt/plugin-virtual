package app

import (
	"encoding/json"
	"testing"

	domain "github.com/slidebolt/sb-domain"
	messenger "github.com/slidebolt/sb-messenger-sdk"
)

func TestHandleCommand_MutatesCoreEntityState(t *testing.T) {
	_, store, _ := env(t)
	app := New()
	app.store = store

	tests := []struct {
		name   string
		entity domain.Entity
		cmd    any
		assert func(t *testing.T, got domain.Entity)
	}{
		{
			name: "light set rgbw",
			entity: domain.Entity{
				ID: "light1", Plugin: PluginID, DeviceID: "dev1", Type: "light",
				State: domain.Light{Power: false, Brightness: 10, Temperature: 300},
			},
			cmd: domain.LightSetRGBW{R: 1, G: 2, B: 3, W: 4, Brightness: 77},
			assert: func(t *testing.T, got domain.Entity) {
				light := got.State.(domain.Light)
				if !light.Power || light.ColorMode != "rgbw" || light.Brightness != 77 {
					t.Fatalf("light state = %+v", light)
				}
				if len(light.RGBW) != 4 || light.RGBW[3] != 4 {
					t.Fatalf("light RGBW = %v", light.RGBW)
				}
			},
		},
		{
			name: "light set hs",
			entity: domain.Entity{
				ID: "light2", Plugin: PluginID, DeviceID: "dev1", Type: "light",
				State: domain.Light{},
			},
			cmd: domain.LightSetHS{Hue: 120, Saturation: 50, Brightness: 88},
			assert: func(t *testing.T, got domain.Entity) {
				light := got.State.(domain.Light)
				if light.ColorMode != "hs" || len(light.HS) != 2 || light.Brightness != 88 {
					t.Fatalf("light state = %+v", light)
				}
			},
		},
		{
			name: "switch toggle",
			entity: domain.Entity{
				ID: "sw1", Plugin: PluginID, DeviceID: "dev1", Type: "switch",
				State: domain.Switch{Power: false},
			},
			cmd: domain.SwitchToggle{},
			assert: func(t *testing.T, got domain.Entity) {
				sw := got.State.(domain.Switch)
				if !sw.Power {
					t.Fatalf("switch state = %+v", sw)
				}
			},
		},
		{
			name: "button press increments",
			entity: domain.Entity{
				ID: "btn1", Plugin: PluginID, DeviceID: "dev1", Type: "button",
				State: domain.Button{Presses: 2},
			},
			cmd: domain.ButtonPress{},
			assert: func(t *testing.T, got domain.Entity) {
				btn := got.State.(domain.Button)
				if btn.Presses != 3 {
					t.Fatalf("button state = %+v", btn)
				}
			},
		},
		{
			name: "binary sensor turn on",
			entity: domain.Entity{
				ID: "bin1", Plugin: PluginID, DeviceID: "dev1", Type: "binary_sensor",
				State: domain.BinarySensor{On: false, DeviceClass: "occupancy"},
			},
			cmd: domain.BinarySensorTurnOn{},
			assert: func(t *testing.T, got domain.Entity) {
				sensor := got.State.(domain.BinarySensor)
				if !sensor.On || sensor.DeviceClass != "occupancy" {
					t.Fatalf("binary sensor state = %+v", sensor)
				}
			},
		},
		{
			name: "binary sensor turn off",
			entity: domain.Entity{
				ID: "bin2", Plugin: PluginID, DeviceID: "dev1", Type: "binary_sensor",
				State: domain.BinarySensor{On: true},
			},
			cmd: domain.BinarySensorTurnOff{},
			assert: func(t *testing.T, got domain.Entity) {
				sensor := got.State.(domain.BinarySensor)
				if sensor.On {
					t.Fatalf("binary sensor state = %+v", sensor)
				}
			},
		},
		{
			name: "climate set temperature",
			entity: domain.Entity{
				ID: "hvac1", Plugin: PluginID, DeviceID: "dev1", Type: "climate",
				State: domain.Climate{HVACMode: "cool", Temperature: 20},
			},
			cmd: domain.ClimateSetTemperature{Temperature: 23.5},
			assert: func(t *testing.T, got domain.Entity) {
				climate := got.State.(domain.Climate)
				if climate.Temperature != 23.5 {
					t.Fatalf("climate state = %+v", climate)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := store.Save(tc.entity); err != nil {
				t.Fatalf("save entity: %v", err)
			}
			addr := messenger.Address{Plugin: tc.entity.Plugin, DeviceID: tc.entity.DeviceID, EntityID: tc.entity.ID}
			if err := app.applyCommand(addr, tc.cmd); err != nil {
				t.Fatalf("applyCommand: %v", err)
			}

			raw, err := store.Get(domain.EntityKey{Plugin: tc.entity.Plugin, DeviceID: tc.entity.DeviceID, ID: tc.entity.ID})
			if err != nil {
				t.Fatalf("get entity: %v", err)
			}
			var got domain.Entity
			if err := json.Unmarshal(raw, &got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			tc.assert(t, got)
		})
	}
}
