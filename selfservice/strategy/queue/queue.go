package queue

import (
	"encoding/json"
	"github.com/ory/x/logrusx"
	"github.com/streadway/amqp"
	"github.com/tidwall/gjson"

	"github.com/ory/kratos/identity"
)

type (
	EmailUserData struct {
		Id            string `json:"id"`
		Email         string `json:"email"`
		FirstName     string `json:"firstName"`
		LastName      string `json:"lastName"`
		Name          string `json:"name"`
		ConfirmSignup string `json:"confirmSignup"`
	}

	PasswordRecoveryEmailData struct {
		Id          string `json:"id"`
		Email       string `json:"email"`
		FirstName   string `json:"firstName"`
		LastName    string `json:"lastName"`
		Name        string `json:"name"`
		RecoveryURL string `json:"recoveryUrl"`
	}

	EmailUser struct {
		User EmailUserData `json:"user"`
	}

	NewAccountEmailData struct {
		EventType string    `json:"eventType"`
		Data      EmailUser `json:"data"`
	}
)

func SendMessageToQueue(rabbitMQURL string, routing string, message string, l *logrusx.Logger) {
	// TODO: add better logging for which message, to where, was sent
	// TODO: the rabbitMQURL should not be passed as parameter, but taked from ENV directly from here
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
			l.Info("connected to RabbitMQ")
			// TODO: the exchange name should be in ENV variables/config
			err := ch.ExchangeDeclare("users_exchange", "direct", true, false, false, false, nil)
			// Handle any errors if we were unable to create the queue
			if err != nil {
				l.Warn("Failed to declare an exchange with name users_exchange")
			} else {
				l.Debug("Exchange declared")
				err = ch.Publish(
					"users_exchange",
					routing,
					false,
					false,
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte(message),
					},
				)

				if err != nil {
					l.Warn(err)
				}
				l.Info("Successfully Published Message to exchange")
			}
		}
	}
}

func SendVerificationQueue(identity *identity.Identity, address *identity.VerifiableAddress, rabbitMQURL string, verifyURL string, l *logrusx.Logger) {
	var firstName = gjson.GetBytes(identity.Traits, "name.first").String()
	var lastName = gjson.GetBytes(identity.Traits, "name.last").String()
	emailUserData := EmailUserData{Id: identity.ID.String(), Email: address.Value, FirstName: firstName,
		LastName: lastName, Name: firstName + " " + lastName, ConfirmSignup: verifyURL}
	emailUser := EmailUser{User: emailUserData}
	data := NewAccountEmailData{EventType: "NewAccount", Data: emailUser}
	stringifyData, _ := json.Marshal(data)
	l.Info(string(stringifyData))
	// TODO: the routing name should be default and in ENV variables/config
	SendMessageToQueue(rabbitMQURL, "MAIL_NEW_ACCOUNT", string(stringifyData), l)
}

func SendRecoveryQueue(identity *identity.Identity, address *identity.RecoveryAddress, rabbitMQURL string, recoveryURL string, l *logrusx.Logger) {
	var firstName = gjson.GetBytes(identity.Traits, "name.first").String()
	var lastName = gjson.GetBytes(identity.Traits, "name.last").String()
	passwordRecoveryEmailData := PasswordRecoveryEmailData{Id: identity.ID.String(), Email: address.Value, FirstName: firstName,
		LastName: lastName, Name: firstName + " " + lastName, RecoveryURL: recoveryURL}
	stringifyData, _ := json.Marshal(passwordRecoveryEmailData)
	l.Info(string(stringifyData))
	// TODO: the routing name should be default and in ENV variables/config
	SendMessageToQueue(rabbitMQURL, "RESET_PASSWORD_REQUEST", string(stringifyData), l)
}
