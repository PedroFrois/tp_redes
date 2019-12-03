package main

import (
	"fmt"
	"strconv"
	"time"

	//"flag"
	"io"
	//"io/ioutil"
	"os"
	"os/exec"
	//"path/filepath"
	"syscall"
	"bytes"
)

const (
	SmtpPort = 587
	ImapPort = 993
	Pop3Port = 995
)

// PDU - UDP
type UdpPdu struct {
	srcPort		uint16
	dstPort		uint16
	size		uint16
	checkSum	uint16
	msg			string
}

func NewUdpPdu (src, dst, size, check, msg string) UdpPdu {

	var u uint64
	var pdu UdpPdu

	u, _ = strconv.ParseUint(src, 10, 16)
	pdu.srcPort = uint16(u)

	u, _ = strconv.ParseUint(dst, 10, 16)
	pdu.dstPort = uint16(u)

	u, _ = strconv.ParseUint(size, 10, 16)
	pdu.size = uint16(u)

	u, _ = strconv.ParseUint(check, 10, 16)
	pdu.checkSum = uint16(u)

	pdu.msg = msg

	return pdu
}

func stringToBin(s string) (binString string) {
    res := ""
    for _, c := range s {
        res = fmt.Sprintf("%s%.8b", res, c)
    }
    return res
}

func calcCheckSum (pdu UdpPdu) uint16 {
	//http://www.jvasconcellos.com.br/fat/FAT_TI/wp-content/uploads/2013/10/checksum-udp.pdf
	//1. determinação do comprimento do datagrama;
	//pdu.size ja fornece o tamanho dos dados em byte

	//2. agrupamento dos campos em blocos de 16 bits;
	//16 bits = 2 bytes = 2 caracteres
	//Todos os campos com excessao da msg ja sao inteiros binarios de 16 bits

	//3. conversão da string menssagem para números binários de 16 bits;
    blocks := []byte(pdu.msg)
    // fmt.Println(blocks)
    var i, j uint16
    j = 0

    var m [32767]uint16

    for i = 0; i < pdu.size-1; i=i+2 {
    	m[j] = uint16(((uint16(blocks[i+1]))*256) | uint16(blocks[i]))
    	// fmt.Println(m[j])
    	j++
    }

	//4. soma dos blocos de 16 bits;
	var sum uint32
	sum = uint32(pdu.srcPort + pdu.dstPort + pdu.size)

	var resp uint16

	for i = 0; i < j; i++ {
		sum = sum+ uint32(m[i])
	}
	if(sum > 65536) {
		resp = uint16(sum)
		resp = resp + uint16(sum >> 16)
	} else {
		resp = uint16(sum)
	}

	//5. complemento de um da soma dos blocos de 16 bits.
	//Ao fazer o complemento de 1, invertemos todos os bits do numero
	//A soma de verificacao eh o complemento de 1 invertido
	//Entao a soma de verificacao eh o resultado ja encontrado

	// fmt.Println(sum)
	return resp
}

func checkCheckSum(pdu UdpPdu) bool {
	sum := calcCheckSum(pdu)

	if sum == pdu.checkSum {
		return true
	} else {
		return false
	}
}

func makeMsgNetwork(pdu UdpPdu) string {
	var aux, src, dst, size, check string

	aux = strconv.FormatInt(int64(pdu.srcPort), 10)
	if(pdu.srcPort > 10000) {
		src = aux
	} else if(pdu.srcPort > 1000) {
		src = "0"+aux
	} else if(pdu.srcPort > 100) {
		src = "00"+aux
	} else if(pdu.srcPort > 10) {
		src = "000"+aux
	} else {
		src = "0000"+aux
	} 

	aux = strconv.FormatInt(int64(pdu.dstPort), 10)
	if(pdu.dstPort >10000) {
		dst = aux
	} else if(pdu.dstPort > 1000) {
		dst = "0"+aux
	} else if(pdu.dstPort > 100) {
		dst = "00"+aux
	} else if(pdu.dstPort > 10) {
		dst = "000"+aux
	} else {
		dst = "0000"+aux
	}

	aux = strconv.FormatInt(int64(pdu.size), 10)
	if(pdu.size > 10000) {
		size = aux
	} else if(pdu.size > 1000) {
		size = "0"+aux
	} else if(pdu.size > 100) {
		size = "00"+aux
	} else if(pdu.size > 10) {
		size = "000"+aux
	} else {
		size = "0000"+aux
	}

	aux = strconv.FormatInt(int64(pdu.checkSum), 10)
	if(pdu.checkSum > 10000) {
		check = aux
	} else if(pdu.checkSum > 1000) {
		check = "0"+aux
	} else if(pdu.checkSum > 100) {
		check = "00"+aux
	} else if(pdu.checkSum > 10) {
		check = "000"+aux
	} else {
		check = "0000"+aux
	}

	msgNetwork := src+dst+size+check+pdu.msg

	return msgNetwork
}

