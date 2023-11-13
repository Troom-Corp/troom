package broker

import (
	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitBroker Rabbit

const NewPassword = "verifyPassword"
const NewEmail = "verifyEmail"

type Rabbit struct {
	Broker *amqp.Connection
}

// Connect открывает соединение с брокером
func (r Rabbit) Connect() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return conn, fiber.NewError(500, "Ошибка при подключении к брокеру")
	}
	r.Broker = conn
	return r.Broker, nil
}

// Close закрывает соединение с брокером
func (r Rabbit) Close() error {
	return r.Broker.Close()
}

// Send отправляет сообщение через брокера сервису EmailSender
func (r Rabbit) SendMsg(queueName string, data []byte) error {
	// queueName - имя очереди, разные названия для разных данных (password, phone, email)
	// data - строка JSON, куда кладется userid и email
	ch, err := r.Broker.Channel()
	if err != nil {
		return fiber.NewError(500, "Ошибка при смене пароля")
	}

	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fiber.NewError(500, "Ошибка при смене пароля")
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
	if err != nil {
		return fiber.NewError(500, "Ошибка при смене пароля")
	}
	return nil
}
