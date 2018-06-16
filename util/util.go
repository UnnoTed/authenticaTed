package util

import (
	"encoding/json"
	"net/http"

	"github.com/go-resty/resty"
	"github.com/labstack/echo"
)

type J map[string]interface{}

func ParseToken(token string) string {
	l := len("Bearer")
	if len(token) > l+1 && token[:l] == "Bearer" {
		return token[l+1:]
	}

	return ""
}

func PegarIDEcho(c echo.Context, key string) (string, error) {
	return PegarID(c.Get("claims").(map[string]interface{}), key)
}

func PegarID(jwt map[string]interface{}, key string) (string, error) {
	if jwt["id"] != nil {
		id, err := Decrypt(jwt["id"].(string), key)
		if err != nil {
			return "", err
		}

		return id, nil
	}

	return "", nil
}

func Error(msg string) J {
	return J{"sucesso": false, "erro": J{"mensagem": msg}}
}

func ErroJSON(c echo.Context, msg string) error {
	return c.JSON(http.StatusInternalServerError, Error(msg))
}

func ErrJSON(c echo.Context, err error) error {
	return c.JSON(http.StatusInternalServerError, Error(err.Error()))
}

func JSONaStruct(resp *resty.Response, s, e interface{}) error {
	var err error

	ok := resp.StatusCode() == http.StatusOK

	if ok {
		err = json.Unmarshal(resp.Body(), &s)
	} else {
		err = json.Unmarshal(resp.Body(), &e)
	}

	return err
}
