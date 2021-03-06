package main

import (
	"github.com/rfparedes/gdg/action"
	"github.com/rfparedes/gdg/setup"
)

func main() {

	//SupportedBinaries := make([]string, 3)
	setup.CreateOrLoadConfig()
	action.Gather()

}
