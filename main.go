package main

import (
	"github.com/marekjalovec/steampipe-plugin-make/make"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: make.Plugin})
}
