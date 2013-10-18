package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

/*

This program is meant to run under remtcl ( See https://github.com/rra/remctl )

It will take values from stdin and the remctl environment to write a status
message to the nagios command file.

( From http://nagios.sourceforge.net/docs/3_0/passivechecks.html )

External applications can submit passive service check results to Nagios by writing a 
PROCESS_SERVICE_CHECK_RESULT external command to the external command file.

The format of the command is as follows:

[<timestamp>] PROCESS_SERVICE_CHECK_RESULT;<host_name>;<svc_description>;<return_code>;<plugin_output>

where...

timestamp is the time in time_t format (seconds since the UNIX epoch) that the service check 
was perfomed (or submitted). Please note the single space after the right bracket.

host_name is the short name of the host associated with the service in the service definition

svc_description is the description of the service as specified in the service definition

return_code is the return code of the check (0=OK, 1=WARNING, 2=CRITICAL, 3=UNKNOWN)

plugin_output is the text output of the service check (i.e. the plugin output)

(From http://www.eyrie.org/~eagle/software/remctl/remctld.html )

ENVIRONMENT

The following environment variables will be set for any commands run via remctld:

REMOTE_USER
REMUSER
Set to the Kerberos principal of the authenticated client. REMUSER has always been set by remctld; 
REMOTE_USER is also set (to the same value) starting with remctl 2.1.

REMOTE_ADDR
The IP address of the remote host. Currently, this is always an IPv4 address, but in the future 
it may be set to an IPv6 address. This environment variable was added in remctl 2.1.

REMOTE_HOST
The hostname of the remote host, if it was available. If reverse name resolution failed, 
this environment variable will not be set. This variable was added in remctl 2.1.

REMCTL_COMMAND
The command string that caused this command to be run. This variable will contain only 
the command, not the subcommand or any additional arguments (which are passed as command arguments). 
This variable was added in remctl 2.16.

*/

func main() {
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
