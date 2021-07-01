package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	vkurlparams "github.com/andrdru/vk-url-params"
)

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
	// see http://localhost:8080/?<vk_url_params>
}

func handler(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	// set debug = true, see http://localhost:8080/?vk_user_id=123
	urlParams, err := vkurlparams.Parse(v, "my_vk_secret_key", false)
	if err != nil {
		if errors.Is(err, vkurlparams.ErrUrlParamsInvalid) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("forbidden"))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("invalid params passed"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("vk user id = %d", urlParams.VKUserID())))
}
