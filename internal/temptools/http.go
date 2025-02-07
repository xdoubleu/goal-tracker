package temptools

import (
	"fmt"
	"net/http"
)

func RedirectWithError(w http.ResponseWriter, r *http.Request, url string, err error) {
	http.Redirect(
		w,
		r,
		fmt.Sprintf("%s?error=%s", url, err.Error()),
		http.StatusSeeOther,
	)
}
