package main

import (
	"fmt"

	"github.com/rfparedes/gdg/setup"
)

func main() {

	var s []string
	//SupportedBinaries := make([]string, 3)
	s = setup.FindSupportedBinaries()
	setup.CreateOrLoadConfig()
	fmt.Println(s)
}
