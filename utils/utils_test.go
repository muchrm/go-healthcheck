package utils

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestGetDurationString(t *testing.T) {
	type args struct {
		t time.Duration
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "500 millisecond should return 500 ms",
			args: args{t: 500 * time.Millisecond},
			want: "500 ms",
		},
		{
			name: "50001 millisecond should return 50 sec 1 ms",
			args: args{t: 50001 * time.Millisecond},
			want: "50 sec 1 ms",
		},
		{
			name: "3*time.Minute + 12*time.Second + 7*time.Millisecond  should return 50 minutes 12 sec 7 ms",
			args: args{t: (3 * time.Minute) + (12 * time.Second) + (7 * time.Millisecond)},
			want: "3 minutes 12 sec 7 ms",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDurationString(tt.args.t); got != tt.want {
				t.Errorf("GetDurationString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetServerUrls(t *testing.T) {
	type args struct {
		reader string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "correct format urls",
			args: args{
				reader: `https://www.google.com
				https://www.facebook.com
				https://www.github.com
				http://localhost
				http://anybyhostfile
				http://localhost:8081`,
			},
			want: []string{
				"https://www.google.com",
				"https://www.facebook.com",
				"https://www.github.com",
				"http://localhost",
				"http://anybyhostfile",
				"http://localhost:8081",
			},
		},
		{
			name: "url incorrect url",
			args: args{
				reader: `incorrecturl`,
			},
			want:    []string{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetServerUrls(strings.NewReader(tt.args.reader))
			if (err != nil) != tt.wantErr {
				t.Errorf("GetServerUrls() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetServerUrls() = %v, want %v", got, tt.want)
			}
		})
	}
}
