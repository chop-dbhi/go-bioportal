package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/cors"
	"github.com/tylerb/graceful"
)

type rpcRequest struct {
	Method string
	Params map[string]interface{}
}

func readAll(ch <-chan *Class) []*Class {
	var a []*Class
	for c := range ch {
		a = append(a, c)
	}
	return a
}

func main() {
	var (
		httpaddr string
		tlscert  string
		tlskey   string
		neoaddr  string
	)

	flag.StringVar(&httpaddr, "http", "127.0.0.1:8080", "HTTP bind address.")
	flag.StringVar(&tlscert, "tlscert", "", "Path to TLS certificate.")
	flag.StringVar(&tlskey, "tlskey", "", "Path to TLS key.")
	flag.StringVar(&neoaddr, "neo4j", "127.0.0.1:7687", "Neo4j Bolt address.")

	flag.Parse()

	// Initialize the service.
	svc, err := NewService(neoaddr, -1)
	if err != nil {
		log.Fatalf("could not connect to the servce: %s", err)
	}

	// Top-level context that is cancellable by signals.
	cxt, cancel := context.WithCancel(context.Background())

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)

	go func() {
		select {
		case <-sig:
			cancel()
			return
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(405)
			return
		}

		var req rpcRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(422)
			return
		}

		var (
			resp interface{}
			err  error
		)

		switch req.Method {
		case "get":
			vocab := req.Params["vocab"].(string)
			code := req.Params["code"].(string)
			resp, err = svc.Get(vocab, code)

		case "validate":
			vocab := req.Params["vocab"].(string)
			var codes []string
			for _, v := range req.Params["codes"].([]interface{}) {
				codes = append(codes, v.(string))
			}
			resp, err = svc.Validate(vocab, codes)

		case "match":
			vocab := req.Params["vocab"].(string)
			pattern := req.Params["pattern"].(string)
			ch := make(chan *Class)
			go func() {
				err = svc.Match(cxt, vocab, pattern, ch)
				close(ch)
			}()
			resp = readAll(ch)

		case "parents":
			vocab := req.Params["vocab"].(string)
			code := req.Params["code"].(string)
			ch := make(chan *Class)
			go func() {
				err = svc.Parents(cxt, vocab, code, ch)
				close(ch)
			}()
			resp = readAll(ch)

		case "children":
			vocab := req.Params["vocab"].(string)
			code := req.Params["code"].(string)
			ch := make(chan *Class)
			go func() {
				err = svc.Children(cxt, vocab, code, ch)
				close(ch)
			}()
			resp = readAll(ch)

		case "ancestors":
			vocab := req.Params["vocab"].(string)
			code := req.Params["code"].(string)
			ch := make(chan *Class)
			go func() {
				err = svc.Ancestors(cxt, vocab, code, ch)
				close(ch)
			}()
			resp = readAll(ch)

		case "descendants":
			vocab := req.Params["vocab"].(string)
			code := req.Params["code"].(string)
			ch := make(chan *Class)
			go func() {
				err = svc.Descendants(cxt, vocab, code, ch)
				close(ch)
			}()
			resp = readAll(ch)

		case "flatten":
			vocab := req.Params["vocab"].(string)
			var codes []string
			for _, v := range req.Params["codes"].([]interface{}) {
				codes = append(codes, v.(string))
			}
			ch := make(chan *Class)
			go func() {
				err = svc.Flatten(cxt, vocab, codes, ch)
				close(ch)
			}()
			resp = readAll(ch)

		default:
			w.WriteHeader(501)
			return
		}

		w.Header().Set("content-type", "application/json")
		enc := json.NewEncoder(w)

		if err != nil {
			w.WriteHeader(500)
			enc.Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		enc.Encode(resp)
	})

	hdlr := cors.Default().Handler(mux)

	srv := &graceful.Server{
		Timeout: 10 * time.Second,
		Server: &http.Server{
			Addr:    httpaddr,
			Handler: hdlr,
		},
	}

	log.Printf("HTTP bind address %s", httpaddr)

	if tlscert == "" {
		err = srv.ListenAndServe()
	} else {
		err = srv.ListenAndServeTLS(tlscert, tlskey)
	}

	if err != nil {
		log.Fatal(err)
	}
}
