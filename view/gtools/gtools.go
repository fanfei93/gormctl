package gtools

import (
	"gitlab.2345.cn/gomod/gormctl/config"
	"gitlab.2345.cn/gomod/gormctl/model"
	"gitlab.2345.cn/gomod/gormctl/model/genmysql"
)

func Execute(c config.Config) {
	option := model.GenModelOption{
		MysqlOption: genmysql.MysqlOption{
			Host:      c.Host,
			Port:      c.Port,
			Database:  c.Database,
			Username:  c.User,
			Password:  c.Password,
			TableName: c.Table,
		},
		PackageName: c.PackageName,
	}

	genModel := model.NewGenModel(&option)
	err := genModel.GenerateModelFile(c)
	if err != nil {
		panic(err)
	}

}
