package main

import "testing"

func Test_fibonacci(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "-1",
			args: args{
				n: -1,
			},
			want: 0,
		},
		{
			name: "0",
			args: args{
				n: 0,
			},
			want: 0,
		},
		{
			name: "1",
			args: args{
				n: 1,
			},
			want: 1,
		},
		{
			name: "2",
			args: args{
				n: 2,
			},
			want: 1,
		},
		{
			name: "5",
			args: args{
				n: 5,
			},
			want: 5,
		},
		{
			name: "10",
			args: args{
				n: 10,
			},
			want: 55,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fibonacci(tt.args.n); got != tt.want {
				t.Errorf("fibonacci() = %v, want %v", got, tt.want)
			}
		})
	}
}
