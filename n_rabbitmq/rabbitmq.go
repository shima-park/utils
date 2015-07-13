package n_rabbitmq

import (
	"errors"
	"log"
	"net/url"
	"time"

	"github.com/streadway/amqp"
)

type Rabbitmq struct {
	Dial_url   string
	User       string
	Password   string
	Exchange   string
	RoutingKey string
	Vhost      string
	conn       *amqp.Connection
}

func (r *Rabbitmq) connect() (*amqp.Connection, error) {
	u, err := url.Parse(r.Dial_url)
	if err != nil {
		return nil, err
	}
	u.User = url.UserPassword(r.User, r.Password)
	u.Path = "/" + r.Vhost
	return amqp.Dial(u.String())
}

func (r *Rabbitmq) Dial() error {
	c := make(chan error, 1)
	go func() {
		if r.conn == nil {
			var err error
			r.conn, err = r.connect()
			c <- err
			if err != nil {
				return
			}
			select {}
		}
	}()
	err := <-c
	return err
}

func (r *Rabbitmq) Channel() (ch *amqp.Channel, err error) {
	if r.conn == nil {
		err = r.Reconn()
		if err != nil {
			return nil, err
		}
	}
	return r.conn.Channel()
}

func (r *Rabbitmq) Reconn() (err error) {
	conn, err := r.connect()
	if err == nil && conn != nil {
		r.conn = conn
	}
	return err
}

func (r *Rabbitmq) NotifyClose(c chan *amqp.Error) chan *amqp.Error {
	return r.conn.NotifyClose(c)
}

func (r *Rabbitmq) Publish(ch *amqp.Channel, msg *amqp.Publishing) (err error) {
	try_count := 3
	if ch == nil {
		log.Printf("[WARNING] send to mq failed: Channel must not be empty\n")
		return errors.New("ampq channel is nil")
	}

	defer func() {
		if f := recover(); f != nil {
			log.Printf("[WARNING] send to mq failed:%v\n", f)
			err = f.(error)
		}
	}()

	for i := 0; i < try_count; i++ {

		err = ch.Publish(r.Exchange, r.RoutingKey, false, false, *msg)

		if err != nil {
			conn_err := r.Reconn()
			ch1, _ := r.Channel()
			*ch = *ch1
			log.Printf("[WARNING] RabbitMQ Reconn err is %v send to bnow failed:%v try Count:%v\n", conn_err, err, i)
			time.Sleep(time.Second)
			continue
		}
		break
	}
	return
}

func (r *Rabbitmq) Close() {
	r.conn.Close()
}
