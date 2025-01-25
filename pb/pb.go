package pb

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/arinji2/law-bot/env"
	"github.com/arinji2/law-bot/network"
)

func SetupPocketbase(pb env.PB) *PocketbaseAdmin {
	parsedURL, err := url.Parse(pb.BaseDomain)
	if err != nil {
		log.Fatal(err)
	}
	parsedURL.Path = "/api/collections/_superusers/auth-with-password"
	type request struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}

	body := request{
		Identity: pb.Email,
		Password: pb.Password,
	}

	responseBody, err := network.MakeRequest(parsedURL, "POST", body)
	if err != nil {
		log.Fatal(err)
	}

	var response PocketbaseAdmin
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		log.Fatal(err)
	}
	response.BaseDomain = pb.BaseDomain
	return &response
}
