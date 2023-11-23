package sms

import (
	"fmt"
	"github.com/ghasedakapi/ghasedak-go"
)

func SendSms(message string, receptor string) {
	clientToken := "509fe8dba1c4d7674d798fb7540b89c443b0b8cf1f6e2ee22c97725dcd6a316d"
	clientID := "10008566"
	client := ghasedak.NewClient(clientToken, clientID)
	response := client.Send(message, receptor)
	fmt.Println(response.Code)
	fmt.Println(response.Message)

	//Param1 := 1337
	//r := c.SendOTP("09xxxxxxxxx", "Your Template", Param1)
	//fmt.Println(r.Message)
	//fmt.Println(r.Code)
}
