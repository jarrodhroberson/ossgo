package identifiers

import (
	"testing"
)

const valid_lei = "W22L00WP2IHZNBB6K585"

func TestLeiCodeFrom(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		args    args
		want    LeiCode
		wantErr bool
	}{
		{
			name: "from string representation",
			args: args{
				code: valid_lei,
			},
			want:    LeiCode(valid_lei),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LeiCodeFrom(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("LeiCodeFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LeiCodeFrom() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLeiCode_EntityIdentifier(t *testing.T) {
	tests := []struct {
		name string
		l    LeiCode
		want string
	}{
		{
			name: "entity identifier",
			l:    LeiCode(valid_lei),
			want: valid_lei[6:18],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.EntityIdentifier(); got != tt.want {
				t.Errorf("EntityIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLeiCode_LocalOperatingUnit(t *testing.T) {
	tests := []struct {
		name string
		l    LeiCode
		want string
	}{
		{
			name: "local operating unit (LOU)",
			l:    LeiCode(valid_lei),
			want: valid_lei[0:4],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.LocalOperatingUnit(); got != tt.want {
				t.Errorf("LocalOperatingUnit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLeiCode_String(t *testing.T) {
	tests := []struct {
		name string
		l    LeiCode
		want string
	}{
		{
			name: "LieCode.String()",
			l:    LeiCode(valid_lei),
			want: valid_lei,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewLeiCode(t *testing.T) {
	type args struct {
		lou        string
		identifier string
	}
	tests := []struct {
		name    string
		args    args
		want    LeiCode
		wantErr bool
	}{
		{
			name: "valid lei",
			args: args{
				lou:        valid_lei[0:4],
				identifier: valid_lei[6:18],
			},
			want:    LeiCode(valid_lei),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewLeiCode(tt.args.lou, tt.args.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLeiCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewLeiCode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateChecksum(t *testing.T) {
	type args struct {
		lei string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid checksum",
			args: args{
				lei: valid_lei[:18],
			},
			want: valid_lei[18:],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateChecksum(tt.args.lei); got != tt.want {
				t.Errorf("calculateChecksum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isChecksumValid(t *testing.T) {
	type args struct {
		leiCode LeiCode
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid checksum",
			args: args{
				leiCode: LeiCode(valid_lei),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := isChecksumValid(tt.args.leiCode); (err != nil) != tt.wantErr {
				t.Errorf("isChecksumValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_isFormatValid(t *testing.T) {
	type args struct {
		leiCode LeiCode
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid format",
			args: args{
				leiCode: LeiCode(valid_lei),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := isFormatValid(tt.args.leiCode); (err != nil) != tt.wantErr {
				t.Errorf("isFormatValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_isUppercaseAlphaNumeric(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid uppercase alphanumeric",
			args: args{
				s: valid_lei[:18],
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := isUppercaseAlphaNumeric(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("isUppercaseAlphaNumeric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_isValidLeiCode(t *testing.T) {
	type args struct {
		leiCode LeiCode
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid lei code",
			args: args{
				leiCode: LeiCode(valid_lei),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := isValidLeiCode(tt.args.leiCode); (err != nil) != tt.wantErr {
				t.Errorf("isValidLeiCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
