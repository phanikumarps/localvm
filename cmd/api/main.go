package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", HomePage).Methods("GET")
	router.HandleFunc("/hello", HelloWorld).Methods("GET")
	router.HandleFunc("/localvm", LocalVM).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))

}

func HomePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Welcome to homepage"))
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Hello World"))
}
func LocalVM(w http.ResponseWriter, r *http.Request) {

	transport := &http.Transport{
		Proxy: http.ProxyURL(ProxyURL()),
	}

	//adding the Transport object to the http Client
	client := &http.Client{
		Transport: transport,
	}

	//generating the HTTP GET request
	request, err := http.NewRequest("GET", GetDestUrlLocalVM().String(), nil)
	if err != nil {
		log.Println(err)
	}

	//calling the URL
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	//printing the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// w.Write(buf.Bytes())
	w.Write(data)
}
func GetDestUrlLocalVM() *url.URL {

	opUrl, err := url.Parse("http://http-host:8001/j")
	if err != nil {
		log.Println(err)
	}
	return opUrl
}

func ProxyURL() *url.URL {
	Url, err := url.Parse("http://connectivity-proxy.kyma-system.svc.cluster.local:20003")
	if err != nil {
		log.Println(err)
	}
	return Url
}
