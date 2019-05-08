package main

import (
	"fmt"
	"log"
	"github.com/streadway/amqp"
)

const (
	//AMQP URI
	uri           =  "amqp://guest:guest@localhost:5672/"
	//Durable AMQP exchange name
	exchangeName =  "test-idoall-exchange-logs"
	//Exchange type - direct|fanout|topic|x-custom
	exchangeType = "fanout"
	//AMQP binding key
	bindingKey   = ""
	//Durable AMQP queue name
	queueName     = ""
)

//如果存在错误，则输出
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main(){
	//调用消息接收者
	consumer(uri, exchangeName, exchangeType, queueName, bindingKey)
}

//接收者方法
//
//@amqpURI, amqp的地址
//@exchange, exchange的名称
//@exchangeType, exchangeType的类型direct|fanout|topic
//@queue, queue的名称
//@key , 绑定的key名称
func consumer(amqpURI string, exchange string, exchangeType string, queue string, key string){
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

	//创建一个exchange
	log.Printf("got Channel, declaring Exchange (%q)", exchange)
	err = channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	);
	failOnError(err, "Exchange Declare:")

	//创建一个queue
	q, err := channel.QueueDeclare(
		queueName, // name
		false,   // durable
		false,   // delete when unused
		true,   // exclusive 当Consumer关闭连接时，这个queue要被deleted
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	//绑定到exchange
	err = channel.QueueBind(
		q.Name, // name of the queue
		key,        // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	);
	failOnError(err, "Failed to bind a queue")

	log.Printf("Queue bound to Exchange, starting Consume")
	//订阅消息
	msgs, err := channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	//创建一个channel
	forever := make(chan bool)

	//调用gorountine
	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	//没有写入数据，一直等待读，阻塞当前线程，目的是让线程不退出
	<-forever
}