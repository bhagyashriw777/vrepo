// Package agent provides functionality for initializing the agent.
package agent

import (
	"context"
	"log"
	"os"
	"os/exec"

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
	vncPort := "6081"
	vnc2Port := "6080"
	if os.Getenv("TERMINAL_PORT") != "" {
		terminalPort = os.Getenv("TERMINAL_PORT")
	}
	if os.Getenv("CODE_PORT") != "" {
		codePort = os.Getenv("CODE_PORT")
	}
	if os.Getenv("NO_VNC_PORT") != "" {
		vncPort = os.Getenv("NO_VNC_PORT")
	}
	if os.Getenv("AGENT_PORT") != "" {
		agentPort = os.Getenv("AGENT_PORT")
	}
	if os.Getenv("NO_VNC2_PORT") != "" {
		vnc2Port = os.Getenv("NO_VNC2_PORT")
	}
	var val string
	for _, val = range apps {
		switch val {
		case "code":
			vscode.Install()
			vscode.Start(codePort)
		// case "term":
		// 	fmt.Println("Terminal")
		case "vnc":
			//execute command to start vnc server
			log.Print("Starting VNC Server")
			out, err := exec.Command("/bin/sh", "-c", "pip install numpy").Output()
			if err != nil {
				log.Print("Error while installing numpy", err)
			}
			log.Printf("Output from numpy: %s\n", out)
			go func() {
				out, err = exec.Command("/bin/sh", "-c", "/opt/vnc/scripts/vncserver.sh").Output()
				if err != nil {
					log.Print("Error while running VNC Server", err)
				}
				log.Printf("Output from vnc server: %s\n", out)
			}()
			go func() {
				out, err = exec.Command("/bin/sh", "-c", "cd /opt/vnc/noVNC/utils/websockify/ && python3 -m websockify --web /opt/vnc/noVNC/ 6081 localhost:5990").Output()
				// out, err = exec.Command("/bin/sh", "nohup supervisord", "-c", " /etc/supervisor/supervisord.conf").Output()
				if err != nil {
					log.Print("Error while running NoHup Supervisord", err)
				}
				log.Printf("Output from supervisord: %s\n", out)
			}()
		case "vnc2":
			//execute command to start vnc server
			log.Print("Starting VNC Server")
			go func() {
				out, err := exec.Command("/bin/sh", "-c", "vncserver :1  -alwaysshared -SecurityTypes None -geometry 1280x1024 -depth 24").Output()
				if err != nil {
					log.Print("Error while running VNC Server", err)
				}
				log.Printf("Output from VNC server: %s\n", out)
			}()
			go func() {
				out, err := exec.Command("/bin/sh", "-c", "websockify --web=/usr/share/novnc/ 6080 localhost:5901").Output()
				// out, err = exec.Command("/bin/sh", "nohup supervisord", "-c", " /etc/supervisor/supervisord.conf").Output()
				if err != nil {
					log.Print("Error while running websockify", err)
				}
				log.Printf("Output from websockify: %s\n", out)
			}()

		}
	}

	go goTTYwin(terminalPort)
	reverseproxy.Serve(terminalPort, codePort, agentPort, vncPort, vnc2Port)

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
			"command":  "powershell",
			"hostname": "localhost",
		},
		// IndexFile: "agent/gotty/resources/index.html",
	}
	fact, _ := localcommand.NewFactory("powershell", []string{"-c", "powershell"}, &localcommand.Options{})
	srv, err := server.New(fact, appOptions)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	srv.Run(context.Background())
}
