package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/streadway/amqp"
)

const (
	//AMQP URI
	uri = "amqp://guest:guest@localhost:5672/"
	//Durable AMQP queue name
	queueName = "rpc-queue"
)

//如果存在错误，则输出
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	//调用发布消息函数
	publish(uri, queueName)
}

//发布者的方法
//
//@amqpURI, amqp的地址
//@queue, queue的名称
func publish(amqpURI string, queue string) {
	//建立连接
	log.Printf("dialing %q", amqpURI)
	connection, err := amqp.Dial(amqpURI)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer connection.Close()

	//创建一个Channel
	log.Printf("got Connection, getting Channel")
	channel, err := connection.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	//创建一个queue
	log.Printf("got queue, declaring %q", queue)
	q, err := channel.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	//均衡处理，每次处理一条消息
	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	//订阅一个消息
	//log.Printf("Queue bound to Exchange, starting Consume")
	msgs, err := channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	// 发布消息
	go func() {
		for d := range msgs {
			n, err := strconv.Atoi(string(d.Body))
			failOnError(err, "Failed to convert body to integer")

			log.Printf(" [.] server端接收到的数据是 (%d)", n)
			response := n * 2

			err = channel.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte(strconv.Itoa(response)),
				})
			failOnError(err, "Failed to publish a message")

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Awaiting RPC requests")

	//没有写入数据，一直等待读，阻塞当前线程，目的是让线程不退出
	<-forever
}
