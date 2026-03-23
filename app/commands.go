package app

import (
	"encoding/json"
	"fmt"
	"log"

	domain "github.com/slidebolt/sb-domain"
	messenger "github.com/slidebolt/sb-messenger-sdk"
)

func (a *App) handleCommand(addr messenger.Address, cmd any) {
	if err := a.applyCommand(addr, cmd); err != nil {
		log.Printf("plugin-virtual: apply %T for %s: %v", cmd, addr.Key(), err)
	}
}

func (a *App) applyCommand(addr messenger.Address, cmd any) error {
	raw, err := a.store.Get(domain.EntityKey{
		Plugin:   addr.Plugin,
		DeviceID: addr.DeviceID,
		ID:       addr.EntityID,
	})
	if err != nil {
		return fmt.Errorf("load entity: %w", err)
	}

	var entity domain.Entity
	if err := json.Unmarshal(raw, &entity); err != nil {
		return fmt.Errorf("decode entity: %w", err)
	}

	switch c := cmd.(type) {
	case domain.LightTurnOn:
		return a.updateLight(entity, func(s *domain.Light) {
			s.Power = true
		})
	case domain.LightTurnOff:
		return a.updateLight(entity, func(s *domain.Light) {
			s.Power = false
		})
	case domain.LightSetBrightness:
		return a.updateLight(entity, func(s *domain.Light) {
			s.Power = true
			s.Brightness = c.Brightness
			if s.ColorMode == "" {
				s.ColorMode = "brightness"
			}
		})
	case domain.LightSetColorTemp:
		return a.updateLight(entity, func(s *domain.Light) {
			s.Power = true
			s.ColorMode = "color_temp"
			s.Temperature = c.Mireds
			if c.Brightness > 0 {
				s.Brightness = c.Brightness
			}
			clearLightColorState(s)
		})
	case domain.LightSetRGB:
		return a.updateLight(entity, func(s *domain.Light) {
			s.Power = true
			s.ColorMode = "rgb"
			s.RGB = []int{c.R, c.G, c.B}
			if c.Brightness > 0 {
				s.Brightness = c.Brightness
			}
			clearLightNonRGBState(s)
		})
	case domain.LightSetRGBW:
		return a.updateLight(entity, func(s *domain.Light) {
			s.Power = true
			s.ColorMode = "rgbw"
			s.RGBW = []int{c.R, c.G, c.B, c.W}
			if c.Brightness > 0 {
				s.Brightness = c.Brightness
			}
			clearLightNonRGBWState(s)
		})
	case domain.LightSetRGBWW:
		return a.updateLight(entity, func(s *domain.Light) {
			s.Power = true
			s.ColorMode = "rgbww"
			s.RGBWW = []int{c.R, c.G, c.B, c.CW, c.WW}
			if c.Brightness > 0 {
				s.Brightness = c.Brightness
			}
			clearLightNonRGBWWState(s)
		})
	case domain.LightSetHS:
		return a.updateLight(entity, func(s *domain.Light) {
			s.Power = true
			s.ColorMode = "hs"
			s.HS = []float64{c.Hue, c.Saturation}
			if c.Brightness > 0 {
				s.Brightness = c.Brightness
			}
			clearLightNonHSState(s)
		})
	case domain.LightSetXY:
		return a.updateLight(entity, func(s *domain.Light) {
			s.Power = true
			s.ColorMode = "xy"
			s.XY = []float64{c.X, c.Y}
			if c.Brightness > 0 {
				s.Brightness = c.Brightness
			}
			clearLightNonXYState(s)
		})
	case domain.LightSetWhite:
		return a.updateLight(entity, func(s *domain.Light) {
			s.Power = true
			s.ColorMode = "white"
			s.White = c.White
			clearLightWhiteState(s)
		})
	case domain.LightSetEffect:
		return a.updateLight(entity, func(s *domain.Light) {
			s.Power = true
			s.Effect = c.Effect
		})
	case domain.SwitchTurnOn:
		return a.updateSwitch(entity, func(s *domain.Switch) { s.Power = true })
	case domain.SwitchTurnOff:
		return a.updateSwitch(entity, func(s *domain.Switch) { s.Power = false })
	case domain.SwitchToggle:
		return a.updateSwitch(entity, func(s *domain.Switch) { s.Power = !s.Power })
	case domain.FanTurnOn:
		return a.updateFan(entity, func(s *domain.Fan) { s.Power = true })
	case domain.FanTurnOff:
		return a.updateFan(entity, func(s *domain.Fan) { s.Power = false })
	case domain.FanSetSpeed:
		return a.updateFan(entity, func(s *domain.Fan) {
			s.Power = c.Percentage > 0
			s.Percentage = c.Percentage
		})
	case domain.CoverOpen:
		return a.updateCover(entity, func(s *domain.Cover) { s.Position = 100 })
	case domain.CoverClose:
		return a.updateCover(entity, func(s *domain.Cover) { s.Position = 0 })
	case domain.CoverSetPosition:
		return a.updateCover(entity, func(s *domain.Cover) { s.Position = c.Position })
	case domain.LockLock:
		return a.updateLock(entity, func(s *domain.Lock) { s.Locked = true })
	case domain.LockUnlock:
		return a.updateLock(entity, func(s *domain.Lock) { s.Locked = false })
	case domain.ButtonPress:
		return a.updateButton(entity, func(s *domain.Button) { s.Presses++ })
	case domain.NumberSetValue:
		return a.updateNumber(entity, func(s *domain.Number) { s.Value = c.Value })
	case domain.SelectOption:
		return a.updateSelect(entity, func(s *domain.Select) { s.Option = c.Option })
	case domain.TextSetValue:
		return a.updateText(entity, func(s *domain.Text) { s.Value = c.Value })
	case domain.ClimateSetMode:
		return a.updateClimate(entity, func(s *domain.Climate) { s.HVACMode = c.HVACMode })
	case domain.ClimateSetTemperature:
		return a.updateClimate(entity, func(s *domain.Climate) { s.Temperature = c.Temperature })
	default:
		return fmt.Errorf("unknown command %T", cmd)
	}
}

