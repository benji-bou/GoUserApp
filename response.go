package app

import (
	"fmt"
	"net/http"
	"tools"
)

//ResponseType explain the type response
type ResponseType uint

const (

	//JSON form for response
	JSON ResponseType = 1 << iota
	//XML form for response
	XML
)

type jsonResponse interface {
	JSONResp() (string, error)
}

type xmlResponse interface {
	XMLResp() (string, error)
}

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

//JSONResp write the interface into a Json Format  it uses json.Marshal if interface get a type name, the json as a parent object with this name
func JSONResp(e http.ResponseWriter, rep interface{}) error {
	e.Header().Set("Content-Type", "application/json")
	json, err := tools.JSONResp(rep)
	if err == nil {
		fmt.Fprint(e, json)
		return nil
	}
	return err
}

//XMLResp write the interface into a XML Format the XML as a parent object with this name
func XMLResp(e http.ResponseWriter, rep interface{}) error {
	e.Header().Set("Content-Type", "application/xml")
	xml, err := tools.XMLResp(rep)
	if err == nil {
		fmt.Fprint(e, xml)
		return nil
	}
	return err
}
