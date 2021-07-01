package vkurlparams

import (
	"context"
	"strconv"
)

// nolint: golint
type (
	ctxKey string

	URLParams interface {
		VKUserID() int64
		IsAppUser() bool
	}

	urlParams struct {
		userID                  int64
		ref                     string
		platform                string
		language                string
		isFavorite              bool
		isAppUser               bool
		areNotificationsEnabled bool
		vkAppID                 int64
		accessTokenSettings     string
	}
)

var ctxKeyVar ctxKey = "urlparams"
var _ URLParams = &urlParams{}

// NewURLParams init
// nolint: golint
func NewURLParams(m map[string]string) *urlParams {
	return &urlParams{
		userID:                  shouldParseInt64(strconv.ParseInt(m["vk_user_id"], 10, 64)),
		ref:                     m["vk_ref"],
		platform:                m["vk_platform"],
		language:                m["vk_language"],
		isFavorite:              boolFromString(m["vk_is_favorite"]),
		isAppUser:               boolFromString(m["vk_is_app_user"]),
		areNotificationsEnabled: boolFromString(m["vk_are_notifications_enabled"]),
		vkAppID:                 shouldParseInt64(strconv.ParseInt(m["vk_app_id"], 10, 64)),
		accessTokenSettings:     m["vk_access_token_settings"],
	}
}

func (u *urlParams) VKUserID() int64 {
	return u.userID
}

func (u *urlParams) IsAppUser() bool {
	return u.isAppUser
}

// CtxGet get urlparams from context
// nolint: golint
func CtxGet(ctx context.Context) *urlParams {
	var i = ctx.Value(ctxKeyVar)
	if i == nil {
		return nil
	}

	return i.(*urlParams)
}

// CtxSet put urlparams to context
func CtxSet(ctx context.Context, params *urlParams) context.Context {
	return context.WithValue(ctx, ctxKeyVar, params)
}

func boolFromString(s string) bool {
	return s == "1"
}

func shouldParseInt64(i int64, err error) int64 {
	if err != nil {
		return 0
	}
	return i
}
