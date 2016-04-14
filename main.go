// Router project main.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type jsonObjectType struct {
	Proxy ProxyType
}
type ProxyType struct {
	Meta         MetaType
	Http_headers HTTPHeadersType
	inSequence   InSequenceType
	outSequence  OutSequenceType
}

type MetaType struct {
	Name       string
	Version    string
	Proxy_type string
	Endpoint   string
}

type HTTPHeadersType struct {
	Set   string
	Unset string
}

type InSequenceType struct {
	sequenceName string
}

type OutSequenceType struct {
	sequneceName string
}

//Server Config
type ServerConfType struct {
	Server ServerType
}
type ServerType struct {
	Host_name string
	Port      string
	Context   string
}

func main() {
	fmt.Println("The Router")
	file, _ := readServerConfig()
	var serverConf ServerConfType
	json.Unmarshal(file, &serverConf)

	fmt.Printf("Server Context: %s\n", serverConf.Server.Context)
	fmt.Printf("Server Port: %s\n", fmt.Sprint(":", serverConf.Server.Port))

	http.HandleFunc(serverConf.Server.Context, handler)
	go http.ListenAndServe(fmt.Sprint(":", serverConf.Server.Port), nil)
	var input string
	fmt.Scanln(&input)

}

func handler(w http.ResponseWriter, r *http.Request) {

	file, _ := readRouteConfig()

	var proxyConf jsonObjectType
	json.Unmarshal(file, &proxyConf)

	var endPoint string

	endPoint = strings.TrimSpace(proxyConf.Proxy.Meta.Endpoint)

	var res *http.Response
	var err error

	if strings.Contains(r.Method, "GET") {
		res, err = http.Get(endPoint)
		if err != nil {
			log.Fatal(err)
		}
	} else if strings.Contains(r.Method, "POST") {
		res, err = http.Post(endPoint, r.Header.Get("Content-Type"), r.Body)
		if err != nil {
			log.Fatal(err)
		}
	}

	response, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%s", response)
	fmt.Fprintf(w, "%s", response)
}

func readRouteConfig() ([]byte, error) {
	file, e := ioutil.ReadFile("route-configs/header_proxy.conf")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	return file, e
}

func readServerConfig() ([]byte, error) {
	file, e := ioutil.ReadFile("server-configs/router.conf")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	return file, e
}
