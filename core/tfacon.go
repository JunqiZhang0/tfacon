package core

import (
	"tfactl/common"
	"tfactl/connectors"

	"github.com/spf13/viper"
)

type TFACon interface {
	UpdateAll(common.UpdatedList)
	GetAllTestIds() []string
	BuildUpdatedList(ids []string) common.UpdatedList
}

func Run(viper *viper.Viper) {
	var con TFACon = GetCon(viper)
	ids := con.GetAllTestIds()
	updated_list_of_issues := con.BuildUpdatedList(ids)
	con.UpdateAll(updated_list_of_issues)
}

func GetCon(viper *viper.Viper) TFACon {
	var con TFACon
	switch viper.Get("ConnectorType") {
	case "RPCon":
		con = &connectors.RPConnector{}
		viper.Unmarshal(con)
		// case "POLCon":
		// 	con = RPConnector{}
		// case "JiraCon":
		// 	con = RPConnector{}

	}
	return con
}
