package main

import (
	"log"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string)  {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
	}
}

func main()  {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Falha ao conectar ao RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Falha ao abrir Canal")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	failOnError(err, "Falha ao declarar a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Falha ao registrar consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs{
			log.Printf("Mensagem Recebida: %s", d.Body)
		}
	}()

	log.Printf("[*] Aguardando por mensagens. Para sair aperte CTRL+C")
	<-forever
}