package todoist

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	httptools "github.com/XDoubleU/essentia/pkg/communication/http"
)

var BASE_URL_REST_API = "https://api.todoist.com/rest/v2"

type Client struct {
	apiToken string
}

func NewClient(apiToken string) Client {
	return Client{
		apiToken: apiToken,
	}
}

func (client Client) sendRequest(ctx context.Context, method string, endpoint string, query string, dst any) error {
	u, err := url.Parse(fmt.Sprintf("%s/%s", BASE_URL_REST_API, endpoint))
	if err != nil {
		return err
	}

	u.RawQuery = query

	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.apiToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	err = httptools.ReadJSON(res.Body, dst)
	if err != nil {
		return err
	}

	return nil
}