func (a *App) updateLight(entity domain.Entity, mutate func(*domain.Light)) error {
	state, ok := entity.State.(domain.Light)
	if !ok {
		return fmt.Errorf("state type %T, want domain.Light", entity.State)
	}
	mutate(&state)
	entity.State = state
	return a.store.Save(entity)
}

func (a *App) updateSwitch(entity domain.Entity, mutate func(*domain.Switch)) error {
	state, ok := entity.State.(domain.Switch)
	if !ok {
		return fmt.Errorf("state type %T, want domain.Switch", entity.State)
	}
	mutate(&state)
	entity.State = state
	return a.store.Save(entity)
}

func (a *App) updateFan(entity domain.Entity, mutate func(*domain.Fan)) error {
	state, ok := entity.State.(domain.Fan)
	if !ok {
		return fmt.Errorf("state type %T, want domain.Fan", entity.State)
	}
	mutate(&state)
	entity.State = state
	return a.store.Save(entity)
}

func (a *App) updateCover(entity domain.Entity, mutate func(*domain.Cover)) error {
	state, ok := entity.State.(domain.Cover)
	if !ok {
		return fmt.Errorf("state type %T, want domain.Cover", entity.State)
	}
	mutate(&state)
	entity.State = state
	return a.store.Save(entity)
}

func (a *App) updateLock(entity domain.Entity, mutate func(*domain.Lock)) error {
	state, ok := entity.State.(domain.Lock)
	if !ok {
		return fmt.Errorf("state type %T, want domain.Lock", entity.State)
	}
	mutate(&state)
	entity.State = state
	return a.store.Save(entity)
}

func (a *App) updateButton(entity domain.Entity, mutate func(*domain.Button)) error {
	state, ok := entity.State.(domain.Button)
	if !ok {
		return fmt.Errorf("state type %T, want domain.Button", entity.State)
	}
	mutate(&state)
	entity.State = state
	return a.store.Save(entity)
}

func (a *App) updateNumber(entity domain.Entity, mutate func(*domain.Number)) error {
	state, ok := entity.State.(domain.Number)
	if !ok {
		return fmt.Errorf("state type %T, want domain.Number", entity.State)
	}
	mutate(&state)
	entity.State = state
	return a.store.Save(entity)
}

func (a *App) updateSelect(entity domain.Entity, mutate func(*domain.Select)) error {
	state, ok := entity.State.(domain.Select)
	if !ok {
		return fmt.Errorf("state type %T, want domain.Select", entity.State)
	}
	mutate(&state)
	entity.State = state
	return a.store.Save(entity)
}

func (a *App) updateText(entity domain.Entity, mutate func(*domain.Text)) error {
	state, ok := entity.State.(domain.Text)
	if !ok {
		return fmt.Errorf("state type %T, want domain.Text", entity.State)
	}
	mutate(&state)
	entity.State = state
	return a.store.Save(entity)
}

func (a *App) updateClimate(entity domain.Entity, mutate func(*domain.Climate)) error {
	state, ok := entity.State.(domain.Climate)
	if !ok {
		return fmt.Errorf("state type %T, want domain.Climate", entity.State)
	}
	mutate(&state)
	entity.State = state
	return a.store.Save(entity)
}

func clearLightColorState(s *domain.Light) {
	s.RGB = nil
	s.RGBW = nil
	s.RGBWW = nil
	s.HS = nil
	s.XY = nil
	s.White = 0
}

func clearLightNonRGBState(s *domain.Light) {
	s.RGBW = nil
	s.RGBWW = nil
	s.HS = nil
	s.XY = nil
	s.White = 0
}

func clearLightNonRGBWState(s *domain.Light) {
	s.RGB = nil
	s.RGBWW = nil
	s.HS = nil
	s.XY = nil
	s.White = 0
}

func clearLightNonRGBWWState(s *domain.Light) {
	s.RGB = nil
	s.RGBW = nil
	s.HS = nil
	s.XY = nil
	s.White = 0
}

func clearLightNonHSState(s *domain.Light) {
	s.RGB = nil
	s.RGBW = nil
	s.RGBWW = nil
	s.XY = nil
	s.White = 0
}

func clearLightNonXYState(s *domain.Light) {
	s.RGB = nil
	s.RGBW = nil
	s.RGBWW = nil
	s.HS = nil
	s.White = 0
}

func clearLightWhiteState(s *domain.Light) {
	s.RGB = nil
	s.RGBW = nil
	s.RGBWW = nil
	s.HS = nil
	s.XY = nil
}
