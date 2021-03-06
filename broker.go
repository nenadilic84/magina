package magina

import (
	"crypto/tls"
	"log"
	"net"

	"github.com/streadway/amqp"
)

// AuthenticateFunc AuthorizePublishFunc AuthorizeSubscribeFunc callback functions for authentication and authorization
type AuthenticateFunc func(client *Client, username string, password string) bool

// AuthorizePublishFunc callback functions for authentication and authorization
type AuthorizePublishFunc func(client *Client, topic string) bool

// AuthorizeSubscribeFunc callback functions for authentication and authorization
type AuthorizeSubscribeFunc func(client *Client, topic string) bool

// OnClientOnlineCB callback function when client connect success
type OnClientOnlineCB func(client *Client)

// OnClientOfflineCB callback function when losing heartbeat from client
type OnClientOfflineCB func(client *Client)

// Broker is MQTT main service
type Broker struct {
	// server address to listen
	Addr string
	// rabbit uri
	RabbitURI string
	// extend the broker to suport RPC. (WARNNING: NOT standard MQTT feature)
	SuportRPC bool
	// rabbitmq connection
	RabbitConnection *amqp.Connection
	// if use mqtts, set this
	TLSConfig *tls.Config
	// callbacks
	Authenticate       AuthenticateFunc
	AuthorizePublish   AuthorizePublishFunc
	AuthorizeSubscribe AuthorizeSubscribeFunc
	OnClientOnline     OnClientOnlineCB
	OnClientOffline    OnClientOfflineCB
}

// InitRabbitConn init rabbitmq connection.
func (b *Broker) InitRabbitConn() {
	if b.RabbitConnection == nil {
		conn, err := amqp.Dial(b.RabbitURI)
		failOnError(err)
		b.RabbitConnection = conn
	}
}

func (b *Broker) handleConnection(conn net.Conn) {
	client := &Client{
		Conn:   conn,
		Broker: b,
	}
	client.Serve()
}

// ListenAndServe serves for mqtt connections.
func (b *Broker) ListenAndServe() {
	b.InitRabbitConn()
	log.Println("listen and serve mqtt broker on " + b.Addr)
	var listener net.Listener
	var err error
	if b.TLSConfig != nil {
		listener, err = tls.Listen("tcp", b.Addr, b.TLSConfig)
	} else {
		listener, err = net.Listen("tcp", b.Addr)
	}
	failOnError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("error accepting new connection: " + err.Error())
		}
		go b.handleConnection(conn)
	}
}
