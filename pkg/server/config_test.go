package application

import (
	"reflect"
	"testing"
)

func Test_overrideConfig(t *testing.T) {
	type args struct {
		init     map[string]string
		override map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{"test1", args{map[string]string{"key1": "1", "key2": "2"}, map[string]string{"key1": "2"}}, map[string]string{"key1": "2", "key2": "2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := overrideConfig(tt.args.init, tt.args.override); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("overrideConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkRequiredConfigsPresent(t *testing.T) {
	type args struct {
		config map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkRequiredConfigsPresent(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("checkRequiredConfigsPresent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
