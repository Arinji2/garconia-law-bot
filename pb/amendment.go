package pb

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/arinji2/law-bot/network"
)

func (p *PocketbaseAdmin) GetAllAmendments(expand bool) ([]AmendmentCollection, error) {
	parsedURL, err := url.Parse(p.BaseDomain)
	if err != nil {
		return nil, err
	}
	parsedURL.Path = "/api/collections/amendment/records"

	if expand {
		params := url.Values{}
		params.Add("expand", "clause")
		parsedURL.RawQuery = params.Encode()
	}

	type request struct{}
	responseBody, err := network.MakeAuthenticatedRequest(parsedURL, "GET", request{}, p.Token)
	if err != nil {
		return nil, err
	}

	var response PbResponse[AmendmentCollection]
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, err
	}

	return response.Items, nil
}

func (p *PocketbaseAdmin) GetAmendmentByNumber(amendmentNumber, clauseNumber, articleNumber string, expand bool) (AmendmentCollection, error) {
	parsedURL, err := url.Parse(p.BaseDomain)
	if err != nil {
		return AmendmentCollection{}, err
	}
	parsedURL.Path = "/api/collections/amendment/records"
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("number='%s' && clause.number='%s' && clause.article.number='%s'", amendmentNumber, clauseNumber, articleNumber))
	if expand {
		params.Add("expand", "clause,clause.article")
	}
	rawQuery := params.Encode()
	parsedURL.RawQuery = rawQuery

	type request struct{}
	responseBody, err := network.MakeAuthenticatedRequest(parsedURL, "GET", request{}, p.Token)
	if err != nil {
		return AmendmentCollection{}, err
	}

	var response PbResponse[AmendmentCollection]
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return AmendmentCollection{}, err
	}

	if len(response.Items) == 0 {
		return AmendmentCollection{}, fmt.Errorf("no amendments found for number: %s", clauseNumber)
	}

	return response.Items[0], nil
}

func (p *PocketbaseAdmin) GetAmendmentsByClause(clause string) ([]AmendmentCollection, error) {
	parsedURL, err := url.Parse(p.BaseDomain)
	if err != nil {
		return nil, err
	}
	parsedURL.Path = "/api/collections/amendment/records"

	params := url.Values{}
	params.Add("filter", fmt.Sprintf("clause.number='%s'", clause))
	params.Add("expand", "clause")
	rawQuery := params.Encode()
	parsedURL.RawQuery = rawQuery

	type request struct{}
	responseBody, err := network.MakeAuthenticatedRequest(parsedURL, "GET", request{}, p.Token)
	if err != nil {
		return nil, err
	}

	var response PbResponse[AmendmentCollection]
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, err
	}

	return response.Items, nil
}
