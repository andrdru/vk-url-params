package vkurlparams

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

var (
	ErrUrlParamsInvalid    = errors.New("invalid url params")
	errGetLengthUnexpected = errors.New("unexpected GET param length")
)

// Parse parse and validate vk params
// r *http.Request; v := r.URL.Query()
func Parse(v url.Values, secret string, isDebug bool) (*urlParams, error) {
	var sign, vkParams, err = ReadParams(v, isDebug)
	if err != nil {
		return nil, fmt.Errorf("params: %w", err)
	}

	if !isDebug {
		err = Validate(vkParams, sign, secret)
		if err != nil {
			return nil, fmt.Errorf("validate: %w", err)
		}
	}

	return NewURLParams(vkParams), nil
}

// Validate validate vk params
func Validate(params map[string]string, sign string, secret string) error {
	if len(params) == 0 {
		return ErrUrlParamsInvalid
	}

	var keys = make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	var u = url.Values{}
	for _, key := range keys {
		u.Add(key, params[key])
	}

	var queryStr = u.Encode()

	h := hmac.New(sha256.New, []byte(secret))
	_, _ = h.Write([]byte(queryStr)) //nolint:errcheck

	var b64 = base64.StdEncoding.EncodeToString(h.Sum(nil))

	var result = strings.ReplaceAll(b64, "+", "-")
	result = strings.ReplaceAll(result, "/", "_")
	result = strings.TrimRight(result, "=")

	if sign != result {
		return ErrUrlParamsInvalid
	}

	return nil
}

// ReadParams read vk_url_params from url.Values
// isDebug should be passed to parse data with invalid vk_url_params (test cases etc)
func ReadParams(v url.Values, isDebug bool) (sign string, vkParams map[string]string, err error) {
	var signParam, ok = v["sign"]
	if !ok && !isDebug {
		return "", nil, ErrUrlParamsInvalid
	}

	if len(signParam) == 0 && !isDebug {
		return "", nil, errGetLengthUnexpected
	}

	if len(signParam) > 0 {
		sign = signParam[0]
	}

	vkParams = make(map[string]string)

	var key string
	var value []string
	for key, value = range v {
		if !strings.HasPrefix(key, "vk_") {
			continue
		}

		if len(value) == 0 && !isDebug {
			return "", nil, errGetLengthUnexpected
		}

		vkParams[key] = value[0]
	}

	return sign, vkParams, nil
}
