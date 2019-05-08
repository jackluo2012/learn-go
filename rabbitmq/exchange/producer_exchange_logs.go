package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"github.com/streadway/amqp"
)


const (
	//AMQP URI
	uri          =  "amqp://guest:guest@localhost:5672/"
	//Durable AMQP exchange name
	exchangeName =  "test-idoall-exchange-logs"
	//Exchange type - direct|fanout|topic|x-custom
	exchangeType = "fanout"
	//AMQP routing key
	routingKey   = ""
)

//如果存在错误，则输出
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main(){
	bodyMsg := bodyFrom(os.Args)
	//调用发布消息函数
	publish(uri, exchangeName, exchangeType, routingKey, bodyMsg)
	log.Printf("published %dB OK", len(bodyMsg))
}


func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello idoall.org"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}

//发布者的方法
//
//@amqpURI, amqp的地址
//@exchange, exchange的名称
//@exchangeType, exchangeType的类型direct|fanout|topic
//@routingKey, routingKey的名称
//@body, 主体内容
func publish(amqpURI string, exchange string, exchangeType string, routingKey string, body string){
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
	log.Printf("got Channel, declaring %q Exchange (%q)", exchangeType, exchange)
	err = channel.ExchangeDeclare(
		exchange,     // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// 发布消息
	log.Printf("declared queue, publishing %dB body (%q)", len(body), body)
	err = channel.Publish(
		exchange,     // exchange
		routingKey, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
			Headers:         amqp.Table{},
			ContentType: "text/plain",
			ContentEncoding: "",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
}