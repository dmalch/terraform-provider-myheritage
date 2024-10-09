package myheritage

type GraphqlRequest struct {
	Query       string            `json:"query"`
	Variables   map[string]string `json:"variables,omitempty"`
	Description string            `json:"description"`
}
