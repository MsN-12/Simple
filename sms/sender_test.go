package sms

import "testing"

func TestSendSms(t *testing.T) {
	msg := "hi this is test message for sms from ghasedak api"
	rec := "09390433026"

	SendSms(msg, rec)
}
