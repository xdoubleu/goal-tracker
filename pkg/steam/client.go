package steam

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	httptools "github.com/XDoubleU/essentia/pkg/communication/http"
)

const BaseURLRESTAPI = "http://api.steampowered.com"

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
	endpoint string,
	query string,
	dst any,
) error {
	u, err := url.Parse(fmt.Sprintf("%s/%s", BaseURLRESTAPI, endpoint))
	if err != nil {
		return err
	}

	u.RawQuery = query
	tempQuery := u.Query()
	tempQuery.Add("key", client.apiToken)
	u.RawQuery = tempQuery.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

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
