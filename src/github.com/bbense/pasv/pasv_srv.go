package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

/*

This program is meant to run under remtcl ( See https://github.com/rra/remctl )

It will take values from the command line and the remctl environment to write a status
message to the nagios command file.

( From http://nagios.sourceforge.net/docs/3_0/passivechecks.html )

External applications can submit passive service check results to Nagios by writing a
PROCESS_SERVICE_CHECK_RESULT external command to the external command file.

The format of the command is as follows:

[<timestamp>] PROCESS_SERVICE_CHECK_RESULT;<hostName>;<svc_description>;<return_code>;<plugin_output>

where...

timestamp is the time in time_t format (seconds since the UNIX epoch) that the service check
was perfomed (or submitted). Please note the single space after the right bracket.

hostName is the short name of the host associated with the service in the service definition

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

/* Need to duplicate this code from nsca


CMT: writes service/host check results to the Nagios command file
static int write_check_result(char *hostName, char *svc_description, int return_code, char *plugin_output, time_t check_time){

        if(aggregate_writes==FALSE){
                if(open_command_file()==ERROR)
                        return ERROR;
                }

        if(!strcmp(svc_description,""))
                fprintf(command_file_fp,"[%lu] PROCESS_HOST_CHECK_RESULT;%s;%d;%s\n",(unsigned long)check_time,hostName,return_code,plugin_output);
        else
                fprintf(command_file_fp,"[%lu] PROCESS_SERVICE_CHECK_RESULT;%s;%s;%d;%s\n",(unsigned long)check_time,hostName,svc_description,return_code,plugin_output);

        if(aggregate_writes==FALSE)
                close_command_file();
        else
                CMT: if we don't fflush() then we're writing in 4k non-CR-terminated blocks, and
                 * anything else (eg. pscwatch) which writes to the file will be writing into
                 * the middle of our commands.

                fflush(command_file_fp);

        return OK;
        }



CMT: opens the command file for writing
static int open_command_file(void){
        struct stat statbuf;

        CMT: file is already open
        if(command_file_fp!=NULL)
                return OK;

        CMT: command file doesn't exist - monitoring app probably isn't running...
        if(stat(command_file,&statbuf)){

                if(debug==TRUE)
                        syslog(LOG_ERR,"Command file '%s' does not exist, attempting to use alternate dump file '%s' for output",command_file,alternate_dump_file);

                CMT: try and write checks to alternate dump file
                command_file_fp=fopen(alternate_dump_file,"a");
                if(command_file_fp==NULL){
                        if(debug==TRUE)
                                syslog(LOG_ERR,"Could not open alternate dump file '%s' for appending",alternate_dump_file);
                        return ERROR;
                        }

                return OK;
                }

        CMT: open the command file for writing or appending
        command_file_fp=fopen(command_file,(append_to_file==TRUE)?"a":"w");
        if(command_file_fp==NULL){
                if(debug==TRUE)
                        syslog(LOG_ERR,"Could not open command file '%s' for %s",command_file,(append_to_file==TRUE)?"appending":"writing");
                return ERROR;
                }

        return OK;
        }



CMT: closes the command file
static void close_command_file(void){

        fclose(command_file_fp);
        command_file_fp=NULL;

        return;
        }



*/

func sendPasv(cmdFile string, message string) bool {

	fd, err := os.OpenFile(cmdFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer fd.Close()

	if _, err = fd.WriteString(message); err != nil {
		panic(err)
	}
	return true
}

func getRemuser() string {
	remuser := os.Getenv("REMUSER")
	if len(remuser) == 0 {
		panic(errors.New("REMUSER not set"))
	}
	//REMUSER should look like host/foobar.slac.stanford.edu
	if strings.Index(remuser, "/") > -1 {
		remuser = strings.Split(remuser, "/")[1]
	}
	return remuser
}

func getHost() string {
	//First see if we have REMOTE_HOST
	remhost := os.Getenv("REMOTE_HOST")
	if len(remhost) == 0 {
		remhost = getRemuser()
	}
	//We should proably complain loudly at this point if we
	//don't have a string that looks like a FQDN.
	// Trim domain
	result := strings.Split(remhost, ".")[0]
	return result
}

//Read results from cmdLine
func getService(cmdline []string) (service string, code int, message string) {
	var err error
	//Validate data

	service = cmdline[0]
	code, err = strconv.Atoi(cmdline[1])
	if err != nil || code < 0 || code > 4 {
		if code < 0 || code > 4 {
			err = fmt.Errorf("nagios check code value is invalid: %d", code)
		}
		panic(err)
	}
	rest := cmdline[2:]
	message = strings.Join(rest, " ")
	return
}

/*
Need to return this string
"[%lu] PROCESS_SERVICE_CHECK_RESULT;%s;%s;%d;%s\n",(unsigned long)check_time,hostName,svc_description,return_code,plugin_output);
*/

func getAlert(args []string) string {
	epoch := int32(time.Now().Unix()) // This is 32 bit seconds since Unix epoch.
	hostName := getHost()
	service, code, message := getService(args)
	alert := fmt.Sprintf("[%d] PROCESS_SERVICE_CHECK_RESULT;%s;%s;%d;%s\n", epoch, hostName, service, code, message)
	return alert
}

// We need to find nagios cmd file, read config file or cmdline arg?
var (
	verbose = flag.Bool("verbose", false, "Verbose")
	cmdfile = flag.String("cmd", "/var/nagios/rw/nagios.cmd", "Path to nagios command file")
	fixed_host = flag.String("fixhost", "", 
		"Used fixed_host name in alert, rather than what remctl provides in the environment")
)

func main() {

	flag.Parse()
	
	if len(*fixed_host) > 0 {
		os.Setenv("REMOTE_HOST",*fixed_host)
	}
	// Read status and message from ARGV.
	alert := getAlert(flag.Args())
	sendPasv(*cmdfile, alert)
	if *verbose {
		fmt.Printf("Sent alert: %q\n", alert)
	}
}
