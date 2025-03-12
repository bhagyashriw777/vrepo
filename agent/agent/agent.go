// Package agent provides functionality for initializing the agent.
package agent

import (
	"context"
	"log"
	"os"

	"github.com/bhagyashriw777/tty2web/backend/localcommand"
	"github.com/bhagyashriw777/tty2web/server"

	vscode "github.com/bhagyashriw777/vrepo/agent/executer"
	reverseproxy "github.com/bhagyashriw777/vrepo/agent/reverse_proxy"
)

// Init initializes the agent.
func Init(apps []string) {

	terminalPort := "8000"
	codePort := "3000"
	agentPort := "9000"
	if os.Getenv("TERMINAL_PORT") != "" {
		terminalPort = os.Getenv("TERMINAL_PORT")
	}
	if os.Getenv("CODE_PORT") != "" {
		codePort = os.Getenv("CODE_PORT")
	}
	if os.Getenv("AGENT_PORT") != "" {
		agentPort = os.Getenv("AGENT_PORT")
	}
	var val string
	for _, val = range apps {
		switch val {
		case "code":
			vscode.Install()
			vscode.Start(codePort)
			// case "term":
			// 	fmt.Println("Terminal")
		}
	}

	go goTTYwin(terminalPort)
	reverseproxy.Serve(terminalPort, codePort, agentPort)

}
func goTTYwin(terminalPort string) {
	appOptions := &server.Options{
		PermitWrite:     true,
		Address:         "localhost",
		Port:            terminalPort,
		TitleFormat:     "V-Collab",
		EnableReconnect: true,
		ReconnectTime:   10,
		TitleVariables: map[string]interface{}{
			"command":  "cmd",
			"hostname": "localhost",
		},
		// IndexFile: "agent/gotty/resources/index.html",
	}
	fact, _ := localcommand.NewFactory("cmd", []string{"-c", "cmd"}, &localcommand.Options{})
	srv, err := server.New(fact, appOptions)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	srv.Run(context.Background())
}
