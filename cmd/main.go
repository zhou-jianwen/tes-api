package main

import (
	"fmt"
	"os"

	"github.com/GBA-BI/tes-api/internal/apiserver"
	"github.com/GBA-BI/tes-api/third_party/forked/sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

//	@title			volcengine TES apiserver
//	@version		1.0
//	@description	This is volcengine TES apiserver using Hertz.

//	@BasePath	/
//	@schemes	http

//	@query.collection.format	multi

func main() {
	ctx := signals.SetupSignalHandler()
	command, err := apiserver.NewServerCommand(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	if err = command.Execute(); err != nil {
		os.Exit(1)
	}
}
