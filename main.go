package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Configuration struct {
	URLS []string
}

type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, r *http.Request)
}

type simpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}
type Loadbalancer struct {
	port            string
	roundRobinCount int
	servers         []Server
}

func NewLoadBalancer(port string, servers []Server) *Loadbalancer {
	return &Loadbalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}
}

func newSimpleServer(addr string) *simpleServer {

	serverUrl, err := url.Parse(addr)
	handleErr(err)
	return &simpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

func handleErr(err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

}

func (s *simpleServer) Address() string {
	return s.addr
}

func (s *simpleServer) IsAlive() bool {
	return true
}

func (s *simpleServer) Serve(rw http.ResponseWriter, r *http.Request) {

	s.proxy.ServeHTTP(rw, r)
}

func (lb *Loadbalancer) getNextAvaiableServer() Server {
	server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	for !server.IsAlive() {
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount++
	return server
}

func (lb *Loadbalancer) serverProxy(rw http.ResponseWriter, r *http.Request) {
	targetServer := lb.getNextAvaiableServer()
	fmt.Printf("forwarding request to address %q\n", targetServer.Address())
	targetServer.Serve(rw, r)

}

func main() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("List of backend servers\n%v", configuration.URLS)
	servers := []Server{
		newSimpleServer(configuration.URLS[0]),
		newSimpleServer(configuration.URLS[1]),
		newSimpleServer(configuration.URLS[2]),
	}
	lb := NewLoadBalancer("8000", servers)
	handleRedirect := func(rw http.ResponseWriter, r *http.Request) {
		lb.serverProxy(rw, r)
	}
	http.HandleFunc("/", handleRedirect)
	fmt.Println("\nServing request at 0.0.0.0:8000")
	http.ListenAndServe("0.0.0.0:"+lb.port, nil)
}
