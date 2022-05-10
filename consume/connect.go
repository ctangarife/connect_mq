package consume

import (
	"fmt"
	"log"

	utl "github.com/ctangarife/connect_mq/utils"
	"github.com/streadway/amqp"
)

// Función de conexión RabbitMQ
func RabbitMQConn() (conn *amqp.Connection, err error) {
	// Nombre de usuario asignado por RabbitMQ
	c, err := utl.ReadConf("config/credentials.yml")
	if err != nil {
		log.Fatal(err)
	}
	var user string = c.Credentials.User
	// Contraseña del usuario de RabbitMQ
	var pwd string = c.Credentials.Pwd
	// La dirección IP de RabbitMQ Broker
	var host string = c.Credentials.Host
	// Puerto monitoreado por RabbitMQ Broker
	var port string = c.Credentials.Port
	url := "amqp://" + user + ":" + pwd + "@" + host + ":" + port + "/?heartbeat=300"
	fmt.Println(url)
	// Crea una nueva conexión
	conn, err = amqp.Dial(url)
	// devuelve conexión y error
	return
}

// Función de manejo de errores
func ErrorHanding(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
