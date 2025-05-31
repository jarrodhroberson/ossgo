package firestore

import (
	"reflect"
	"testing"

	"cloud.google.com/go/firestore"

	"github.com/jarrodhroberson/ossgo/functions/must"
)

func TestNewProjection(t *testing.T) {
	type args struct {
		paths []string
	}
	tests := []struct {
		name string
		args args
		want Projection
	}{
		{
			name: "Id&LastUpdatedId",
			args: args{paths: []string{"snippet.id", "last_updated_id"}},
			want: Projection{
				fieldPaths: []firestore.FieldPath{[]string{"snippet", "id"}, []string{"last_updated_id"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProjection(tt.args.paths...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProjection() = %s, want %s", must.MarshalJson(got.fieldPaths), must.MarshalJson(tt.want.fieldPaths))
			}
		})
	}
}
