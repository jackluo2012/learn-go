package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"time"
)

//消息生产者

const (
	//AMQP URI
	uri = "amqp://guest:guest@localhost:5672/"
	//Durable AMQP exchange name
	exchangeName = ""
	//Durable AMQP queue name
	queueName = "test-idoall-queues"
	//Body of message
	bodyMsg string = "hello jackluo ."
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
	publish(uri, exchangeName, queueName, bodyMsg)
	log.Printf("published %dB OK", len(bodyMsg))
}

//发布者的方法
//
//@amqpURI, amqp的地址
//@exchange, exchange的名称
//@queue, queue的名称
//@body, 主体内容
func publish(amqpURI string, exchange string, queue string, body string) {
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

	log.Printf("got queue, declaring %q", queue)

	//创建一个queue
	q, err := channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	log.Printf("declared queue, publishing %dB body (%q)", len(body), body)

	// Producer只能发送到exchange，它是不能直接发送到queue的。
	// 现在我们使用默认的exchange（名字是空字符）。这个默认的exchange允许我们发送给指定的queue。
	// routing_key就是指定的queue名字。

	rand.Seed(time.Now().UnixNano())
	i := 1
	tick := time.NewTicker(time.Millisecond * time.Duration(rand.Intn(1000)))
	for {
		<-tick.C
		err = channel.Publish(
			exchange, // exchange
			q.Name,   // routing key
			false,    // mandatory
			false,    // immediate
			amqp.Publishing{
				Headers:         amqp.Table{},
				ContentType:     "text/plain",
				ContentEncoding: "",
				Body:            []byte(body + fmt.Sprintf("%d", i)),
			})
		i++
		if err != nil {
			failOnError(err, "Failed to publish a message")
		}
	}

}
