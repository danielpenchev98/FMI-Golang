package response

//SuccessResponse is send to the client when his request was successfully fulfilled
type SuccessResponse struct {
	Status int `json:"status"`
}
