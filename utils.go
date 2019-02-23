package main

import (
	"os"
	"os/exec"
	"regexp"
)

//Masscan outputs some text along with IP addr, this function strips unwanted text.
func filterIP(input string) string {
	numBlock := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
	regexPattern := numBlock + "\\." + numBlock + "\\." + numBlock + "\\." + numBlock

	regEx := regexp.MustCompile(regexPattern)
	return regEx.FindString(input)
}

// Function to check if masscan is installed
func masscanInstalled() bool {
	cmd := exec.Command("/bin/bash", "-c", "sudo masscan -v ")
	_, err := cmd.StdoutPipe()

	if err != nil {
		print(err)
		print(`Masscan not found!
If you are running Ubuntu or Kali linux, Install Masscan by running -

sudo apt install masscan

For other disctributions check masscan's git repo for install instructions - 
https://github.com/robertdavidgraham/masscan
`)
		os.Exit(0)
		return false
	}
	return true
}

// WelcomeMsg prints welcome msg :D (Go-Lint compatabilty).
func WelcomeMsg() {
	print(`
███╗   ███╗ ██████╗ ███╗   ██╗ ██████╗  ██████╗       ██████╗ ██╗   ██╗███████╗████████╗███████╗██████╗ 
████╗ ████║██╔═══██╗████╗  ██║██╔════╝ ██╔═══██╗      ██╔══██╗██║   ██║██╔════╝╚══██╔══╝██╔════╝██╔══██╗
██╔████╔██║██║   ██║██╔██╗ ██║██║  ███╗██║   ██║█████╗██████╔╝██║   ██║███████╗   ██║   █████╗  ██████╔╝
██║╚██╔╝██║██║   ██║██║╚██╗██║██║   ██║██║   ██║╚════╝██╔══██╗██║   ██║╚════██║   ██║   ██╔══╝  ██╔══██╗
██║ ╚═╝ ██║╚██████╔╝██║ ╚████║╚██████╔╝╚██████╔╝      ██████╔╝╚██████╔╝███████║   ██║   ███████╗██║  ██║
╚═╝     ╚═╝ ╚═════╝ ╚═╝  ╚═══╝ ╚═════╝  ╚═════╝       ╚═════╝  ╚═════╝ ╚══════╝   ╚═╝   ╚══════╝╚═╝  ╚═╝
																																			   

Started Scannig servers.
`)
}
