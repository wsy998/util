package main

import (
	"os"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"

	"util/internal/cmd"
	_ "util/internal/logic"
	"util/internal/model"
	"util/internal/service"
)

func main() {
	ctx := gctx.GetInitCtx()
	defer func() {

		err := service.Os().KillAll()
		if err != nil {
			g.Log().Fatal(ctx, err)
		}
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				g.Log().Fatal(ctx, err)
			}
		}
	}()
	var config *model.TypeConfig
	configPath := gfile.Join(gfile.Pwd(), ".de-config.json")
	if gfile.Exists(configPath) {
		if json, err := gjson.DecodeToJson(gfile.GetBytes(configPath)); err != nil {
			panic(err)
		} else {
			if err := json.Scan(&config); err != nil {
				panic(err)
			}
		}
	} else {
		config = new(model.TypeConfig)
	}
	var err error
	switch os.Args[1] {
	case "dev":
		err = cmd.RunDev(ctx, config)
	case "build":
		err = cmd.RunBuild(ctx, config)
	case "go":
		err = cmd.RunGo(ctx, config)
	case "fe":
		err = cmd.RunFe(ctx, config)
	default:
		err = cmd.RunAlias(ctx, config)
	}
	if err != nil {
		panic(err)
	}
}
