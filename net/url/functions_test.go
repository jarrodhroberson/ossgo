package url

import (
	neturl "net/url"
	"reflect"
	"testing"
)

func TestMustJoin(t *testing.T) {
	type args struct {
		url   *neturl.URL
		elems []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Append Path to base URL",
			args: args{
				url:   MustParse("https://stripe.com"),
				elems: []string{"cust_2340987872303"},
			},
			want: "https://stripe.com/cust_2340987872303",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MustJoin(tt.args.url, tt.args.elems...); got != tt.want {
				t.Errorf("MustJoin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustParse(t *testing.T) {
	type args struct {
		rawURL string
	}
	tests := []struct {
		name string
		args args
		want *neturl.URL
	}{
		{
			name: "Parse Base URL",
			args: args{
				rawURL: "https://stripe.com",
			},
			want: func() *neturl.URL {
				if u, err := neturl.Parse("https://stripe.com"); err != nil {
					panic(err)
				} else {
					return u
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MustParse(tt.args.rawURL); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MustParse() = %v, want %v", got, tt.want)
			}
		})
	}
}
