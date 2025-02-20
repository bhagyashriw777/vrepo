// Package executer runs command received as input to Run() function and returns the output
package executer

import (
	"log"
	"os/exec"
	"time"
)

// Install runs the command to install code-server and retries until it succeeds
func Install() {
	const maxRetries = 5
	const retryDelay = 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		var out []byte
		out, err := exec.Command("/bin/sh", "-c", "curl -fsSL https://code-server.dev/install.sh | sh -s -- --method=standalone --prefix=/var/tmp/code-server").Output()
		if err != nil {
			log.Printf("Attempt %d: Installation failed: %v", i+1, err)
			time.Sleep(retryDelay)
			continue
		}
		log.Printf("Output: %s\n", out)
		log.Println("Installation succeeded")
		return
	}
	log.Fatal("Installation failed after maximum retries")
}

// Start runs the command to start code-server
func Start(code_port string) {

	err := exec.Command("/bin/sh", "-c", "/var/tmp/code-server/bin/code-server --auth none --port "+code_port+" $HOME >/var/tmp/code-server.log 2>&1 &").Start()
	if err != nil {
		log.Fatal(err)
	}
}
