package pb

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/arinji2/law-bot/network"
)

func (p *PocketbaseAdmin) GetAllArticles() ([]BaseCollection, error) {
	parsedURL, err := url.Parse(p.BaseDomain)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base domain: %w", err)
	}
	parsedURL.Path = "/api/collections/article/records"

	type request struct{}
	responseBody, err := network.MakeAuthenticatedRequest(parsedURL, "GET", request{}, p.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to make authenticated request: %w", err)
	}

	var response PbResponse[BaseCollection]
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Items, nil
}