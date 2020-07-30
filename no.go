package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	accountSid = "ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	authToken  = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	domain     = "fcc.doesnt.care" // Change these if it's public facing/whatever

	scammerNumber = "+1202#######"
	twilioNumber  = "+1202#######"

	transportPrefix     = "http://"
	binDirectory        = "/root/regret/" // audio file directory
	audioDirectory      = "/audio"
	audioFullDirectory  = "/root/regret/audio"
	audioRoute          = "/audio/"
	callRoute           = "/call"
	recordCallbackRoute = "/rec"
	twimlRoute          = "/twiml"
	smsReceivedRoute    = "/sms"
)

type TwiML struct {
	XMLName xml.Name `xml:"Response"`

	Say  string `xml:",omitempty"`
	Play string `xml:",omitempty"`
}

func main() {
	fs := http.FileServer(http.Dir(audioFullDirectory))
	http.Handle(audioRoute, http.StripPrefix(audioDirectory, fs))
	http.HandleFunc(smsReceivedRoute, smsReceived)
	http.HandleFunc(twimlRoute, twiml)
	http.HandleFunc(callRoute, call)
	http.HandleFunc(recordCallbackRoute, rec)
	http.ListenAndServe(":80", nil)
}

func smsReceived(w http.ResponseWriter, r *http.Request) {
	fmt.Println(formatRequest(r))
}

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

func twiml(w http.ResponseWriter, r *http.Request) {
	twiml := `<?xml version="1.0" encoding="UTF-8"?>
<Response>
	<Say>Hello?</Say>
	<Pause length="2"/>
	<Say>Hi?</Say>
	<Pause length="1"/>
	<Say>Hello!</Say>
	<Pause length="1"/>
	<Play loop="30">` + transportPrefix + domain + audioRoute + `scam.mp3</Play>
</Response>`

	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(twiml))
}

func rec(w http.ResponseWriter, r *http.Request) {
	fmt.Println(formatRequest(r))
}

func call(w http.ResponseWriter, r *http.Request) {
	// Let's set some initial default variables
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Calls.json"

	// Build out the data for our message
	v := url.Values{}
	v.Set("Record", "true")
	v.Set("RecordingStatusCallback", transportPrefix+domain+recordCallbackRoute)
	v.Set("To", scammerNumber)
	v.Set("From", twilioNumber)
	v.Set("Url", transportPrefix+domain+twimlRoute)
	rb := *strings.NewReader(v.Encode())

	// Create Client
	client := &http.Client{}

	req, _ := http.NewRequest("POST", urlStr, &rb)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// make request
	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		err := json.Unmarshal(bodyBytes, &data)
		if err == nil {
			fmt.Println(data["sid"])
		}
	} else {
		fmt.Println(resp.Status)
		w.Write([]byte("Bad requestion : " + resp.Status))
	}
}
