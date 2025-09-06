package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)


func WriteJson(o any, w http.ResponseWriter) error {
	js, er := json.Marshal(o)
	if er != nil {
		return er
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return nil
}

func ReadJSON(dst any, r *http.Request) error {

	er := json.NewDecoder(r.Body).Decode(&dst)

	if er != nil {

		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshall *json.InvalidUnmarshalError

		switch {
		case errors.As(er, &syntaxError):
			{
				return fmt.Errorf("Json is not proper at %d", syntaxError.Offset)
			}
		case errors.Is(er, io.ErrUnexpectedEOF):
			{
				return fmt.Errorf("Json is not proper")
			}
		case errors.As(er, &unmarshalTypeError):
			{
				fmt.Println("unmarshalTypeError", unmarshalTypeError.Field)
				if unmarshalTypeError.Field != "" {
					return fmt.Errorf("Incorrect json type for field [%s]", unmarshalTypeError.Field)
				}
				return fmt.Errorf("Incorrect json type at [%d]", unmarshalTypeError.Offset)
			}
		case errors.Is(er, io.EOF):
			{
				return fmt.Errorf("Body cannot be empty")
			}
		case errors.As(er, &invalidUnmarshall):
			{
				return fmt.Errorf("Invalid")
			}
		default:
			{
				return er
			}
		}
	}
	return nil

}

type Success struct {
	Msg string `json:"msg"`
}

type Error struct {
	Message string `json:"msg"`
	Code    string `json:"code"`
	Key     string `json:"key"`
}
