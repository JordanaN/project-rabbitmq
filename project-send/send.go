package main

import (
"log"

"fmt"
"github.com/streadway/amqp"
)

func failOnError(err error, msg string)  {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
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
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Falha ao declarar a Queue")

	body := "Hello"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf("[x] Sent %s", body)
	failOnError(err, "Falha ao publicar a mensagem")
}