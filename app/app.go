package app

import (
	"encoding/json"
	"fmt"
	"log"

	contract "github.com/slidebolt/sb-contract"
	domain "github.com/slidebolt/sb-domain"
	messenger "github.com/slidebolt/sb-messenger-sdk"
	storage "github.com/slidebolt/sb-storage-sdk"
)

const PluginID = "plugin-virtual"

// App is the importable runtime for the plugin-virtual binary.
// Keep production behavior here so tests can exercise it without importing cmd/.
type App struct {
	msg   messenger.Messenger
	store storage.Storage
	cmds  *messenger.Commands
	subs  []messenger.Subscription
}

func New() *App {
	return &App{}
}

func (a *App) Hello() contract.HelloResponse {
	return contract.HelloResponse{
		ID:              PluginID,
		Kind:            contract.KindPlugin,
		ContractVersion: contract.ContractVersion,
		DependsOn:       []string{"messenger", "storage"},
	}
}

func (a *App) OnStart(deps map[string]json.RawMessage) (json.RawMessage, error) {
	msg, err := messenger.Connect(deps)
	if err != nil {
		return nil, fmt.Errorf("connect messenger: %w", err)
	}
	a.msg = msg

	store, err := storage.Connect(deps)
	if err != nil {
		return nil, fmt.Errorf("connect storage: %w", err)
	}
	a.store = store

	// Register plugin-specific types before wiring commands.
	// domain.Register("my_type", MyState{})
	// domain.RegisterCommand("my_action", MyCommand{})

	a.cmds = messenger.NewCommands(msg, domain.LookupCommand)
	sub, err := a.cmds.Receive(PluginID+".>", a.handleCommand)
	if err != nil {
		return nil, fmt.Errorf("subscribe commands: %w", err)
	}
	a.subs = append(a.subs, sub)
	if err := a.msg.Flush(); err != nil {
		return nil, fmt.Errorf("flush subscriptions: %w", err)
	}

	log.Println("plugin-virtual: started")
	return nil, nil
}

func (a *App) OnShutdown() error {
	for _, sub := range a.subs {
		sub.Unsubscribe()
	}
	if a.store != nil {
		a.store.Close()
	}
	if a.msg != nil {
		a.msg.Close()
	}
	return nil
}
