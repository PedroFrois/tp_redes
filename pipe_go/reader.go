package main

import "fmt"
import "flag"
import "os/exec"
import "os"
import "io"
import "path/filepath"
import "syscall"
import "bytes"
import "io/ioutil"

var writerPath string

//Código cria um pipe num diretorio temporario. Executa o código do writer.go e lê do pipe
func main() {
	flag.StringVar(&writerPath, "writer", "./writer", "path to writer")
	flag.Parse()

	tmpDir, _ := ioutil.TempDir("", "named-pipes")

	// Create named pipe
	namedPipe := filepath.Join(tmpDir, "stdout")
	syscall.Mkfifo(namedPipe, 0600)

	go func() {
		cmd := exec.Command(writerPath, namedPipe)
		// Just to forward the stdout
		cmd.Stdout = os.Stdout
		cmd.Run()
	}()

	// Open named pipe for reading
	fmt.Println("Opening named pipe for reading")
	stdout, _ := os.OpenFile(namedPipe, os.O_RDONLY, 0600)
	//OpenFile(nome do pipe, leitura ou escrita, permissão do arquivo que será criado caso já não exista)
	fmt.Println("Reading")

	var buff bytes.Buffer
	fmt.Println("Waiting for someone to write something")
	io.Copy(&buff, stdout)
	stdout.Close()
	fmt.Printf("Data: %s\n", buff.String())
}
