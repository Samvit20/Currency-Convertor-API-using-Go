package main
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"os"
)
type CurrencyConversionResponse struct {
	Success bool              `json:"success"`
	Rates   map[string]float64 `json:"rates"`
	Error   string            `json:"error"`
}

const currencyConvertorAPI = "https://free.currconv.com/api/v7/convert"

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/convert", handleConvert).Methods("GET")
	http.ListenAndServe(":8080", router)
}

func handleConvert(w http.ResponseWriter, r *http.Request) {
	
	params := r.URL.Query()
	from := params.Get("from")
	to := params.Get("to")
	amount := params.Get("amount")

	api_key := os.Getenv("API_KEY")
	
	response, err := http.Get(fmt.Sprintf("%s?q=%s_%s&compact=ultra&apiKey=%s", currencyConvertorAPI, from, to, api_key))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()
	
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var conversionResponse CurrencyConversionResponse
	err = json.Unmarshal(body, &conversionResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if !conversionResponse.Success {
		http.Error(w, conversionResponse.Error, http.StatusBadRequest)
		return
	}
	
	conversionRate := conversionResponse.Rates[from+"_"+to]
	convertedAmount := conversionRate * toFloat(amount)
	
	responseJSON := fmt.Sprintf(`{
		"from_currency": "%s",
		"to_currency": "%s",
		"amount": "%s",
		"converted_amount": "%f"
	}`, from, to, amount, convertedAmount)
	
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, responseJSON)
}
func toFloat(s string) float64 {
	result, _ := strconv.ParseFloat(s, 64)
	return result
}
