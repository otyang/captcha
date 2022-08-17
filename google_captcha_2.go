package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"
)

type GoogleRecaptchaV2 struct{}

func (c *GoogleRecaptchaV2) IndexHandler(w http.ResponseWriter, r *http.Request) {

	htmlInString := strings.ReplaceAll(c.HTMLTemlplate(), "[SITE_KEY]", G2_SITE_KEY)

	tpl, err := template.New("").Parse(htmlInString)
	if err != nil {
		log.Fatalln(err)
	}

	tpl.Execute(w, nil)
}

func (c *GoogleRecaptchaV2) SendHandler(w http.ResponseWriter, r *http.Request) {
	captcha := r.PostFormValue("g-recaptcha-response")
	fmt.Println(captcha)

	valid := c.CheckTheCaptcha(captcha)
	fmt.Println(valid)
	if valid {
		fmt.Fprintf(w, "The captcha was correct!")
	} else {
		fmt.Fprintf(w, "This captcha was NOT correct")
	}
}

func (c *GoogleRecaptchaV2) CheckTheCaptcha(response string) bool {
	req, _ := http.NewRequest("POST", "https://www.google.com/recaptcha/api/siteverify", nil)
	q := req.URL.Query()
	q.Add("secret", G2_API_KEY)
	q.Add("response", response)
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	var googleResponse map[string]interface{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &googleResponse)
	return googleResponse["success"].(bool)

}

func (c *GoogleRecaptchaV2) HTMLTemlplate() string {
	return recaptcha2Template
}

var recaptcha2Template = `
<html>

<head>
  <title>reCaptcha</title>
  <script src='https://www.google.com/recaptcha/api.js'></script>
</head>

<body>
  <form action="/send" method="POST">
    <input required type="text" name="email" placeholder="Email" />
    <br>
    <div name="jj" class="g-recaptcha" data-sitekey="[SITE_KEY]"></div>
    <input type="submit">
  </form>
  <script src='https://www.google.com/recaptcha/api.js'></script>
</body>

</html>
`
