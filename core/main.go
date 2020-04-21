package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nexmo-community/nexmo-go"
	"log"
	"net/http"
	"strconv"
)

const NEXMO_API_SECRET = "NEXMO_API_SECRET"
const NEXMO_API_KEY = "NEXMO_API_KEY"

type BrazilCases struct{
 	Country      string `json:"country"`
 	NumberCases  int `json:"cases"`
 	Deaths	 	 int `json:"deaths"`
	TodayCases 	 int `json:"todayCases"`
}

func main(){
	handleRoutes()
}

func handleRoutes(){
	r := mux.NewRouter()
	r.HandleFunc("/healthz", healthCheck).Methods("GET")
	r.HandleFunc("/getData", sendRequest).Methods("GET")
	err := http.ListenAndServe(":8080", r)
	if err != nil{
		fmt.Println("some errors has found")
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
}

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

	if err != nil {
		fmt.Print(err)
	}
	auth := nexmo.NewAuthSet()
	auth.SetAPISecret(NEXMO_API_KEY, NEXMO_API_SECRET)

	client := nexmo.NewClient(http.DefaultClient, auth)
	smsReq := nexmo.SendSMSRequest {
		From: "5513981281982",
		To:   "5513982002638",
		Text: "Country:" + brazilCases.Country + "Cases:" + strconv.Itoa(brazilCases.NumberCases) + "Deaths: " + strconv.Itoa(brazilCases.Deaths) + "TodayCases: " + strconv.Itoa(brazilCases.TodayCases),
	}

	callR, _, err := client.SMS.SendSMS(smsReq)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Status:", callR.Messages[0].Status)
	fmt.Print(callR)
}

