package url

import (
	neturl "net/url"
)

func MustParse(rawURL string) *neturl.URL {
	if u, err := neturl.Parse(rawURL); err != nil {
		panic(err)
	} else {
		return u
	}

}

func MustJoin(url *neturl.URL, elems ...string) string {
	if s, err := neturl.JoinPath(url.String(), elems...); err != nil {
		panic(err)
	} else {
		return s
	}
}
