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
	w.WriteHeader(statusCode)
	u, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Can't marshal to json: {}", err)
	}
	w.Write(u)
}
