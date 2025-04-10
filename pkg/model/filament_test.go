package model

import "testing"

func TestExploreType(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name      string
		args      args
		wantMajor string
		wantMinor string
	}{
		{
			name: "PLA",
			args: args{
				name: "PLA",
			},
			wantMajor: "PLA",
			wantMinor: FilamentTypeDefaultMinor,
		},
		{
			name: "PLA-BASIC",
			args: args{
				name: "PLA-BASIC",
			},
			wantMajor: "PLA",
			wantMinor: "BASIC",
		},
		{
			name: "PETG-哑光",
			args: args{
				name: "PETG-哑光",
			},
			wantMajor: "PETG",
			wantMinor: "哑光",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMajor, gotMinor := ExploreType(tt.args.name)
			if gotMajor != tt.wantMajor {
				t.Errorf("ExploreType() gotMajor = %v, want %v", gotMajor, tt.wantMajor)
			}
			if gotMinor != tt.wantMinor {
				t.Errorf("ExploreType() gotMinor = %v, want %v", gotMinor, tt.wantMinor)
			}
		})
	}
}
