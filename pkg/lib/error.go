package lib

// @Description Http Error structure
type HttpError struct {
	Code    int    `json:"code"`    // Http status code
	Name    string `json:"name"`    // Http status name
	Message string `json:"message"` // Error message
}
