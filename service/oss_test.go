package service

import "testing"

func Test_strToSecond(t *testing.T) {
	tests := []struct {
		name string
		args string
		want int
	}{
		{
			name: "1",
			args: "0:01.89",
			want: 2,
		},
		{
			name: "2",
			args: "0:01.512",
			want: 2,
		},
		{
			name: "3",
			args: "0:31.89",
			want: 32,
		},
		{
			name: "4",
			args: "10:31.89",
			want: 632,
		},
		{
			name: "5",
			args: "1:10:31.89",
			want: 4232,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strToSecond(tt.args); got != tt.want {
				t.Errorf("strToSecond() = %v, want %v", got, tt.want)
			}
		})
	}
}
