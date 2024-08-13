package secrets

import (
	"testing"
)

func TestPath_String(t *testing.T) {
	type fields struct {
		ProjectNumber int
		Name          string
		Version       int
	}
	var tests = []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "With Version",
			fields: fields{
				ProjectNumber: 0,
				Name:          "WithVersion",
				Version:       1,
			},
			want: "projects/0/secrets/WithVersion/versions/1",
		},
		{
			name: "Without Version",
			fields: fields{
				ProjectNumber: 0,
				Name:          "WithoutVersion",
				Version:       -1,
			},
			want: "projects/0/secrets/WithoutVersion",
		},
		{
			name: "Latest Version",
			fields: fields{
				ProjectNumber: 0,
				Name:          "WithLatest",
				Version:       0,
			},
			want: "projects/0/secrets/WithLatest/versions/latest",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Path{
				ProjectNumber: tt.fields.ProjectNumber,
				Name:          tt.fields.Name,
				Version:       tt.fields.Version,
			}
			if got := p.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
