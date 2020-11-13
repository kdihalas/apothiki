package http

type ErrorResponse struct {
	Errors []ErrorMessage `json:"errors"`
}

type ErrorMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json"detail"`
}

type RepositoriesResponse struct {
	Repositories []string `json:"repositories"`
}

type RepositoriesTagResponse struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}
