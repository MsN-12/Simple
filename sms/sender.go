package sms

import (
	"fmt"
	"github.com/ghasedakapi/ghasedak-go"
)

func SendSms(message string, receptor string) {
	c := ghasedak.NewClient("509fe8dba1c4d7674d798fb7540b89c443b0b8cf1f6e2ee22c97725dcd6a316d",
		"10008566")

	r := c.Send(message, receptor)
	fmt.Println(r.Code)
	fmt.Println(r.Message)

	//Param1 := 1337
	//r := c.SendOTP("09xxxxxxxxx", "Your Template", Param1)
	//fmt.Println(r.Message)
	//fmt.Println(r.Code)
}
