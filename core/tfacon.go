package core

import (
	"fmt"
	"net/http"

	"github.com/JunqiZhang0/tfacon/connectors"

	"github.com/JunqiZhang0/tfacon/common"

	"github.com/spf13/viper"
)

type TFACon interface {
	GetAllTestIds() []string
	BuildUpdatedList(ids []string, concurrent bool, add_attributes bool) common.GeneralUpdatedList
	UpdateAll(common.GeneralUpdatedList, bool)
	String() string
	InitConnector()
	Validate(verbose bool) (bool, error)
}

func Run(viperRun, viperConfig *viper.Viper) {
	var con TFACon = GetCon(viperRun)
	// fmt.Printf("%+v\n", con)
	// fmt.Println("===========================")
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

func GetInfo(viper *viper.Viper) TFACon {
	var con TFACon = GetCon(viper)
	return con
}

func Validate(viper *viper.Viper) (bool, error) {
	var con TFACon = GetCon(viper)
	success, err := con.Validate(viper.GetBool("verbose"))
	return success, err
}

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
