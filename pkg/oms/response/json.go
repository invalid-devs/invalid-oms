package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Data interface {
	JSON() []byte
}

func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	u, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Can't marshal to json: {}", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(u)
}
