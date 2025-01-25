package pb

import (
	"encoding/json"
	"log"

	"github.com/arinji2/law-bot/env"
	"github.com/arinji2/law-bot/network"
)

var baseURL string

func SetupPocketbase(pb env.PB) *PocketbaseAdmin {
	baseURL = pb.BaseDomain
	url := baseURL + "/api/collections/_superusers/auth-with-password"
	type request struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}

	body := request{
		Identity: pb.Email,
		Password: pb.Password,
	}

	responseBody, err := network.MakeRequest(url, "POST", body)
	if err != nil {
		log.Fatal(err)
	}

	var response PocketbaseAdmin
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		log.Fatal(err)
	}

	return &response
}
