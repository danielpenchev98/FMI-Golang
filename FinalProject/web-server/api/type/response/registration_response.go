package response

//RegistrationResponse is returned to the client when the registration was successfull
//it contains the statius of hist request and the jwt token
type RegistrationResponse struct {
	StatusCode int    `json:"status"`
	JWTToken   string `json:"jwt_token"`
}
