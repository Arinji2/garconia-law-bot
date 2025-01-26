package pb

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/arinji2/law-bot/network"
)

func (p *PocketbaseAdmin) GetAllClauses(expand bool) ([]ClauseCollection, error) {
	parsedURL, err := url.Parse(p.BaseDomain)
	if err != nil {
		return nil, err
	}
	parsedURL.Path = "/api/collections/clause/records"

	if expand {
		params := url.Values{}
		params.Add("expand", "article")
		parsedURL.RawQuery = params.Encode()
	}

	type request struct{}
	responseBody, err := network.MakeAuthenticatedRequest(parsedURL, "GET", request{}, p.Token)
	if err != nil {
		return nil, err
	}

	var response PbResponse[ClauseCollection]
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, err
	}

	return response.Items, nil
}

func (p *PocketbaseAdmin) GetClauseByNumber(clauseNumber, articleNumber string, expand bool) (ClauseCollection, error) {
	parsedURL, err := url.Parse(p.BaseDomain)
	if err != nil {
		return ClauseCollection{}, err
	}
	parsedURL.Path = "/api/collections/clause/records"
	params := url.Values{}
	params.Add("filter", fmt.Sprintf("number='%s' && article.number='%s'", clauseNumber, articleNumber))
	if expand {
		params.Add("expand", "article")
	}
	rawQuery := params.Encode()
	parsedURL.RawQuery = rawQuery

	type request struct{}
	responseBody, err := network.MakeAuthenticatedRequest(parsedURL, "GET", request{}, p.Token)
	if err != nil {
		return ClauseCollection{}, err
	}

	var response PbResponse[ClauseCollection]
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return ClauseCollection{}, err
	}

	if len(response.Items) == 0 {
		return ClauseCollection{}, fmt.Errorf("no clauses found for number: %s", clauseNumber)
	}

	return response.Items[0], nil
}

func (p *PocketbaseAdmin) GetClausesByArticle(article string) ([]ClauseCollection, error) {
	parsedURL, err := url.Parse(p.BaseDomain)
	if err != nil {
		return nil, err
	}
	parsedURL.Path = "/api/collections/clause/records"

	params := url.Values{}
	params.Add("filter", fmt.Sprintf("article.number='%s'", article))
	params.Add("expand", "article")
	rawQuery := params.Encode()
	parsedURL.RawQuery = rawQuery

	type request struct{}
	responseBody, err := network.MakeAuthenticatedRequest(parsedURL, "GET", request{}, p.Token)
	if err != nil {
		return nil, err
	}

	var response PbResponse[ClauseCollection]
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, err
	}

	return response.Items, nil
}
