// modified version of https://github.com/afiskon/promtail-client

package promtail

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const LOG_ENTRIES_CHAN_SIZE = 5000

type ClientConfig struct {

	// E.g. http://localhost:3100/api/prom/push
	PushURL string
	// E.g. "{job=\"somejob\"}"
	Labels             string
	BatchWait          time.Duration
	BatchEntriesNumber int
}

// http.Client wrapper for adding new methods, particularly sendJsonReq
type httpClient struct {
	parent http.Client
}

// A bit more convenient method for sending requests to the HTTP server
func (client *httpClient) sendJsonReq(method, url string, ctype string, reqBody []byte) (resp *http.Response, resBody []byte, err error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", ctype)

	resp, err = client.parent.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	resBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return resp, resBody, nil
}

type jsonLogEntry struct {
	Ts   time.Time `json:"ts"`
	Line string    `json:"line"`
}

type promtailStream struct {
	Labels  string          `json:"labels"`
	Entries []*jsonLogEntry `json:"entries"`
}

type promtailMsg struct {
	Streams []promtailStream `json:"streams"`
}

type JSONClient struct {
	config    *ClientConfig
	quit      chan struct{}
	entries   chan *jsonLogEntry
	waitGroup sync.WaitGroup
	client    httpClient
}

func NewClientJson(conf ClientConfig) (*JSONClient, error) {
	client := JSONClient{
		config:  &conf,
		quit:    make(chan struct{}),
		entries: make(chan *jsonLogEntry, LOG_ENTRIES_CHAN_SIZE),
		client:  httpClient{},
	}

	client.waitGroup.Add(1)
	go client.run()

	return &client, nil
}

func (c *JSONClient) Send(str string) {
	c.entries <- &jsonLogEntry{
		Ts:   time.Now(),
		Line: str,
	}
}

func (c *JSONClient) Shutdown() {
	close(c.quit)
	c.waitGroup.Wait()
}

func (c *JSONClient) run() {
	var batch []*jsonLogEntry
	batchSize := 0
	maxWait := time.NewTimer(c.config.BatchWait)

	defer func() {
		if batchSize > 0 {
			c.send(batch)
		}

		c.waitGroup.Done()
	}()

	for {
		select {
		case <-c.quit:
			return
		case entry := <-c.entries:
			batch = append(batch, entry)
			batchSize++
			if batchSize >= c.config.BatchEntriesNumber {
				c.send(batch)
				batch = []*jsonLogEntry{}
				batchSize = 0
				maxWait.Reset(c.config.BatchWait)
			}
		case <-maxWait.C:
			if batchSize > 0 {
				c.send(batch)
				batch = []*jsonLogEntry{}
				batchSize = 0
			}
			maxWait.Reset(c.config.BatchWait)
		}
	}
}

func (c *JSONClient) send(entries []*jsonLogEntry) {
	var streams []promtailStream
	streams = append(streams, promtailStream{
		Labels:  c.config.Labels,
		Entries: entries,
	})

	msg := promtailMsg{Streams: streams}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Printf("promtail.ClientJson: unable to marshal a JSON document: %s\n", err)
		return
	}

	resp, body, err := c.client.sendJsonReq("POST", c.config.PushURL, "application/json", jsonMsg)
	if err != nil {
		log.Printf("promtail.ClientJson: unable to send an HTTP request: %s\n", err)
		return
	}

	if resp.StatusCode != 204 {
		log.Printf("promtail.ClientJson: Unexpected HTTP status code: %d, message: %s\n", resp.StatusCode, body)
		return
	}
}
