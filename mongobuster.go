package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"os/exec"
	"regexp"

	"github.com/mongodb/mongo-go-driver/mongo"
)

func main() {
	ipAddr := make(chan string)

	masscanInstalled() // Check if masscan binary is installed.
	WelcomeMsg()

	go execMasscan(ipAddr)
	workDispatcher(ipAddr) // Dont call this func inside execMasscan coz exec.Command is a blocking statement.
}

func execMasscan(ipAddr chan string) {
	cmd := exec.Command("/bin/bash", "-c", "sudo masscan -p27017 0.0.0.0/0 --exclude 255.255.255.255 --open-only")
	ok, err := cmd.StdoutPipe()

	if err != nil {
		log.Fatalf("could not get stdout pipe: %v", err)
	}

	go func() {
		scanner := bufio.NewScanner(ok)

		for scanner.Scan() {
			msg := scanner.Text()
			ipAddr <- msg
		}
		close(ipAddr)
	}()
	if err := cmd.Run(); err != nil {
		log.Fatalf("could not run cmd: %v", err)
	}
	if err != nil {
		log.Fatalf("could not wait for cmd: %v", err)
	}
}

func workDispatcher(ipAddr chan string) {
	num := 0
	for value := range ipAddr {
		num++
		print("\rTotal MongoDB servers found - ", num)
		go testIP(filterIP(value))
	}
}

//Masscan outputs some text along with IP addr, this function strips unwanted text.
func filterIP(input string) string {
	numBlock := "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
	regexPattern := numBlock + "\\." + numBlock + "\\." + numBlock + "\\." + numBlock

	regEx := regexp.MustCompile(regexPattern)
	return regEx.FindString(input)
}

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

func testIP(input string) {

	client, err := mongo.Connect(context.TODO(), "mongodb://"+input)

	if err != nil {
		// log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		// log.Fatal(err)
	} else {
		println("")
		println(input + " is VULNERABLE")
		println("")
	}

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
