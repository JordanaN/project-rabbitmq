package main

import (
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string)  {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
	}
}

func main()  {

	//faz a coneção com RabbitMQ
	conn, err := amqp.Dial("amqp://guest@localhost:5672")
	failOnError(err, "Falhou a conexão com RabbitMQ")
	defer conn.Close()

	//cria a conexão com o canal de mensagens
	ch, err := conn.Channel()
	failOnError(err, "Falhou a conexão com o canal")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Falha ao declarar a Queue")

	body := bodyFrom(os.Args)
	err = ch.Publish(
		"",           // exchange
		q.Name,       // routing key
		false,        // mandatory
		false,
		amqp.Publishing {
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	failOnError(err, "Falha ao publicar mensagem")
	log.Printf(" [x] Sent %s", body)
}

func bodyFrom(args []string) string  {
	var s string
	if (len(args)<2) || os.Args[1] == "" {
		s = "Hello"
	}else{
		s = strings.Join(args[1:], "")
	}
	return s
}