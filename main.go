package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"swarm/bundle"
	"swarm/config"
	"swarm/monitor"
	"swarm/source"
	"swarm/web"
	"syscall"
)

func main() {
	fmt.Print("\nSwarm welcomes you.\n\n")
	log.SetOutput(os.Stdout)

	// configuration
	swarmConfig, err := config.TryLoadSwarmConfigFromCWD()
	if err != nil {
		log.Fatalf("Failed to load swarm.json file: %s", err)
		return
	}
	runtimeConfig := chooseBuild(swarmConfig.Builds)
	moduleDescrs, err := config.LoadBuildDescriptionFile(runtimeConfig.BuildPath)
	if err != nil {
		log.Fatalf("Failed to load build description file: '%s'", runtimeConfig.BuildPath)
		return
	}

	var server *web.Server

	// monitor
	ws := source.NewWorkspace(swarmConfig.RootPath)
	mon := monitor.NewMonitor(ws, swarmConfig.Monitor)
	moduleSet := bundle.CreateModuleSet(
		ws,
		moduleDescrs.NormaliseModules(ws.RootPath()),
		runtimeConfig,
	)
	go mon.NotifyOnChanges(func(changes *monitor.EventChangeset) {
		moduleSet.NotifyChanges(changes)
		server.NotifyReload()
	})
	moduleSet.NotifyChanges(nil) // trigger first time build

	// web server
	handlers := moduleSet.GenerateHTTPHandlers()
	server = web.CreateServer(swarmConfig.RootPath, &web.ServerOptions{
		Port:     swarmConfig.Server.Port,
		Handlers: handlers,
	})
	go server.Start()
	fmt.Printf("Listening on http://localhost:%d\n", server.Port())

	// sleep
	waitForExit(server)
}

func chooseBuild(builds map[string]*config.RuntimeConfig) *config.RuntimeConfig {
	// build specified as first command line argument
	if len(os.Args) > 1 {
		buildName := os.Args[1]
		if build, found := builds[buildName]; found {
			return build
		}
	}

	// single build?
	if len(builds) == 1 {
		for k := range builds {
			return builds[k]
		}
	}

	fmt.Println("Please choose a build:")
	for {
		buildNames := make([]string, len(builds))
		i := 0
		for name := range builds {
			buildNames[i] = name
			i++
		}
		sort.Strings(buildNames)

		appmap := map[string]*config.RuntimeConfig{}
		for i, name := range buildNames {
			fmt.Printf("  %d) %s\n", (i + 1), name)
			appmap[strconv.Itoa(i+1)] = builds[name]
			appmap[name] = builds[name]
			i++
		}
		fmt.Println("  -----------------------")
		fmt.Print("  >")
		reader := bufio.NewReader(os.Stdin)
		lineBytes, _, err := reader.ReadLine()
		if err != nil {
			log.Fatal("Bad input")
		}
		if build, found := appmap[string(lineBytes)]; found {
			return build
		}
	}
}

func waitForExit(server *web.Server) {

	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal

	server.Stop()
}
