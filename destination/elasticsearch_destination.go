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
	Index               string
	TimestampField      string
	SynchronousIndexing bool
	ElasticsearchConfig *elasticsearch.Config
	BulkIndexerConfig   *esutil.BulkIndexerConfig
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
	if config.ElasticsearchConfig == nil {
		config.ElasticsearchConfig = &elasticsearch.Config{
			Addresses: []string{
				"http://127.0.0.1:9200",
			},
		}
	}
	client, err := elasticsearch.NewClient(*config.ElasticsearchConfig)
	if err != nil {
		panic(err)
	}
	d := ElasticseachDestination{
		Config: config,
		client: client,
	}
	if !d.Config.SynchronousIndexing {
		if d.Config.BulkIndexerConfig == nil {
			d.Config.BulkIndexerConfig = &esutil.BulkIndexerConfig{}
		}
		d.Config.BulkIndexerConfig.Index = d.Config.Index
		d.Config.BulkIndexerConfig.Client = d.client
		bulkIndexer, err := esutil.NewBulkIndexer(*d.Config.BulkIndexerConfig)
		if err != nil {
			panic(err)
		}
		d.bulkIndexer = bulkIndexer
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

	if d.Config.SynchronousIndexing {
		resp, err := d.client.Index(
			d.Config.Index,
			bytes.NewReader(data),
			d.client.Index.WithDocumentID(documentID),
		)
		if err != nil {
			log.Printf("%v %v", resp, err)
			return
		}
		if resp.IsError() {
			log.Printf("%v %v", resp, err)
			return
		}
	} else {
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
