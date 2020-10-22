package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/henryhiraki22/covid-api/core/domain"
	"github.com/nexmo-community/nexmo-go"
	_"github.com/nexmo-community/nexmo-go"
	"log"
	"net/http"
	"strconv"
	_"strconv"

)

const NEXMO_API_SECRET = "NEXMO_API_SECRET"
const NEXMO_API_KEY = "NEXMO_API_KEY"

func main(){
	err := handleRoutes()
	if err != nil{
		fmt.Println(err.Error())
	}
}
func handleRoutes() error{
	r := mux.NewRouter()
	r.HandleFunc("/healthz", healthCheck).Methods("GET")
	r.HandleFunc("/sendCases", sendCases).Methods("GET")
	r.HandleFunc("/sendSms", sendSms)
	err := http.ListenAndServe(":8080", r)
	if err != nil{
		fmt.Println("some errors has found")
	}
	return nil
}

func healthCheck(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
}

func getData() []byte{
	var brazilCases domain.BrazilCases
	resp, err := http.Get("https://coronavirus-19-api.herokuapp.com/countries/brazil")
	if err != nil {
		fmt.Print(err.Error())
	}
	if err != nil {
		fmt.Printf("something is wrong")
	}
	if err := json.NewDecoder(resp.Body).Decode(&brazilCases); err != nil {
		log.Println(err)
	}
	fmt.Println(resp.Body)
	data, err := json.Marshal(brazilCases)
	return data
}

func sendCases(w http.ResponseWriter, r *http.Request){
	request := getData()
	w.Header().Set("Content-Type", "application/json")
	w.Write(request)
}

func sendSms(w http.ResponseWriter, r* http.Request){
		var cases domain.BrazilCases
		reqCases := getData()

		err := json.Unmarshal(reqCases, &cases)

		auth := nexmo.NewAuthSet()
		auth.SetAPISecret(NEXMO_API_KEY, NEXMO_API_SECRET)

		client := nexmo.NewClient(http.DefaultClient, auth)
		smsReq := nexmo.SendSMSRequest {
			From: "",
			To:   "",
			Text: "Country:" + cases.Country + "Cases:" + strconv.Itoa(cases.NumberCases) + "Deaths: " + strconv.Itoa(cases.Deaths) + "TodayCases: " + strconv.Itoa(cases.TodayCases),
		}

		callR, _, err := client.SMS.SendSMS(smsReq)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Status:", callR.Messages[0].Status)
		fmt.Print(callR)

}
