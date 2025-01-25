package pb

type PocketbaseAdmin struct {
	Token  string `json:"token"`
	Record struct {
		ID string `json:"id"`
	} `json:"record"`
}
