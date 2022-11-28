package destination

import (
	"encoding/json"
	"fmt"
	"log"
)

type StdoutDestinationConfig struct {
}

type StdoutDestination struct {
	Config *StdoutDestinationConfig
}

func NewStdoutDestination(config *StdoutDestinationConfig) StdoutDestination {
	if config == nil {
		config = &StdoutDestinationConfig{}
	}
	return StdoutDestination{
		Config: config,
	}
}

func (d *StdoutDestination) Publish(msg map[string]interface{}) {
	result, err := json.Marshal(msg)
	if err != nil {
		log.Panicf("%v\n\n%v", msg, err)
	}
	fmt.Println(string(result))
}
