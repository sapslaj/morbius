package destination

type Destination interface {
	Publish(map[string]interface{})
}
