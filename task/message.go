package task

import (
	"github.com/mikitu/piesebeta_crawler/qutils"
	"encoding/gob"
	"github.com/streadway/amqp"
	"log"
	"bytes"
)

func NewMessagePublishModelTask(message []string, queueName string) *MessagePublishTask{
	msg := map[string]string{"year":message[0], "url": message[1], "model": message[2]}
	t := &MessagePublishTask{msg: msg, queueName: queueName}
	return t
}
type MessageModelsTask struct{
	Model string `json:"model"`
	Year string `json:"year"`
	Url string `json:"url"`
}
type MessagePublishTask struct{
	msg map[string]string
	queueName string
}

func (this MessagePublishTask) Execute() {
	conn, ch := qutils.GetChannel(qutils.Qurl)
	defer ch.Close()
	defer conn.Close()
	q := qutils.GetQueue(this.queueName, ch, false)
	msg := getMessage(this.msg)
	publish(ch, q, msg)
	log.Printf("Message sent: %v\n", this.msg["url"])
}

func getMessage(message interface{}) amqp.Publishing {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	enc.Encode(message)

	msg := amqp.Publishing{
		Body: buf.Bytes(),
	}
	return msg
}

func publish(ch *amqp.Channel, q *amqp.Queue, msg amqp.Publishing) {
	ch.Publish(
		"",             //exchange string,
		q.Name, 	//key string,
		false,          //mandatory bool,
		false,          //immediate bool,
		msg)            //msg amqp.Publishing)

}