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
	
	fmt.Println("Opening network pipe for reading")
	stdoutRD, _ := os.OpenFile(tra_net, os.O_RDONLY, 0600)
	fmt.Println("Reading")

	var buffNetwork bytes.Buffer
	io.Copy(&buffNetwork, stdoutRD)
	msgNetwork:= buffNetwork.String()

	fmt.Print("mensagem recebida:\t")
	fmt.Println(msgNetwork)

//---------
	
	fmt.Println("Opening network pipe for writing")
	stdoutWR, _ := os.OpenFile(net_tra, os.O_WRONLY, 0600)
	fmt.Println("Writing")
	stdoutWR.Write([]byte("005870600000000000010000000001000501101060635acknowledge"))
	fmt.Println("Message written")
	
//----------

	fmt.Println("Reading")

	io.Copy(&buffNetwork, stdoutRD)
	msgNetwork = buffNetwork.String()

	fmt.Print("mensagem recebida:\t")
	fmt.Println(msgNetwork)

//---------
	
	fmt.Println("Writing")
	stdoutWR.Write([]byte("005870600000000000010000000001000501101060635conneccao ok"))
	fmt.Println("Message written")
	
//----------

	fmt.Println("Reading")

	io.Copy(&buffNetwork, stdoutRD)
	msgNetwork = buffNetwork.String()

	fmt.Print("mensagem recebida:\t")
	fmt.Println(msgNetwork)

//---------
	
	fmt.Println("Writing")
	stdoutWR.Write([]byte("005870600000000000010000000001000501101060635email ok"))
	fmt.Println("Message written")
	
//----------

	fmt.Println("Reading")

	io.Copy(&buffNetwork, stdoutRD)
	msgNetwork = buffNetwork.String()

	fmt.Print("mensagem recebida:\t")
	fmt.Println(msgNetwork)

//---------
	
	fmt.Println("Writing")
	stdoutWR.Write([]byte("005870600000000000010000000001000501100160635email ok"))
	fmt.Println("Message written")
	
//----------

	stdoutRD.Close()
	stdoutWR.Close()
}