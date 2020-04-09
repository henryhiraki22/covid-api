package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type BrazilCases struct{
 	Country  string `json:"country"`
 	NumberCases 	 int `json:"cases"`
 	Deaths	 int `json:"deaths"`
}

func main(){
	handleRoutes()
}

func handleRoutes(){
	r := mux.NewRouter()
	r.HandleFunc("/healthz", healthCheck).Methods("GET")
	r.HandleFunc("/getDate", sendRequest).Methods("GET")
	err := http.ListenAndServe(":8080", r)
	if err != nil{
		fmt.Println("some errors has found")
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`{"message": "i'm alive'}`))
	if err != nil{
		fmt.Println("error")
	}

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
}
