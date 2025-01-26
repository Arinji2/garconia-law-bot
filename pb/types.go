package pb

type PocketbaseAdmin struct {
	Token  string `json:"token"`
	Record struct {
		ID string `json:"id"`
	} `json:"record"`
	BaseDomain string
}

type PbResponse[T any] struct {
	Items      []T `json:"items"`
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}
type BaseCollection struct {
	ID          string `json:"id"`
	Number      string `json:"number"`
	Description string `json:"description"`
	CreatedAt   string `json:"created"`
	UpdatedAt   string `json:"updated"`
}

type ClauseCollection struct {
	Article string `json:"article"`
	BaseCollection
	Expand struct {
		Article BaseCollection
	} `json:"expand"`
}

type AmendmentCollection struct {
	Clause string `json:"clause"`
	BaseCollection
	Expand struct {
		Clause BaseCollection
	} `json:"expand"`
}
