package api100

type apiResponse struct {
	Version string      `json:"version"`
	Path    string      `json:"path"`
	Data    interface{} `json:"data"`
}

type authResponse struct {
	Token string `json:"token"`
}
