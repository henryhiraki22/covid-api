package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"github.com/nexmo-community/nexmo-go"
)

const NEXMO_API_SECRET = "NEXMO_API_SECRET"
const NEXMO_API_KEY = "NEXMO_API_KEY"

type BrazilCases struct{
 	Country      string `json:"country"`
 	NumberCases  int `json:"cases"`
 	Deaths	 	 int `json:"deaths"`
	TodayCases 	 int `json:"todayCases"`
}

type SmsValue struct {
	From string `json:"from"`
	Text string `json:"text"`
	To   string `json:"to"`
	ApiKey  string `json:"apiKey"`
	ApiSecret string `json:"apiSecret"`
}

func main(){
	handleRoutes()
}

func handleRoutes(){
	r := mux.NewRouter()
	r.HandleFunc("/healthz", healthCheck).Methods("GET")
	r.HandleFunc("/getData", sendRequest).Methods("GET")
	r.HandleFunc("/sendSms", sendSms).Methods("POST")
	err := http.ListenAndServe(":8080", r)
	if err != nil{
		fmt.Println("some errors has found")
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
}

//TODO fazer o docker e subir no kubernets local.
func sendRequest(w http.ResponseWriter, r *http.Request){
	var brazilCases BrazilCases
	resp, err := http.Get("https://coronavirus-19-api.herokuapp.com/countries/brazil")
	if err != nil {
		fmt.Print(err.Error())
	}
	if err != nil{
		fmt.Printf("something is wrong")
	}
	if err := json.NewDecoder(resp.Body).Decode(&brazilCases); err != nil {
		log.Println(err)
	}
	data, err := json.Marshal(brazilCases)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

	err = sendSmsFunction()

	if err != nil {
		fmt.Print(err)
	}
}


func sendSms(w http.ResponseWriter, r *http.Request) {
	var content = SmsValue{
		From:      "test sender",
		Text:      "hello world",
		To:        "5513982002638",
	}
	data, err := json.Marshal(content)
	if err != nil {
		fmt.Print("something is wrong")
	}
	req, err := http.NewRequest("POST", "https://rest.nexmo.com/sms/json", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(NEXMO_API_KEY+":"+NEXMO_API_SECRET)))

	response, err := http.DefaultClient.Do(req)
	fmt.Print(response)
}

func sendSmsFunction() error{
	auth := nexmo.NewAuthSet()
	auth.SetAPISecret(NEXMO_API_KEY, NEXMO_API_SECRET)

	client := nexmo.NewClient(http.DefaultClient, auth)
	smsReq := nexmo.SendSMSRequest {
	From: "55139822002631",
	To: "5513988300054",
	Text: "papito é gay",
	}

	callR, _, err := client.SMS.SendSMS(smsReq)

	if err != nil {
	log.Fatal(err)
	}

	fmt.Println("Status:", callR.Messages[0].Status)
	fmt.Print(callR)

	return nil
}