func log(pdu UdpPdu, msg string) {
	fmt.Println(time.Now())
	fmt.Print("\tPDU\t")
	fmt.Println(pdu)
	fmt.Println("\tMSG\t"+msg)
}


//****************************************************************

var applicationPathWR string
var applicationPathRD string
var networkPathWR string
var networkPathRD string

const app_tra = "./app_tra"
const tra_app = "./tra_app"
const tra_net = "./tra_net"
const net_tra = "./net_tra"

func main() {
	var pdu UdpPdu

	//****************************************************************
	//----------Open Pipe Network Read

	networkPathRD = "./network"
	
	// Create named pipe
	syscall.Mkfifo(net_tra, 0600)
	
	go func() {
		cmd := exec.Command(networkPathRD, net_tra)
		// Just to forward the stdout
	 	cmd.Stdout = os.Stdout
	 	cmd.Run()
	}()

	log(pdu, "Abrindo net_tra")
	stdoutNetworkRD, _ := os.OpenFile(net_tra, os.O_RDONLY, 0600)

	//****************************************************************
	//Ler Pipe Network
	log(pdu, "Lendo net_tra")
	var buffNetworkRD bytes.Buffer
	log(pdu, "Esperando alguem escrever algo")

	io.Copy(&buffNetworkRD, stdoutNetworkRD)

	msgNetwork := buffNetworkRD.String()
	stdoutNetworkRD.Close()
	log(pdu, "Mensagem recebida de Network Layer")
	log(pdu, msgNetwork)

	//****************************************************************
	//Tirar PDU
	pdu = NewUdpPdu(msgNetwork[0:5], msgNetwork[5:10], msgNetwork[10:15], msgNetwork[15:20], msgNetwork[20:])
	msgApp := pdu.msg

	//****************************************************************
	//Conferir Check Sum
	check := checkCheckSum(pdu)

	//****************************************************************
	//----------Open Pipe App Write

	applicationPathWR = "./application"
		
	// Create named pipe
	syscall.Mkfifo(tra_app, 0600)
	
	go func() {
		cmd := exec.Command(applicationPathWR, tra_app)
		// Just to forward the stdout
		cmd.Stdout = os.Stdout
		cmd.Run()
	}()

	log(pdu, "Abrindo tra_app")
	stdoutAppWR, _ := os.OpenFile(tra_app, os.O_WRONLY, 0600)

	//****************************************************************
	//Escrever Pipe App	
	if(check) {
		log(pdu, "Mensagem recebida corretamente")
	} else {
		log(pdu, "Mensagem recebida com erros")
	}
	stdoutAppWR.Write([]byte(msgApp))
	stdoutAppWR.Close()
	log(pdu, "Mensagem enviada para Application Layer")
	log(pdu, msgApp)

	//****************************************************************
	//----------Open Pipe App Read

	// Create named pipe
	syscall.Mkfifo(app_tra, 0600)
	
	log(pdu, "Abrindo app_tra")
	stdoutAppRD, _ := os.OpenFile(app_tra, os.O_RDONLY, 0600)

	//****************************************************************
	//Ler Pipe App
	log(pdu, "Lendo app_tra")
	var buffAppRD bytes.Buffer
	log(pdu, "Esperando alguem escrever algo")

	io.Copy(&buffAppRD, stdoutAppRD)
	stdoutAppRD.Close()

	msgApp = buffAppRD.String()
	log(pdu, "Mensagem recebida de Application Layer")
	log(pdu, msgApp)

	//****************************************************************
	//Montar PDU
	
	pdu.srcPort = 6500	
	pdu.dstPort = SmtpPort

	pdu.size = uint16(len(msgApp))		//len() retorna o numero de bytes da string
	pdu.msg = msgApp

	pdu.checkSum = calcCheckSum(pdu)

	log(pdu, "PDU montada")

	//****************************************************************
	//----------Open Pipe Network Write

	// Create named pipe
	syscall.Mkfifo(tra_net, 0600)

	log(pdu, "Abrindo tra_net")
	stdoutNetworkWR, _ := os.OpenFile(tra_net, os.O_WRONLY, 0600)

	//****************************************************************
	//Escrever Pipe Network
	msgNetwork = makeMsgNetwork(pdu)
	stdoutNetworkWR.Write([]byte(msgNetwork))
	stdoutNetworkWR.Close()
	log(pdu, "Mensagem enviada para Network Layer")
	log(pdu, msgNetwork)

}


