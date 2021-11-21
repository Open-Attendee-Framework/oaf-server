package api100

type errorResponse struct {
	Httpstatus   string       `json:"httpstatus"`
	Errorcode    apiErrorcode `json:"errorcode"`
	Errormessage string       `json:"errormessage"`
}

type apiErrorcode int

const (
	errorWrongCredentials apiErrorcode = 1 + iota
	errorDBQueryFailed
	errorMalformedAuth
	errorNoHash
	errorNoToken
	errorJSONError
	errorUserNotAuthorized
	errorInvalidParameter
	errorNotFound
)

func (e *apiErrorcode) String() string {
	switch *e {
	case errorWrongCredentials:
		return "Wrong Username or Password"
	case errorDBQueryFailed:
		return "Database Query failed"
	case errorMalformedAuth:
		return "Authorization request is malformed"
	case errorNoHash:
		return "Could not generate Hash from Password"
	case errorNoToken:
		return "Could not generate Token"
	case errorJSONError:
		return "JSON Marshal error"
	case errorUserNotAuthorized:
		return "User not authorized"
	case errorInvalidParameter:
		return "Invalid parameter"
	case errorNotFound:
		return "Resource not found"
	default:
		return "unknown error"
	}
}
