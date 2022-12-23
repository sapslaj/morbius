package destination

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/go-logfmt/logfmt"
)

type StdoutDestinationConfig struct {
	Format string `yaml:"format"`
}

type StdoutDestination struct {
	Config      *StdoutDestinationConfig
	publishFunc func(map[string]interface{}) (string, error)
	Writer      io.Writer
}

func NewStdoutDestination(config *StdoutDestinationConfig) StdoutDestination {
	if config == nil {
		config = &StdoutDestinationConfig{}
	}
	d := StdoutDestination{
		Config: config,
		Writer: os.Stdout, // Exposed so it can be overridden in tests
	}
	switch config.Format {
	case "logfmt":
		d.publishFunc = d.publishLogfmt
	case "json":
		fallthrough
	default:
		d.publishFunc = d.publishJSON
	}

	return d
}

func (d *StdoutDestination) Publish(msg map[string]interface{}) {
	result, err := d.publishFunc(msg)
	if err != nil {
		log.Panicf("%v\n\n%v", msg, err)
	}
	fmt.Fprintln(d.Writer, string(result))
}

func (d *StdoutDestination) publishJSON(msg map[string]interface{}) (string, error) {
	result, err := json.Marshal(msg)
	return string(result), err
}

func (d *StdoutDestination) publishLogfmt(msg map[string]interface{}) (string, error) {
	buf := &bytes.Buffer{}
	var err error
	encoder := logfmt.NewEncoder(buf)
	for key, value := range msg {
		err = encoder.EncodeKeyval(key, value)
		if err != nil {
			return buf.String(), err
		}
	}
	return buf.String(), err
}
