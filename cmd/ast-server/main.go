package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/prometheus/prometheus/promql"
)

func main() {
	listenAddr := flag.String("listen-addr", ":8000", "Web API listen address.")
	http.HandleFunc("/parse", func(w http.ResponseWriter, r *http.Request) {
		expr, err := promql.ParseExpr(r.FormValue("expr"))
		if err != nil {
			errJSON, err := json.Marshal(map[string]string{"type": "error", "message": fmt.Sprintf("Expression incomplete or buggy: %v", err)})
			if err != nil {
				http.Error(w, fmt.Sprintf("Error marshaling error JSON: %v", err), http.StatusInternalServerError)
				return
			}
			http.Error(w, string(errJSON), http.StatusBadRequest)
			return
		}
		buf, err := json.Marshal(translateAST(expr))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error marshaling AST: %v", err), http.StatusBadRequest)
			return
		}
		w.Write(buf)
	})
	http.ListenAndServe(*listenAddr, nil)
}
