package todoist

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	httptools "github.com/XDoubleU/essentia/pkg/communication/http"
)

const BaseURLRESTAPI = "https://api.todoist.com/rest/v2"

type client struct {
	apiToken string
}

func New(apiToken string) Client {
	return client{
		apiToken: apiToken,
	}
}

func (client client) sendRequest(
	ctx context.Context,
	_ string,
	endpoint string,
	query string,
	dst any,
) error {
	u, err := url.Parse(fmt.Sprintf("%s/%s", BaseURLRESTAPI, endpoint))
	if err != nil {
		return err
	}

	u.RawQuery = query

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.apiToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	err = httptools.ReadJSON(res.Body, dst)
	if err != nil {
		return err
	}

	return nil
}
