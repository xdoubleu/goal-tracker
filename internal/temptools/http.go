package temptools

import (
	"fmt"
	"net/http"

	"github.com/gorilla/schema"
)

//nolint:gochecknoglobals //ok
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

func RedirectWithError(w http.ResponseWriter, r *http.Request, url string, err error) {
	http.Redirect(
		w,
		r,
		fmt.Sprintf("%s?error=%s", url, err.Error()),
		http.StatusSeeOther,
	)
}
