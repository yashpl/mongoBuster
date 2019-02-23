package main

import (
	"bufio"
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os/exec"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func main() {
	ipAddr := make(chan string)
	data := make(chan string)

	masscanInstalled() // Check if masscan binary is installed.
	WelcomeMsg()

	maxPtr := flag.String("max-rate", "100", "Max rate at which packets will be sent")
	outFile := flag.String("out-file", "IPs.log", "Name of file to which vulnerable IPs will be exported")
	flag.Parse()

	go execMasscan(ipAddr, maxPtr)
	go fileWriter(data, outFile)
	workDispatcher(ipAddr, data) // Dont call this func inside execMasscan coz exec.Command is a blocking statement.
}

func execMasscan(ipAddr chan string, maxPtr *string) {
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
	_, err = client.ListDatabaseNames(context.TODO(), bson.D{{}})

	if err != nil {
		print("\r\033[K" + input + ": ")
		println(err.Error())
	} else {
		println("\r\033[K" + input + " is VULNERABLE")
		println("")
		data <- input

	}

}

func fileWriter(data chan string, outFile *string) {

	for value := range data {

		toWrite := []byte(value)
		err := ioutil.WriteFile(*outFile, toWrite, 0644)
		check(err)
	}

}
