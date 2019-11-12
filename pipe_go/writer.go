package main

import "os"
import "fmt"
import "flag"

func main() {
	flag.Parse()
 	namedPipe := flag.Args()[0]
 
 	fmt.Println("Opening named pipe for writing")
 	stdout, _ := os.OpenFile(namedPipe, os.O_RDWR, 0600)
 	fmt.Println("Writing")
 	stdout.Write([]byte("hello"))
 	stdout.Close()
}
