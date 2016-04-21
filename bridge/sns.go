package bridge

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/techcampman/twitter-d-server/constant"
	"github.com/techcampman/twitter-d-server/logger"
)

type (
	// PushMessage interface for any message sender
	PushMessage interface {
		Send(message string, deviceType string, targetArn string) error
	}

	// MessageBySNS ... AWS SNS session structure
	MessageBySNS struct {
		*sns.SNS
	}

	iOSWrapper struct {
		APNS        string `json:"APNS"`
		APNSSandbox string `json:"APNS_SANDBOX"`
	}

	push struct {
		Alert *string     `json:"alert"`
		Sound *string     `json:"sound,omitempty"`
		Data  interface{} `json:"custom_data,omitempty"`
	}

	iosPush struct {
		APS push `json:"aps"`
	}
)

// NewPushMessageBySNS allocates and returns new SNS
func NewPushMessageBySNS(c *aws.Config) *MessageBySNS {
	return &MessageBySNS{
		SNS: sns.New(session.New(c)),
	}
}

// Send ... send notification message
func (s *MessageBySNS) Send(message string, deviceType string, targetArn string) (err error) {
	if targetArn == "" {
		return
	}
	if message == "" {
		return
	}
	if deviceType != constant.DeviceTypeiOS {
		err = s.sendPushNotification(message, targetArn)
	}

	return
}

func (s *MessageBySNS) sendPushNotification(message string, targetArn string) (err error) {
	data := new(push)
	data.Alert = &message
	msg := iOSWrapper{}
	ios := iosPush{
		APS: *data,
	}
	b, err := json.Marshal(ios)
	if err != nil {
		return err
	}
	msg.APNS = string(b[:])
	msg.APNSSandbox = string(b[:])
	pushData, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	m := string(pushData[:])
	arn := targetArn
	params := &sns.PublishInput{
		Message:          aws.String(m),
		MessageStructure: aws.String("json"),
		TargetArn:        aws.String(arn),
	}
	_, err = s.Publish(params)
	if err != nil {
		logger.Error(err)
	}

	return
}
