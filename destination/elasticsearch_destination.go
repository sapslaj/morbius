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

type ElasticseachDestinationConfig struct {
	Index          string
	TimestampField string
}

type ElasticseachDestination struct {
	Config      *ElasticseachDestinationConfig
	client      *elasticsearch.Client
	bulkIndexer esutil.BulkIndexer
}

func NewElasticsearchDestination(config *ElasticseachDestinationConfig) ElasticseachDestination {
	if config == nil {
		config = &ElasticseachDestinationConfig{}
	}
	if config.Index == "" {
		config.Index = "netflow"
	}
	if config.TimestampField == "" {
		config.TimestampField = "@timestamp"
	}
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		panic(err)
	}
	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  config.Index,
		Client: client,
	})
	if err != nil {
		panic(err)
	}
	d := ElasticseachDestination{
		Config:      config,
		client:      client,
		bulkIndexer: bulkIndexer,
	}
	d.setupIndex()
	return d
}

func (d *ElasticseachDestination) Publish(msg map[string]interface{}) {
	msg[d.Config.TimestampField] = fmt.Sprint(time.Now().UnixMilli())
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
	resp, err := d.client.Indices.Exists([]string{d.Config.Index})
	log.Print(resp)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode == 404 {
		resp, err = d.client.Indices.Create(d.Config.Index)
		log.Print(resp)
		if err != nil {
			panic(err)
		}
	}
	mappingsBody, _ := json.Marshal(map[string]interface{}{
		"properties": map[string]interface{}{
			d.Config.TimestampField: map[string]string{
				"type": "date",
			},
		},
	})
	resp, err = d.client.Indices.PutMapping([]string{d.Config.Index}, bytes.NewReader(mappingsBody))
	log.Print(resp)
	if err != nil {
		panic(err)
	}
}
