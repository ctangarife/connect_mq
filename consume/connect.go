package consume

import (
	"errors"
	"fmt"
	"log"

	utl "github.com/ctangarife/connect_mq/utils"
	"github.com/streadway/amqp"
)

//MessageBody is the struct for the body passed in the AMQP message. The type will be set on the Request header
type MessageBody struct {
	Data []byte
}

//Message is the amqp request to publish
type Message struct {
	Queue         string
	ReplyTo       string
	ContentType   string
	CorrelationID string
	Priority      uint8
	Body          MessageBody
}

//Connection is the connection created
type Connection struct {
	name     string
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	queues   []string
	err      chan error
}

var (
	connectionPool = make(map[string]*Connection)
)

//NewConnection returns the new connection object
func NewConnection(name, exchange string, queues []string) *Connection {
	if c, ok := connectionPool[name]; ok {
		return c
	}

	c := &Connection{
		exchange: exchange,
		queues:   queues,
		err:      make(chan error),
	}
	connectionPool[name] = c
	return c
}

//GetConnection returns the connection which was instantiated
func GetConnection(name string) *Connection {
	return connectionPool[name]
}

func (c *Connection) Connect() error {
	var err error
	conf, errc := utl.ReadConf("config/config.yml")
	if errc != nil {
		log.Fatal(err)
	}
	var user string = conf.Credentials.User
	// Contraseña del usuario de RabbitMQ
	var pwd string = conf.Credentials.Pwd
	// La dirección IP de RabbitMQ Broker
	var host string = conf.Credentials.Host
	// Puerto monitoreado por RabbitMQ Broker
	var port string = conf.Credentials.Port

	//Conexión a RabbitMQ
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pwd, host, port)

	c.conn, err = amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("Error in creating rabbitmq connection with %s : %s", url, err.Error())
	}
	go func() {
		<-c.conn.NotifyClose(make(chan *amqp.Error)) //Listen to NotifyClose
		c.err <- errors.New("Connection Closed")
	}()
	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}
	if err := c.channel.ExchangeDeclare(
		c.exchange, // name
		"direct",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // noWait
		nil,        // arguments
	); err != nil {
		fmt.Println(err, "Error en la conexión")
		return fmt.Errorf("Error in Exchange Declare: %s", err)
	}
	return nil
}

func (c *Connection) BindQueue() error {
	for _, q := range c.queues {
		if _, err := c.channel.QueueDeclare(q, true, false, false, false, nil); err != nil {
			return fmt.Errorf("error in declaring the queue %s", err)
		}
		if err := c.channel.QueueBind(q, "my_routing_key", c.exchange, false, nil); err != nil {
			fmt.Println(err, "Error en la conexión BindQueue")
			return fmt.Errorf("Queue  Bind error: %s", err)
		}
	}
	return nil
}

//Reconnect reconnects the connection
func (c *Connection) Reconnect() error {
	if err := c.Connect(); err != nil {
		return err
	}
	if err := c.BindQueue(); err != nil {
		return err
	}
	return nil
}
