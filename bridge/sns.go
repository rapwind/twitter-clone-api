package bridge

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/techcampman/twitter-d-server/constant"
	"github.com/techcampman/twitter-d-server/logger"
	"gopkg.in/mgo.v2/bson"
)

type (
	// PushMessage interface for any message sender
	PushMessage interface {
		Send(message string, count int, notifType string, notifID bson.ObjectId, deviceType string, targetArn string) error
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
		Alert            *string     `json:"alert"`
		Badge            int         `json:"badge,omitempty"`
		Sound            string      `json:"sound,omitempty"`
		ContentAvailable int         `json:"contentAvailable,omitempty"`
		Data             interface{} `json:"custom_data,omitempty"`
	}

	iosPush struct {
		APS  push          `json:"aps"`
		Type string        `json:"type"`
		ID   bson.ObjectId `json:"id"`
	}
)

// NewPushMessageBySNS allocates and returns new SNS
func NewPushMessageBySNS(c *aws.Config) *MessageBySNS {
	return &MessageBySNS{
		SNS: sns.New(session.New(c)),
	}
}

// Send ... send notification message
func (s *MessageBySNS) Send(message string, count int, notifType string, notifID bson.ObjectId, deviceType string, targetArn string) (err error) {
	if targetArn == "" {
		return
	}
	if message == "" {
		return
	}
	if deviceType == constant.DeviceTypeiOS {
		err = s.sendPushNotification(message, count, notifType, notifID, targetArn)
	}

	return
}

func (s *MessageBySNS) sendPushNotification(alert string, badge int, notifType string, notifID bson.ObjectId, targetArn string) (err error) {
	data := new(push)
	data.Alert = &alert
	data.Badge = badge
	data.ContentAvailable = 1
	data.Sound = "default"
	msg := iOSWrapper{}
	ios := iosPush{
		APS:  *data,
		Type: notifType,
		ID:   notifID,
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
