package counter

import (
	"RabbitExample/pkg/broker"
	"RabbitExample/services/counter/model"
	"encoding/json"
	"log"
)

type listener interface {
	Listen(eventName string, handler broker.EventHandler)
	Close()
}

type Counter struct {
	eventListener listener
}

func New(l listener) *Counter {
	c := Counter{
		eventListener: l,
	}
	c.run()

	return &c
}

func (c *Counter) Close() {
	c.eventListener.Close()
}

func (c *Counter) run() {
	log.Println("Counter service is running")
	c.eventListener.Listen("increment", c.increment)
}

func (c *Counter) increment(request []byte) ([]byte, error) {
	req := model.CountReq{}
	err := json.Unmarshal(request, &req)
	if err != nil {
		return []byte{}, err
	}

	log.Println("Incrementing count", req.Count)

	rsp := model.CountRsp{
		Count: req.Count + 1,
	}

	b, err := json.Marshal(rsp)
	if err != nil {
		return []byte{}, err
	}

	return b, nil
}
