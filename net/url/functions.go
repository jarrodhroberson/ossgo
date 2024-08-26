package url

import (
	neturl "net/url"
)

func MustJoin(url neturl.URL, elems ...string) string {
	if s, err := neturl.JoinPath(url.String(), elems...); err != nil {
		panic(err)
	} else {
		return s
	}
}
