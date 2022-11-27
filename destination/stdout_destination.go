package destination

import (
	"encoding/json"
	"fmt"
	"log"
)

type StdoutDestination struct {
}

func (d *StdoutDestination) Publish(msg map[string]interface{}) {
	result, err := json.Marshal(msg)
	if err != nil {
		log.Panicf("%v\n\n%v", msg, err)
	}
	fmt.Println(string(result))
}
