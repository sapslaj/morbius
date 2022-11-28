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
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type ElasticseachDestination struct {
	client      *elasticsearch.Client
	index       string
	bulkIndexer esutil.BulkIndexer
}

func NewElasticsearchDestination() ElasticseachDestination {
	index := "netflow"
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		panic(err)
	}
	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  index,
		Client: client,
	})
	if err != nil {
		panic(err)
	}
	d := ElasticseachDestination{
		client:      client,
		index:       index,
		bulkIndexer: bulkIndexer,
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

	err = d.bulkIndexer.Add(context.Background(), esutil.BulkIndexerItem{
		Action:     "index",
		DocumentID: documentID,
		Body:       bytes.NewReader(data),
		OnFailure: func(ctx context.Context, bii esutil.BulkIndexerItem, biri esutil.BulkIndexerResponseItem, err error) {
			log.Printf("%v %v %v", bii, biri, err)
		},
	})
	if err != nil {
		panic(err)
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
