package main

import (
	"encoding/base64"
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
	router.HandleFunc("/auth", LocalVMAuth).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))

}

func HomePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Welcome to homepage, ci-cd with k8"))
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Hello World w/ kubectl ci-cd"))
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

func LocalVMAuth(w http.ResponseWriter, r *http.Request) {

	transport := &http.Transport{
		Proxy: http.ProxyURL(ProxyURL()),
	}

	//adding the Transport object to the http Client
	client := &http.Client{
		Transport: transport,
	}

	//generating the HTTP GET request
	request, err := http.NewRequest("GET", GetDestUrlLocalVMAuth().String(), nil)
	if err != nil {
		log.Println(err)
	}
	authorization := "Basic" + " " + basicauth("abc", "123")
	h := map[string][]string{
		"Authorization": {authorization},
	}
	request.Header = h
	//calling the URL
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Println(err)
	}

	//printing the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("GET on LocalVM with Basic Auth worked"))
}

func basicauth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func GetDestUrlLocalVM() *url.URL {

	opUrl, err := url.Parse("http://http-host:8001/j")
	if err != nil {
		log.Println(err)
	}
	return opUrl
}

func GetDestUrlLocalVMAuth() *url.URL {

	opUrl, err := url.Parse("http://http-host:8001/auth")
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
