package helpers

import (
	"encoding/json"
	"io"
	"net/http"
)

func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)

	return e.Encode(i)
}

func FromJSON(i interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(i)
}

func JSONError(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}
