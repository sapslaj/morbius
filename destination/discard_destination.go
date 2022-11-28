package destination

import (
	"encoding/json"
	"log"
)

type DiscardDestinationConfig struct {
}

type DiscardDestination struct {
	Config *DiscardDestinationConfig
}

func NewDiscardDestination(config *DiscardDestinationConfig) DiscardDestination {
	if config == nil {
		config = &DiscardDestinationConfig{}
	}
	return DiscardDestination{
		Config: config,
	}
}

func (d *DiscardDestination) Publish(msg map[string]interface{}) {
	_, err := json.Marshal(msg)
	if err != nil {
		log.Panicf("%v\n\n%v", msg, err)
	}
}
