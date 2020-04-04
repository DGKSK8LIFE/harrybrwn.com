package auth

import (
	"encoding/json"
	"net/http"
)

// RedirectHandler handles Oauth 2.0 redirects.
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if query.Get("error") != "" {
		jsonerr := map[string]string{
			"error":       query.Get("error"),
			"description": query.Get("error_description"),
		}
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(jsonerr)
		return
	}

	var jsonResp = map[string]string{"error": "redirect server refused request"}
	if validate(r) {
		jsonResp = map[string]string{"code": query.Get("code")}
		w.WriteHeader(200)
	} else {
		w.WriteHeader(403)
	}
	json.NewEncoder(w).Encode(jsonResp)
}
