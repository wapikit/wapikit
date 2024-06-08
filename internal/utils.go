package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetRequestBodyJson[T interface{}](context CustomContext) (map[string]interface{}, error) {
	body, err := io.ReadAll(context.Request().Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var payload interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
	}
	return payload.(map[string]interface{}), nil
}

func FetchWhatsappBusinessAccountDetails() {

}
