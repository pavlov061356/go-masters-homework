package config

import (
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	type args struct {
		path string
	}

	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "valid config file",
			args: args{
				path: "test_data/valid_config.yaml",
			},
			want: &Config{
				Port:   8080,
				DBPath: "test.db",
			},
			wantErr: false,
		},
		{
			name: "invalid config file, no db_path",
			args: args{
				path: "test_data/invalid_config_no_db_path.yaml",
			},
			wantErr: true,
		},
		{
			name: "invalid config file, no port",
			args: args{
				path: "test_data/invalid_config_invalid_port.yaml",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.args.path)

			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
}
