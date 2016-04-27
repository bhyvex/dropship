package dropship

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/consul/api"
)

type ConsulEventHook struct {
	HookMeta
	config *api.Config
}

func NewConsulEventHook() *ConsulEventHook {
	config := initializeConsulConfig(nil)
	return &ConsulEventHook{HookMeta{"consul-event"}, config}
}

func (h ConsulEventHook) Execute(config HookConfig, service Config) error {
	client, err := api.NewClient(h.config)
	if err != nil {
		return err
	}

	payload := map[string]string{
		"hash": service.Hash,
	}

	plBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	name, ok := config["name"]
	serv, ok := config["service"]
	tag, ok := config["tag"]
	node, ok := config["node"]

	if !ok {
		return errors.New("Consul Hook: invalid config")
	}

	id, meta, err := client.Event().Fire(&api.UserEvent{
		Name:          name,
		Payload:       plBytes,
		ServiceFilter: serv,
		TagFilter:     tag,
		NodeFilter:    node,
	}, nil)

	fmt.Println(id, meta, err)

	return err
}
