// reverse proxy
package reverseproxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func Serve(terminal_port string, agent_port string) {
	terminal := NewReverseProxy("http://localhost:" + terminal_port)
	http.Handle("/term/", terminal)
	http.Handle("/health", http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("You lost your way, didn't you zoro?"))
		}))
	// http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("You lost your way, didn't you zoro?"))
	// }))
	log.Default().Println("Agent is listening on port " + agent_port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", agent_port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func NewReverseProxy(target string) *httputil.ReverseProxy {
	url, _ := url.Parse(target)
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = url.Scheme
			req.URL.Host = url.Host
		},
	}
}

// paas.vodafone.com/workspace/ws1/apps/terminal
