package main

import "testing"

func Test_queryLibPostal(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "address with name",
			args: args{
				address: "Kalkas Business Services 245 SE 1st Street Ste,225 Miami Florida",
			},
			want: "245 se 1st street, miami, florida",
		},
		{
			name: "address with c/o in name",
			args: args{
				address: "C/O Kalkas Business Services 245 SE 1st Street Ste,225 Miami Florida",
			},
			want: "245 se 1st street, miami, florida",
		},
		{
			name: "misspelled work suite",
			args: args{
				address: "c/o Georgian Pine Invesments, 2200 Sand Hill Road, suit 240 MENLO PARK CA 94025",
			},
			want: "240 sand hill road suit, menlo park, ca",
		},
		// C/O Law Offices of Ernesto Sanchez 815 Ponce de Leon Blvd., Suite 306 Coral Gables, FL   33134 USA
		{
			name: "misspelled work suite",
			args: args{
				address: "C/O Law Offices of Ernesto Sanchez 815 Ponce de Leon Blvd., Suite 306 Coral Gables, FL   33134 USA",
			},
			want: "815 ponce de leon blvd., coral gables, fl",
		},
		// 1145 South Camino Del Rio # 145 Durango, Colorado 81303
		{
			name: "correct",
			args: args{
				address: "1145 South Camino Del Rio # 145 Durango, Colorado 81303",
			},
			want: "1145 south camino del rio, durango, colorado",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := queryLibPostal(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("queryLibPostal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("queryLibPostal() = %v, want %v", got, tt.want)
			}
		})
	}
}
