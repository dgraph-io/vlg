package model

import (
	"testing"
	"time"
)

func TestDateTime_UnmarshalCSV(t *testing.T) {
	type fields struct {
		Time time.Time
	}
	type args struct {
		csv string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid date",
			args: args{
				csv: "5/3/0004",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			date := &DateTime{}
			if err := date.UnmarshalCSV(tt.args.csv); (err != nil) != tt.wantErr {
				t.Errorf("DateTime.UnmarshalCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
