package main

import (
	//"flag"
	"fmt"
	"os"
	"bytes"
	"io"
)

const app_tra = "./app_tra"
const tra_app = "./tra_app"

func main() {

	fmt.Println("Opening app pipe for reading")
	stdoutRD, _ := os.OpenFile(tra_app, os.O_RDONLY, 0600)
	fmt.Println("Reading")

	var buffApp bytes.Buffer
	io.Copy(&buffApp, stdoutRD)
	msgApp := buffApp.String()

	fmt.Print("mensagem recebida:\t")
	fmt.Println(msgApp)
	stdoutRD.Close()
	
//----------
	
	fmt.Println("Opening app pipe for writing")
	stdoutWR, _ := os.OpenFile(app_tra, os.O_WRONLY, 0600)
	fmt.Println("Writing")
	
	stdoutWR.Write([]byte("hello"))
	
	fmt.Println("Message written")
	stdoutWR.Close()


}