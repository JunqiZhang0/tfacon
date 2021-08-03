package core

import (
	"net/http"

	"github.com/JunqiZhang0/tfacon/connectors"

	"github.com/JunqiZhang0/tfacon/common"

	"github.com/spf13/viper"
)

type TFACon interface {
	GetAllTestIds() []string
	BuildUpdatedList(ids []string) common.GeneralUpdatedList
	UpdateAll(common.GeneralUpdatedList)
	String() string
}

func Run(viper *viper.Viper) {
	var con TFACon = GetCon(viper)
	// fmt.Printf("%+v\n", con)
	// fmt.Println("===========================")
	ids := con.GetAllTestIds()
	updated_list_of_issues := con.BuildUpdatedList(ids)
	con.UpdateAll(updated_list_of_issues)
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
