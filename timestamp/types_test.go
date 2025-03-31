package timestamp

import (
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestTimestamp_After(t *testing.T) {
	type fields struct {
		t time.Time
	}
	type args struct {
		ots *Timestamp
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &Timestamp{
				t: tt.fields.t,
			}
			if got := ts.After(tt.args.ots); got != tt.want {
				t.Errorf("After() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimestamp_Before(t *testing.T) {
	type fields struct {
		t time.Time
	}
	type args struct {
		ots *Timestamp
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &Timestamp{
				t: tt.fields.t,
			}
			if got := ts.Before(tt.args.ots); got != tt.want {
				t.Errorf("Before() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimestamp_MarshalBinary(t *testing.T) {
	type fields struct {
		t time.Time
	}
	tests := []struct {
		name     string
		fields   fields
		wantData []byte
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := Timestamp{
				t: tt.fields.t,
			}
			gotData, err := ts.MarshalBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("MarshalBinary() gotData = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func TestTimestamp_MarshalJSON(t *testing.T) {
	type fields struct {
		t time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "end of time",
			fields: fields{
				t: Enums().EndOfTime().t,
			},
			want:    []byte(strconv.Quote(Enums().EndOfTime().t.Format(time.RFC3339Nano))),
			wantErr: false,
		},
		{
			name: "beginning of time",
			fields: fields{
				t: Enums().BeginningOfTime().t,
			},
			want:    []byte(strconv.Quote(Enums().BeginningOfTime().t.Format(time.RFC3339Nano))),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := Timestamp{
				t: tt.fields.t,
			}
			got, err := ts.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimestamp_MarshalText(t *testing.T) {
	type fields struct {
		t time.Time
	}
	tests := []struct {
		name     string
		fields   fields
		wantText []byte
		wantErr  bool
	}{
		{
			name: "end of time",
			fields: fields{
				t: Enums().EndOfTime().t,
			},
			wantText: []byte(Enums().EndOfTime().t.Format(time.RFC3339Nano)),
			wantErr:  false,
		},
		{
			name: "beginning of time",
			fields: fields{
				t: Enums().BeginningOfTime().t,
			},
			wantText: []byte(Enums().BeginningOfTime().t.Format(time.RFC3339Nano)),
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := Timestamp{
				t: tt.fields.t,
			}
			gotText, err := ts.MarshalText()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotText, tt.wantText) {
				t.Errorf("MarshalText() gotText = %v, want %v", gotText, tt.wantText)
			}
		})
	}
}

func TestTimestamp_String(t *testing.T) {
	type fields struct {
		t time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := Timestamp{
				t: tt.fields.t,
			}
			if got := ts.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimestamp_UnmarshalBinary(t *testing.T) {
	type fields struct {
		t time.Time
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &Timestamp{
				t: tt.fields.t,
			}
			if err := ts.UnmarshalBinary(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTimestamp_UnmarshalText(t *testing.T) {
	type fields struct {
		t time.Time
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &Timestamp{
				t: tt.fields.t,
			}
			if err := ts.UnmarshalText(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTimestamp_UnmarshallJSON(t *testing.T) {
	type fields struct {
		t time.Time
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &Timestamp{
				t: tt.fields.t,
			}
			if err := ts.UnmarshallJSON(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshallJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_tsenums_BeginningOfTime(t *testing.T) {
	tests := []struct {
		name string
		want Timestamp
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := enums{}
			if got := i.BeginningOfTime(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BeginningOfTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tsenums_EndOfTime(t *testing.T) {
	tests := []struct {
		name string
		want Timestamp
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := enums{}
			if got := i.EndOfTime(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EndOfTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
