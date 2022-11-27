package destination

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type ElasticseachDestination struct {
	client *elasticsearch.Client
	index  string
}

func NewElasticsearchDestination() ElasticseachDestination {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		panic(err)
	}
	d := ElasticseachDestination{
		client: client,
		index:  "netflow",
	}
	d.setupIndex()
	return d
}

func (d *ElasticseachDestination) Publish(msg map[string]interface{}) {
	msg["@timestamp"] = fmt.Sprint(time.Now().UnixMilli())
	data, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	documentID := fmt.Sprintf("%x", sha256.Sum256(data))
	req := esapi.IndexRequest{
		Index:      d.index,
		DocumentID: documentID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}
	res, err := req.Do(context.TODO(), d.client)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if res.IsError() {
		panic(res.Status())
	}
}

func (d *ElasticseachDestination) setupIndex() {
	resp, err := d.client.Indices.Exists([]string{d.index})
	log.Print(resp)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode == 404 {
		resp, err = d.client.Indices.Create(d.index)
		log.Print(resp)
		if err != nil {
			panic(err)
		}
	}
	mappingsBody, _ := json.Marshal(map[string]interface{}{
		"properties": map[string]interface{}{
			"@timestamp": map[string]string{
				"type": "date",
			},
		},
	})
	resp, err = d.client.Indices.PutMapping([]string{d.index}, bytes.NewReader(mappingsBody))
	log.Print(resp)
	if err != nil {
		panic(err)
	}
}
