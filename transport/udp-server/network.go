package main

import (
	//"flag"
	"fmt"
	"os"
	"bytes"
	"io"
)

const tra_net = "./tra_net"
const net_tra = "./net_tra"

func main() {

	fmt.Println("Opening network pipe for writing")
	stdoutWR, _ := os.OpenFile(net_tra, os.O_WRONLY, 0600)
	fmt.Println("Writing")
	stdoutWR.Write([]byte("00587065000000560808goodbye"))
	fmt.Println("Message written")
	stdoutWR.Close()

//---------
	
	fmt.Println("Opening network pipe for reading")
	stdoutRD, _ := os.OpenFile(tra_net, os.O_RDONLY, 0600)
	fmt.Println("Reading")

	var buffNetwork bytes.Buffer
	io.Copy(&buffNetwork, stdoutRD)
	msgNetwork:= buffNetwork.String()

	fmt.Print("mensagem recebida:\t")
	fmt.Println(msgNetwork)
	stdoutRD.Close()
}