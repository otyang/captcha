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

type HCaptcha struct{}

func (c *HCaptcha) IndexHandler(w http.ResponseWriter, r *http.Request) {

	htmlInString := strings.ReplaceAll(c.HTMLTemlplate(), "[SITE_KEY]", H_SITE_KEY)

	tpl, err := template.New("").Parse(htmlInString)
	if err != nil {
		log.Fatalln(err)
	}

	tpl.Execute(w, nil)
}

func (c *HCaptcha) SendHandler(w http.ResponseWriter, r *http.Request) {
	captcha := r.PostFormValue("h-captcha-response")
	fmt.Println(captcha)

	valid := c.CheckTheCaptcha(captcha)
	fmt.Println(valid)
	if valid {
		fmt.Fprintf(w, "The captcha was correct!")
	} else {
		fmt.Fprintf(w, "This captcha was NOT correct")
	}
}

func (c *HCaptcha) CheckTheCaptcha(response string) bool {
	req, _ := http.NewRequest("POST", "https://hcaptcha.com/siteverify", nil)
	q := req.URL.Query()
	q.Add("secret", H_API_KEY)
	q.Add("response", response)
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	var hCaptchaResponse map[string]interface{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &hCaptchaResponse)
	return hCaptchaResponse["success"].(bool)

}

func (c *HCaptcha) HTMLTemlplate() string {
	return html
}

var html = `
<html>

<head>
  <title>hCaptcha</title>
  <script src="https://js.hcaptcha.com/1/api.js" async defer></script>
</head>

<body>
  <form action="/send" method="POST">
    <input required type="text" name="email" placeholder="Email" />
    <br>
    <div class="h-captcha" data-sitekey="[SITE_KEY]"></div>
    <input type="submit" value="Submit">
  </form>
  <script src='https://www.google.com/recaptcha/api.js'></script>
</body>

</html>
`
