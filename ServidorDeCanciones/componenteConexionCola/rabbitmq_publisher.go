package componenteconexioncola

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

var rabbitMQURL string

func init() {
	rabbitMQURL = os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://admin2:1234@192.168.56.1:5672/"
	}
}

func SetRabbitMQURL(url string) {
	if url != "" {
		rabbitMQURL = url
	}
}

type RabbitPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

// Representa el payload enviado a RabbitMQ
type NotificacionCancion struct {
	Titulo        string `json:"titulo"`
	Artista       string `json:"artista"`
	Album         string `json:"album"`
	Anio          int    `json:"anio"`
	Duracion      string `json:"duracion"`
	Genero        string `json:"genero"`
	FechaRegistro string `json:"fecha_registro"`
	Mensaje       string `json:"mensaje"`
}

// Crea la conexión y declara la cola.
// La URL de RabbitMQ se toma de la variable de entorno RABBITMQ_URL si está definida,
// si no, se usa una URL por defecto (puede ajustarse según el entorno).
func NewRabbitPublisher() (*RabbitPublisher, error) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("error conectando a RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("error abriendo canal: %v", err)
	}

	q, err := ch.QueueDeclare(
		"notificaciones_canciones",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("error declarando cola: %v", err)
	}

	return &RabbitPublisher{conn: conn, channel: ch, queue: q}, nil
}

func (p *RabbitPublisher) PublicarNotificacion(msg NotificacionCancion) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error convirtiendo mensaje a JSON: %v", err)
	}

	err = p.channel.Publish(
		"",
		p.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("error publicando mensaje: %v", err)
	}

	return nil
}

func (p *RabbitPublisher) Cerrar() {
	if p.channel != nil {
		_ = p.channel.Close()
	}
	if p.conn != nil {
		_ = p.conn.Close()
	}
}
