package core

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/JunqiZhang0/tfacon/common"
	"github.com/JunqiZhang0/tfacon/connectors"
	"github.com/spf13/viper"
)

// var viperconfigfortest *viper.Viper = viper.New()

// func Test_runHelper(t *testing.T) {
// 	viperconfigfortest.SetConfigType("ini")
// 	viperconfigfortest.SetDefault("config.concurrency", true)
// 	viperconfigfortest.SetDefault("config.retry_times", 1)
// 	viperconfigfortest.SetDefault("config.add_attributes", false)
// 	type args struct {
// 		viperConfig *viper.Viper
// 		ids         []string
// 		con         TFACon
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "test with single id in ids",
// 			args: args{
// 				viperConfig: viperconfigfortest,
// 				ids:         []string{"12306"},
// 				con:         GetCon(vipertfaconfortest),
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			defer func() {
// 				if r := recover(); r != nil {
// 					t.Fatalf("%s failed on: %s\n", tt.name, r)
// 				}
// 			}()
// 			runHelper(tt.args.viperConfig, tt.args.ids, tt.args.con)

// 		})
// 	}
// }

var vipertfaconfortest *viper.Viper = viper.New()

func TestGetCon(t *testing.T) {
	// vipertfaconfortest.SetConfigFile("yaml")
	vipertfaconfortest.SetConfigName("tfacon")
	vipertfaconfortest.AddConfigPath("../test_data/workspace_data")
	err := vipertfaconfortest.ReadInConfig()
	common.HandleError(err)
	rpcon := &connectors.RPConnector{Client: &http.Client{}}
	err = vipertfaconfortest.Unmarshal(rpcon)
	common.HandleError(err)
	type args struct {
		viper *viper.Viper
	}
	tests := []struct {
		name string
		args args
		want TFACon
	}{
		// TODO: Add test cases.
		{
			name: "test tfacon get with rpcon",
			args: args{
				viper: vipertfaconfortest,
			},
			want: rpcon,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCon(tt.args.viper); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCon() = %v, want %v", got, tt.want)
			} else {
				t.Logf("GetCon() = %v, want %v", got, tt.want)
			}
		})
	}
}
