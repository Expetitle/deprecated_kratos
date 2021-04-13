package queue

import (
	"encoding/json"
	"fmt"
	"github.com/ory/x/logrusx"
	"github.com/streadway/amqp"
	"github.com/tidwall/gjson"

	"github.com/ory/kratos/identity"
)

type (
	EmailUserData struct {
		Email         string `json:"email"`
		FirstName     string `json:"firstName"`
		LastName      string `json:"lastName"`
		Name          string `json:"name"`
		ConfirmSignup string `json:"confirmSignup"`
	}

	EmailUser struct {
		User EmailUserData `json:"user"`
	}

	NewAccountEmailData struct {
		EventType string    `json:"eventType"`
		Data      EmailUser `json:"data"`
	}
)

func SendVerificationQueue(identity *identity.Identity, address *identity.VerifiableAddress, rabbitMQURL string, verifyURL string, l *logrusx.Logger) {
	var firstName = gjson.GetBytes(identity.Traits, "name.first").String()
	var lastName = gjson.GetBytes(identity.Traits, "name.last").String()
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		l.Warn("Failed to connect to RabbitMQ")
		defer conn.Close()
	} else {
		ch, err := conn.Channel()
		if err != nil {
			l.Warn("Failed to open a channel")
			defer ch.Close()
		} else {
			fmt.Println("connected to RabbitMQ")
			q, err := ch.QueueDeclare("MAIL_NEW_ACCOUNT", false, false, false, false, nil)
			fmt.Println(q)
			// Handle any errors if we were unable to create the queue
			if err != nil {
				l.Warn("Failed to create a queue with name MAIL_NEW_ACCOUNT")
			} else {
				emailUserData := EmailUserData{Email: address.Value, FirstName: firstName, LastName: lastName, Name: firstName + " " + lastName, ConfirmSignup: verifyURL}
				emailUser := EmailUser{User: emailUserData}
				data := NewAccountEmailData{EventType: "NewAccount", Data: emailUser}
				stringifyData, _ := json.Marshal(data)
				fmt.Println(string(stringifyData))
				err = ch.Publish(
					"",
					"MAIL_NEW_ACCOUNT",
					false,
					false,
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte(string(stringifyData)),
					},
				)

				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("Successfully Published Message to Queue")
			}
		}
	}
}

func SendRecoveryQueue(identity *identity.Identity, address *identity.RecoveryAddress, rabbitMQURL string, recoveryURL string, l *logrusx.Logger) {
	var firstName = gjson.GetBytes(identity.Traits, "name.first").String()
	var lastName = gjson.GetBytes(identity.Traits, "name.last").String()
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		l.Warn("Failed to connect to RabbitMQ")
		defer conn.Close()
	} else {
		ch, err := conn.Channel()
		if err != nil {
			l.Warn("Failed to open a channel")
			defer ch.Close()
		} else {
			fmt.Println("connected to RabbitMQ")
			q, err := ch.QueueDeclare("RESET_PASSWORD_REQUEST", false, false, false, false, nil)
			fmt.Println(q)
			// Handle any errors if we were unable to create the queue
			if err != nil {
				l.Warn("Failed to create a queue with name RESET_PASSWORD_REQUEST")
			} else {
				emailUserData := EmailUserData{Email: address.Value, FirstName: firstName, LastName: lastName, Name: firstName + " " + lastName, ConfirmSignup: recoveryURL}
				emailUser := EmailUser{User: emailUserData}
				data := NewAccountEmailData{EventType: "resetPasswordRequest", Data: emailUser}
				stringifyData, _ := json.Marshal(data)
				fmt.Println(string(stringifyData))
				err = ch.Publish(
					"",
					"RESET_PASSWORD_REQUEST",
					false,
					false,
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte(string(stringifyData)),
					},
				)

				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("Successfully Published Message to Queue")
			}
		}
	}
}
