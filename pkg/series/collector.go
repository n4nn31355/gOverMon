package series

type ICollector interface {
	Collect()
}

type Collector struct {
	size int `validate:"gte=0"`
}

func (c *Collector) Collect() error {
	panic("not implemented")
}
