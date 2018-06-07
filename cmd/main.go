package main

import (
	"context"
	"fmt"
	"go-project-builder/builder"
	"os"
	"os/signal"
	"syscall"

	"github.com/namsral/flag"
)

func main() {
	var config string = ""
	fs := flag.NewFlagSet("go-project-builder", flag.ContinueOnError)
	fs.StringVar(&config, "config", "", "toml file name")
	err := fs.Parse(os.Args[1:])
	if nil != err {
		fmt.Println(err)
		return
	}

	resolver, err := builder.NewResolver(context.Background(), config)
	if nil != err {
		fmt.Println(err)
		return
	}

	if err = signalNotify(resolver); nil != err {
		fmt.Println(err)
		return
	}

	err = resolver.Start()
	if nil != err {
		fmt.Println(err)
	}

	return
}

func signalNotify(resolver *builder.Resolver) error {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigChan
		signal.Stop(sigChan)
		fmt.Println(fmt.Sprintf("receive stop signal:%v, the programm will be quit.", sig))
		if nil != resolver {
			resolver.Stop()
		}
	}()
	return nil
}
