package dolar

import (
	"encoding/json"
	"fmt"
	"github.com/thiagopereiramartinez/dolar/bcb"
	"net/http"
	"time"
)

func Dolar(w http.ResponseWriter, r *http.Request) {

	cotacao := make(chan bcb.CotacaoDolar)

	data := time.Now().Add(-4 * time.Hour)
	go bcb.ObterCotacao(data, cotacao)

	w.Header().Set("Content-Type", "application/json")

	if j, err := json.Marshal(<-cotacao); err != nil {
		fmt.Fprintf(w, "{'error':'General error'}")
		return
	} else {
		fmt.Fprintf(w, string(j))
	}

}
