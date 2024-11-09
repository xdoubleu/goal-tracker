package temptools

import (
	"net/http"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func ReadForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = decoder.Decode(dst, r.PostForm)
	if err != nil {
		return err
	}

	return nil
}
