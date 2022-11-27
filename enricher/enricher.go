package enricher

type Enricher interface {
	Process(map[string]interface{}) map[string]interface{}
}
