package common

import (
	"strings"
)

const (
	NEW_CDN_DOMAIN = "https://rescdn.dokidokilive.com/upload_image/"
)

func GetImageUrl(uri string) string {
	if uri == "" {
		return uri
	}
	if !strings.HasPrefix(uri, "http") {
		uri = NEW_CDN_DOMAIN + uri
	}
	return uri
}
