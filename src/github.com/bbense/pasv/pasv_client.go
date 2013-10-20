package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// We will need three things from the command line, the status, the service and the message.
var (
	nstatus  = flag.Int("status", 0, "Status value for nagios service check")
	nservice = flag.String("checkname", "default", "Name of service check")
)

// We should have defaults for the remctl stuff.

func main() {

	flag.Parse()
	message := strings.Join(flag.Args(), " ")

	fmt.Printf("Sending %s : %d : %s\n", *nservice, *nstatus, message)
	cmd := exec.Command("tr", "a-z", "A-Z")
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("in all caps: %q\n", out.String())
}
