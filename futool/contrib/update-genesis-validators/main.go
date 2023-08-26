package main

import (
	"github.com/percosis-labs/fury/app"
	"github.com/percosis-labs/futool/contrib/update-genesis-validators/cmd"
)

func main() {
	app.SetSDKConfig()
	cmd.Execute()
}
