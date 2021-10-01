package core

import (
	"fmt"
	"net/http"

	"github.com/JunqiZhang0/tfacon/common"
	"github.com/JunqiZhang0/tfacon/connectors"
	"github.com/spf13/viper"
)

// TFACon is the general interface for all TFA Classifer
// connectors, any connector to any test management platform
// should inpement this interface.
type TFACon interface {
	GetAllTestIds() []string
	BuildUpdatedList(ids []string, concurrent bool, add_attributes bool) common.GeneralUpdatedList
	UpdateAll(common.GeneralUpdatedList, bool)
	String() string
	InitConnector()
	Validate(verbose bool) (bool, error)
}

// Run method is the run operation for any type of connector that
// implements TFACon interface.
func Run(viperRun, viperConfig *viper.Viper) {
	var con TFACon = GetCon(viperRun)
	con.InitConnector()
	runHelper(viperConfig, con.GetAllTestIds(), con)
}

func runHelper(viperConfig *viper.Viper, ids []string, con TFACon) {
	if len(ids) == 0 {
		return
	}
	updated_list_of_issues := con.BuildUpdatedList(ids, viperConfig.GetBool("config.concurrency"), viperConfig.GetBool("config.add_attributes"))
	// Doing this because the api can only take 20 items per request
	con.UpdateAll(updated_list_of_issues, viperConfig.GetBool("config.verbose"))
	runHelper(viperConfig, con.GetAllTestIds(), con)
}

// GetInfo method is the get info operation for any type of connector that
// implements TFACon interface.
func GetInfo(viper *viper.Viper) TFACon {
	con := GetCon(viper)
	return con
}

// Validate method is the validate operation for any type of connector that
// implements TFACon interface.
func Validate(con TFACon, viper *viper.Viper) (bool, error) {
	// var con TFACon = GetCon(viper)
	success, err := con.Validate(viper.GetBool("config.verbose"))
	return success, err
}

// GetCon method is the get con operation for any type of connector that
// implements TFACon interface, it returns the actual tfa connector instance.
func GetCon(viper *viper.Viper) TFACon {
	var con TFACon
	switch viper.Get("CONNECTOR_TYPE") {
	case "RPCon":
		con = &connectors.RPConnector{Client: &http.Client{}}
		err := viper.Unmarshal(con)
		if err != nil {
			fmt.Println(err)
		}
	// case "POLCon":
	// 	con = RPConnector{}
	// case "JiraCon":
	// 	con = RPConnector{}
	default:
		con = &connectors.RPConnector{Client: &http.Client{}}
		err := viper.Unmarshal(con)
		if err != nil {
			fmt.Println(err)
		}
	}
	return con
}
