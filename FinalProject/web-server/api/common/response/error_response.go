package response

/*ErrorResponse is sent to the client of the REST API when
there is an error with request or the server*/
type ErrorResponse struct {
	ErrorCode int    `json:"errorcode"` //status code of the request - 4xx or 5xx
	ErrorMsg  string `json:"message"`   //desription of the error
}
