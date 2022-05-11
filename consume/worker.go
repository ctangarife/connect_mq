package consume

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ok93-01-18/event_reporter"
	"github.com/streadway/amqp"
)

// event name
const CustomError = "custom-error"

type QueueRun struct {
	queue string
	run   bool
}

// TestSender is a simple sender for error reporter
type TestSender struct {
}

// Send is method for sending message
func (ts *TestSender) Send(ctx context.Context, subject string, msg string) error {
	fmt.Println(subject, msg)
	return nil
}

func Worker(queue []string) {
	forever := make(chan bool)
	conn := NewConnection("monitor", "monitoring", queue)
	if err := conn.Connect(); err != nil {
		panic(err)
	}
	if err := conn.BindQueue(); err != nil {
		panic(err)
	}
	deliveries, err := conn.Consume()
	if err != nil {
		panic(err)
	}
	for q, d := range deliveries {
		go conn.HandleConsumedDeliveries(q, d, messageHandler)
	}
	<-forever
}

func messageHandler(c Connection, q string, deliveries <-chan amqp.Delivery) {
	var wg sync.WaitGroup
	wg.Add(1)
	for d := range deliveries {
		m := Message{
			Queue:         q,
			Body:          MessageBody{Data: d.Body},
			ContentType:   d.ContentType,
			Priority:      d.Priority,
			CorrelationID: d.CorrelationId,
		}
		//handle the custom message
		message := string(m.Body.Data)

		log.Println("Got message from queue ", m.Queue, ": ", message)
		d.Ack(true)
		go buttonPanic(5, m.Queue, &wg)
		wg.Wait()
	}
}

func buttonPanic(timer int, queue string, wg *sync.WaitGroup) {
	reporter := event_reporter.New()
	err := reporter.Add(CustomError, &event_reporter.ReportConfig{
		Subject:   "Сustom error",                         // subject of message
		LogSize:   5,                                      // max count event execute before send
		ResetTime: 20 * time.Second,                       // reset MaxCount interval
		Senders:   []event_reporter.Sender{&TestSender{}}, // senders slice
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	defer wg.Done()

	for i := 0; i < timer; i++ {
		time.Sleep(1 * time.Second)
		if i == timer-1 {
			fmt.Printf("La cola no recibe mensajes hace %d segundos ", timer)
			reporter.Publish(CustomError, "error happened")
		}
	}
}
