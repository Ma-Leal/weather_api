package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CEPRequest struct {
	CEP string `json:"cep"`
}

func main() {
	http.HandleFunc("/cep", handleCEPRequest)
	fmt.Println("Serviço A rodando em http://localhost:8081/cep")
	http.ListenAndServe(":8081", nil)
}

func handleCEPRequest(w http.ResponseWriter, r *http.Request) {
	var request CEPRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(request.CEP) != 8 {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	resp, err := http.Post("http://weather:8080/?cep="+request.CEP, "application/json", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	w.Write(body)

}
