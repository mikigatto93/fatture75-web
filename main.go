package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type RequestHandler func(w http.ResponseWriter, r *http.Request)

var routerFunctions = make(map[string]RequestHandler)

func init() {
	routerFunctions["quote_data"] = handleQuoteDataRequest
	routerFunctions["fill_spreadsheet"] = handleFillSpredsheetRequest
}

func main() {

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/", fs)

	mux.HandleFunc("/api/", handleApiCall)

	fmt.Println("Listening on port 8888...")
	err := http.ListenAndServe(":8888", mux)
	log.Fatal(err)

}

func handleApiCall(w http.ResponseWriter, r *http.Request) {

	urlSegments := strings.Split(r.URL.Path, "/")
	lastSegment := urlSegments[len(urlSegments)-1]

	switch r.Method {
	case http.MethodGet:
		// Handle the GET request
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)

	case http.MethodPost:
		// Handle the POST request
		routeRequest(lastSegment, w, r)

	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)

	default:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

}

func routeRequest(segment string, w http.ResponseWriter, r *http.Request) {
	handler, ok := routerFunctions[segment]

	if ok {
		handler(w, r)
	} else {
		fmt.Printf("Request /api/%s has no handler!\n", segment)
	}

}

func handleQuoteDataRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20) // limit your max input length!

	file, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	name := strings.Split(header.Filename, ".")
	fmt.Printf("File name %s\n", name[0])

	var buf bytes.Buffer

	io.Copy(&buf, file)

	xmlCol := NewXmlCollector(buf)

	err = xmlCol.LoadData()
	if err != nil {
		panic(err)
	}
	fmt.Println(xmlCol.ProductData)
	buf.Reset()

	type JsonResponse struct {
		Title      string           `json:"title"`
		Commission string           `json:"commission"`
		Products   []XmlProductData `json:"products"`
	}

	response := JsonResponse{
		Products:   xmlCol.ProductData,
		Title:      xmlCol.HeaderData.Title,
		Commission: xmlCol.HeaderData.Commission,
	}

	encodedResponse, err := json.Marshal(response)

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(encodedResponse)
	}
}

func handleFillSpredsheetRequest(w http.ResponseWriter, r *http.Request) {

	collector := NewJsonCollector()

	err := collector.LoadData(r.Body)

	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(collector.Products)
	fmt.Println(collector.ExcelFileName)

	writer, err := NewExcelWriter("model.xlsm", "conversion_map.json")
	if err != nil {
		fmt.Print(err)
	}
	err = writer.InsertProducts(collector, collector.ExcelFileName)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("done")

}

/*https://stackoverflow.com/questions/24116147/how-to-download-file-in-browser-from-go-server*/
