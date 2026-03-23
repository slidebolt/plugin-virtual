package app

import (
	"encoding/json"
	"testing"

	domain "github.com/slidebolt/sb-domain"
	managersdk "github.com/slidebolt/sb-manager-sdk"
)

func TestOnStart_DoesNotSeedDemoEntities(t *testing.T) {
	env := managersdk.NewTestEnv(t)
	env.Start("messenger")
	env.Start("storage")

	deps := map[string]json.RawMessage{
		"messenger": env.MessengerPayload(),
	}

	app := New()
	if _, err := app.OnStart(deps); err != nil {
		t.Fatalf("OnStart: %v", err)
	}
	t.Cleanup(func() { _ = app.OnShutdown() })

	_, err := env.Storage().Get(domain.EntityKey{
		Plugin:   PluginID,
		DeviceID: "demo_device",
		ID:       "demo_light",
	})
	if err == nil {
		t.Fatal("expected no demo entities to be seeded")
	}
}
