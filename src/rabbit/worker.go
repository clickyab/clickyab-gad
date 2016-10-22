package rabbit

import (
	"assert"
	"encoding/json"
	"reflect"
	"utils"

	"sync"

	"github.com/Sirupsen/logrus"
)

// goodFunc verifies that the function satisfies the signature, represented as a slice of types.
// The last type is the single result type; the others are the input types.
// A final type of nil means any result type is accepted.
func goodFunc(fn reflect.Value, types ...reflect.Type) bool {
	if fn.Kind() != reflect.Func {
		return false
	}
	// Last type is return, the rest are ins.
	if fn.Type().NumIn() != len(types)-1 || fn.Type().NumOut() != 1 {
		return false
	}
	for i := 0; i < len(types)-1; i++ {
		if fn.Type().In(i) != types[i] {
			return false
		}
	}
	outType := types[len(types)-1]
	if outType != nil && fn.Type().Out(0) != outType {
		return false
	}
	return true
}

// RunWorker listen on a topic in amqp
func RunWorker(exchange, topic, queue string, jobPattern interface{}, function interface{}, prefetch int, quit chan chan struct{}) error {
	in := reflect.ValueOf(jobPattern)

	fn := reflect.ValueOf(function)
	elemType := in.Type()

	var t bool
	if !goodFunc(fn, elemType, reflect.ValueOf(t).Type()) {
		logrus.Panic("function must be of type func(" + in.Type().Elem().String() + ") bool")
	}

	c, err := conn.Channel()
	if err != nil {
		return err
	}
	err = c.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	if err != nil {
		return err
	}

	q, err := c.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	// prefetch count
	// **WARNING**
	// If ignore this, then there is a problem with rabbit. prefetch all jobs for this worker then.
	// the next worker get nothing at all!
	// **WARNING**
	c.Qos(prefetch, 0, false)

	err = c.QueueBind(
		q.Name,   // queue name
		topic,    // routing key
		exchange, // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}
	consumerTag := <-utils.ID
	delivery, err := c.Consume(q.Name, consumerTag, false, false, false, false, nil)
	if err != nil {
		return err
	}

	waiter := sync.WaitGroup{}
bigLoop:
	for {
		select {
		case job := <-delivery:
			cp := reflect.New(reflect.TypeOf(jobPattern)).Elem().Addr().Interface()
			err := json.Unmarshal(job.Body, cp)
			if err != nil {
				assert.Nil(job.Reject(false))
			}
			input := []reflect.Value{reflect.ValueOf(cp).Elem()}
			waiter.Add(1)
			go func() {
				defer waiter.Done()

				out := fn.Call(input)
				if out[0].Interface().(bool) {
					assert.Nil(job.Ack(false))
				} else {
					assert.Nil(job.Nack(false, false))
				}
			}()
		case ok := <-quit:
			c.Cancel(consumerTag, false)
			waiter.Wait()
			FinalizeWait()
			ok <- struct{}{}
			break bigLoop
		}

	}

	return nil
}
