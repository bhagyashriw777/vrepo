// Package agent provides functionality for initializing the agent.
package agent

import (
	"context"
	"log"
	"os"

	"github.com/bhagyashriw777/tty2web/backend/localcommand"
	"github.com/bhagyashriw777/tty2web/server"

	reverseproxy "github.com/bhagyashriw777/vrepo/agent/reverse_proxy"
)

// Init initializes the agent.
func Init(apps []string) {

	terminalPort := "8000"
	agentPort := "9000"
	if os.Getenv("TERMINAL_PORT") != "" {
		terminalPort = os.Getenv("TERMINAL_PORT")
	}
	if os.Getenv("AGENT_PORT") != "" {
		agentPort = os.Getenv("AGENT_PORT")
	}

	go goTTYwin(terminalPort)
	reverseproxy.Serve(terminalPort, agentPort)

}
func goTTYwin(terminalPort string) {
	appOptions := &server.Options{
		PermitWrite:     true,
		Address:         "localhost",
		Port:            terminalPort,
		TitleFormat:     "V-Collab",
		EnableReconnect: true,
		Url:             "term",
		ReconnectTime:   10,
		TitleVariables: map[string]interface{}{
			"command":  "powershell.exe",
			"hostname": "localhost",
		},
		// IndexFile: "agent/gotty/resources/index.html",
	}
	fact, _ := localcommand.NewFactory("powershell.exe", []string{"-c", "powershell.exe"}, &localcommand.Options{})
	srv, err := server.New(fact, appOptions)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	srv.Run(context.Background())
}
