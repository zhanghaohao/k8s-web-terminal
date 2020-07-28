package ctx

import (
	"net/http"
	"encoding/json"
	"util/logger"
)

type CustomResponse struct {
	ReturnCode int 				`json:"return_code"`
	ReturnMessage interface{} 	`json:"return_message"`
	Data interface{}			`json:"data"`
}

func WriteJSON(w http.ResponseWriter, code int, message interface{}, data interface{})  {
	response := CustomResponse{ReturnCode: code, ReturnMessage: message, Data: data}
	ret, err := json.Marshal(response)
	if err != nil {
		logger.Error.Println(err)
		WriteJSON(w, 500, err.Error(), err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Accept", "application/json")
	w.Write(ret)
	return
}