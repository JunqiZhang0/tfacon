package core

import (
	"net/http"

	"github.com/JunqiZhang0/tfacon/connectors"

	"github.com/JunqiZhang0/tfacon/common"

	"github.com/spf13/viper"
)

type TFACon interface {
	GetAllTestIds() []string
	BuildUpdatedList(ids []string, concurrent bool) common.GeneralUpdatedList
	UpdateAll(common.GeneralUpdatedList)
	String() string
	InitConnector()
	Validate() error
}

func Run(viperRun, viperConfig *viper.Viper) {
	var con TFACon = GetCon(viperRun)
	// fmt.Printf("%+v\n", con)
	// fmt.Println("===========================")
	con.InitConnector()
	runHelper(viperRun, viperConfig, con.GetAllTestIds(), con)

}

func runHelper(viperRun, viperConfig *viper.Viper, ids []string, con TFACon) {
	if len(ids) == 0 {
		return
	}
	updated_list_of_issues := con.BuildUpdatedList(ids, viperConfig.GetBool("config.concurrency"))
	// Doing this because the api can only take 20 items per request

	con.UpdateAll(updated_list_of_issues)
	runHelper(viperRun, viperConfig, con.GetAllTestIds(), con)
}

func GetInfo(viper *viper.Viper) TFACon {
	var con TFACon = GetCon(viper)
	return con
}

func GetCon(viper *viper.Viper) TFACon {
	var con TFACon
	switch viper.Get("CONNECTOR_TYPE") {
	case "RPCon":
		con = &connectors.RPConnector{Client: &http.Client{}}
		viper.Unmarshal(con)
		// case "POLCon":
		// 	con = RPConnector{}
		// case "JiraCon":
		// 	con = RPConnector{}

	}
	return con
}
