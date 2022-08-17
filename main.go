package main

import (
	"log"
	"net/http"
)

const (
	H_API_KEY  = "0x055a7D2c94f486D36e2d0B3110F72BB50d4cf52B" // HCaptcha Api Key
	H_SITE_KEY = "f53b9c2f-9ad5-47f5-869d-d3521d3ac4c8"       // HCaptcha Site Key

	G2_API_KEY  = "6LfjJ1IeAAAAAKcvTkZRVzmTmPw_hVIiKDzViMcT" // G2 Captcha ApiKey
	G2_SITE_KEY = "6LfjJ1IeAAAAAOM9508egaxProJTpvgaQlYeznrb" // G2 Site Key
)

type Captcha interface { // Ensures / Enforces Compliances
	IndexHandler(w http.ResponseWriter, r *http.Request)
	SendHandler(w http.ResponseWriter, r *http.Request)
	CheckTheCaptcha(response string) bool
	HTMLTemlplate() string
}

var _ Captcha = (*HCaptcha)(nil) // Interface Check
var _ Captcha = (*GoogleRecaptchaV2)(nil)

func main() {

	// HTML Server
	{
		htmlJs := http.NewServeMux()
		htmlJs.Handle("/", http.FileServer(http.Dir("./")))
		go func() {
			err := http.ListenAndServe(":9000", htmlJs)
			if err != nil {
				log.Fatal("ListenAndServe: Error Initialising Html Server", err)
			}
		}()
	}

	// HCaptcha Server
	{
		var hCaptchaStruct = &HCaptcha{}
		hServer := http.NewServeMux()
		hServer.HandleFunc("/", hCaptchaStruct.IndexHandler)
		hServer.HandleFunc("/send", hCaptchaStruct.SendHandler)
		go func() {
			err := http.ListenAndServe(":9001", hServer)
			if err != nil {
				log.Fatal("ListenAndServe: Error Initialising HCaptcha Server", err)
			}
		}()
	}

	// Recaptcha 2 Server
	{
		var g2CaptchaStruct = &GoogleRecaptchaV2{}
		g2Server := http.NewServeMux()
		g2Server.HandleFunc("/", g2CaptchaStruct.IndexHandler)
		g2Server.HandleFunc("/send", g2CaptchaStruct.SendHandler)

		err := http.ListenAndServe(":9002", g2Server)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}

}
