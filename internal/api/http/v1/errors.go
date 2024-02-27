package v1

//var (
//	ErrorValidate = errors.New("some fields are incorrect")
//)
//TODO: Create spatial errors understandable for users

type ErrorResponse struct {
	Error string `json:"error"`
}
