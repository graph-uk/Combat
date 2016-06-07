package LauncherServer

import (
	"net/http"
	"strconv"

	"github.com/graph-uk/Combat/CombatWorkerLauncher/LauncherServer/config"
)

type LauncherServer struct {
	config *config.Config
}

func NewLauncherServer() (*LauncherServer, error) {
	var result LauncherServer
	var err error
	result.config, err = config.LoadConfig()
	if err != nil {
		return &result, err
	}

	return &result, nil
}

func (t *LauncherServer) Serve() error {
	http.HandleFunc("/LaunchWorkers/", t.launchWorkersHandler)
	err := http.ListenAndServe(":"+strconv.Itoa(t.config.Port), nil)
	return err
}
