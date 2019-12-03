package main

import (
	"fmt"
	"os"
	"bytes"
	"io"
)

const app_tra = "./app_tra"
const tra_app = "./tra_app"

func main() {
	
	fmt.Println("Opening app pipe for writing")
	stdoutWR, _ := os.OpenFile(app_tra, os.O_WRONLY, 0600)
	fmt.Println("Writing")
	
	stdoutWR.Write([]byte("parte1"))
	
	fmt.Println("Message written")

//----------
	
	fmt.Println("Opening app pipe for reading")
	stdoutRD, _ := os.OpenFile(tra_app, os.O_RDONLY, 0600)
	fmt.Println("Reading")

	var buffApp bytes.Buffer
	io.Copy(&buffApp, stdoutRD)
	msgApp := buffApp.String()

	fmt.Print("mensagem recebida:\t")
	fmt.Println(msgApp)

//----------

	fmt.Println("Writing")
	stdoutWR.Write([]byte("parte2"))
	fmt.Println("Message written")

//----------

	fmt.Println("Reading")

	io.Copy(&buffApp, stdoutRD)
	msgApp = buffApp.String()

	fmt.Print("mensagem recebida:\t")
	fmt.Println(msgApp)

//----------

	fmt.Println("Writing")
	stdoutWR.Write([]byte("corpo do email"))
	fmt.Println("Message written")

//----------

	fmt.Println("Reading")

	io.Copy(&buffApp, stdoutRD)
	msgApp = buffApp.String()

	fmt.Print("mensagem recebida:\t")
	fmt.Println(msgApp)

//----------

	fmt.Println("Writing")
	stdoutWR.Write([]byte("finalizacao parte1"))
	fmt.Println("Message written")

//----------

	fmt.Println("Reading")

	io.Copy(&buffApp, stdoutRD)
	msgApp = buffApp.String()

	fmt.Print("mensagem recebida:\t")
	fmt.Println(msgApp)

//----------

	fmt.Println("Writing")
	stdoutWR.Write([]byte("finalizacao parte2"))
	fmt.Println("Message written")

//----------

	fmt.Println("Reading")

	io.Copy(&buffApp, stdoutRD)
	msgApp = buffApp.String()

	fmt.Print("mensagem recebida:\t")
	fmt.Println(msgApp)

//----------

	stdoutRD.Close()
	stdoutWR.Close()
}
