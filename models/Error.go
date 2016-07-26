package models

import (
	"encoding/json"
	"fmt"
)

//RequestError Implement WSResponse, it's an Object that represent an error to send back to the client
type RequestError struct {
	//Title of the error
	Title string `json:"title"`
	//Description of the error
	Description string `json:"description"`
	//Code of the error
	Code int `json:"code"`
}

//Error return well format error string
func (e RequestError) Error() string {
	return fmt.Sprint("Error", e.Code, e.Title, e.Description)
}

func (e RequestError) MarshalJSON() ([]byte, error) {
	type Alias RequestError
	return json.Marshal(&struct {
		Error Alias
	}{
		Error: (Alias)(e),
	})
}
