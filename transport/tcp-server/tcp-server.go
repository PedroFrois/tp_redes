package main

import (
	"fmt"
	"strconv"
	"time"

	"io"
	"os"
	"os/exec"
	"syscall"
	"bytes" 
)

var sequence 	uint32
var ackowledge	uint32

// PDU - TCP
type TcpPdu struct {
	SrcPort		uint16
	DstPort		uint16
	SeqNumber	uint32
	AckNumber	uint32
	Window		uint16
	Urg			bool 		//Urgent
	Ack			bool 		//Acknowledgement
	Psh			bool 		//Push
	Rst			bool 		//Reset
	Syn			bool 		//Synchronization
	Fin			bool 		//Finish
	CheckSum	uint16
	Msg			string
}

const (
	SmtpPort = 587
	ImapPort = 993
	Pop3Port = 995
)


func NewTcpPdu (src, dst, seqN, ackN, window, u, a, p, r, s, f, check, msg string) TcpPdu {

	var c uint64
	var pdu TcpPdu

	c, _ = strconv.ParseUint(src, 10, 16)
	pdu.SrcPort = uint16(c)
	c, _ = strconv.ParseUint(dst, 10, 16)
	pdu.DstPort = uint16(c)

	c, _ = strconv.ParseUint(seqN, 10, 32)
	pdu.SeqNumber = uint32(c)
	c, _ = strconv.ParseUint(ackN, 10, 32)
	pdu.AckNumber = uint32(c)

	c, _ = strconv.ParseUint(window, 10, 16)
	pdu.Window = uint16(c)

	c, _ = strconv.ParseUint(u, 10, 1)	
	if(c == 1){ pdu.Urg = true} else {pdu.Urg = false}
	c, _ = strconv.ParseUint(a, 10, 1)	
	if(c == 1){ pdu.Ack = true} else {pdu.Ack = false}
	c, _ = strconv.ParseUint(p, 10, 1)	
	if(c == 1){ pdu.Psh = true} else {pdu.Psh = false}
	c, _ = strconv.ParseUint(r, 10, 1)	
	if(c == 1){ pdu.Rst = true} else {pdu.Rst = false}
	c, _ = strconv.ParseUint(s, 10, 1)	
	if(c == 1){ pdu.Syn = true} else {pdu.Syn = false}
	c, _ = strconv.ParseUint(f, 10, 1)	
	if(c == 1){ pdu.Fin = true} else {pdu.Fin = false}

	c, _ = strconv.ParseUint(check, 10, 16)
	pdu.CheckSum = uint16(c)

	pdu.Msg = msg

	return pdu
}

func BoolToInt(b bool) int8 {
	if(b == false) {
		return 0
	} else {
		return 1
	}
}

func CalcCheckSum (pdu TcpPdu) uint16 {
	//http://www.jvasconcellos.com.br/fat/FAT_TI/wp-content/uploads/2013/10/checksum-udp.pdf
	//1. determinação do comprimento do datagrama;
	//pdu.size ja fornece o tamanho dos dados em byte

	//2. agrupamento dos campos em blocos de 16 bits;
	//16 bits = 2 bytes = 2 caracteres
	//Todos os campos com excessao da msg ja sao inteiros binarios de 16 bits

	//3. conversão da string menssagem para números binários de 16 bits;
    blocks := []byte(pdu.Msg)
    var i, j uint16
    j = 0

    var m [32767]uint16

    for i = 0; i < uint16(len(pdu.Msg)-1); i=i+2 {
    	m[j] = uint16(((uint16(blocks[i+1]))*256) | uint16(blocks[i]))
    	j++
    }

    //Conversão da Window e das Flags em um numero de 16 bits
    aux := (pdu.Window << 6) + uint16((BoolToInt(pdu.Urg) << 5) + (BoolToInt(pdu.Ack) << 4) + (BoolToInt(pdu.Psh) << 3) + (BoolToInt(pdu.Rst) << 2) + (BoolToInt(pdu.Syn) << 1) + (BoolToInt(pdu.Fin)))


	//4. soma dos blocos de 16 bits;
	var sum uint32
	sum = uint32(pdu.SrcPort + pdu.DstPort + uint16(pdu.SeqNumber >> 16) + uint16(pdu.SeqNumber) + uint16(pdu.AckNumber >> 16) + uint16(pdu.AckNumber) + aux + pdu.CheckSum)
	
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

	return resp
}

func CheckCheckSum(pdu TcpPdu) bool {
	sum := CalcCheckSum(pdu)

	if sum == pdu.CheckSum {
		return true
	} else {
		return false
	}
}

func MakeMsgNetwork(pdu TcpPdu) string {
	var aux, src, dst, seqN, ackN, window, check string
	var u, a, p, r, s, f string

	aux = strconv.FormatInt(int64(pdu.SrcPort), 10)
	if(pdu.SrcPort > 10000) {
		src = aux
	} else if(pdu.SrcPort > 1000) {
		src = "0"+aux
	} else if(pdu.SrcPort > 100) {
		src = "00"+aux
	} else if(pdu.SrcPort > 10) {
		src = "000"+aux
	} else {
		src = "0000"+aux
	} 

	aux = strconv.FormatInt(int64(pdu.DstPort), 10)
	if(pdu.DstPort > 10000) {
		dst = aux
	} else if(pdu.DstPort > 1000) {
		dst = "0"+aux
	} else if(pdu.DstPort > 100) {
		dst = "00"+aux
	} else if(pdu.DstPort > 10) {
		dst = "000"+aux
	} else {
		dst = "0000"+aux
	}

	aux = strconv.FormatInt(int64(pdu.SeqNumber), 10)
	if(pdu.SeqNumber > 1000000000) {
		seqN = aux
	} else if(pdu.SeqNumber > 100000000) {
		seqN = "0"+aux
	} else if(pdu.SeqNumber > 10000000) {
		seqN = "00"+aux
	} else if(pdu.SeqNumber > 1000000) {
		seqN = "000"+aux
	} else if(pdu.SeqNumber > 100000) {
		seqN = "0000"+aux
	} else if(pdu.SeqNumber > 10000) {
		seqN = "00000"+aux
	} else if(pdu.SeqNumber > 1000) {
		seqN = "000000"+aux
	} else if(pdu.SeqNumber > 100) {
		seqN = "0000000"+aux
	} else if(pdu.SeqNumber > 10) {
		seqN = "00000000"+aux
	} else {
		seqN = "000000000"+aux
	}

	aux = strconv.FormatInt(int64(pdu.AckNumber), 10)
	if(pdu.AckNumber > 1000000000) {
		ackN = aux
	} else if(pdu.AckNumber > 100000000) {
		ackN = "0"+aux
	} else if(pdu.AckNumber > 10000000) {
		ackN = "00"+aux
	} else if(pdu.AckNumber > 1000000) {
		ackN = "000"+aux
	} else if(pdu.AckNumber > 100000) {
		ackN = "0000"+aux
	} else if(pdu.AckNumber > 10000) {
		ackN = "00000"+aux
	} else if(pdu.AckNumber > 1000) {
		ackN = "000000"+aux
	} else if(pdu.AckNumber > 100) {
		ackN = "0000000"+aux
	} else if(pdu.AckNumber > 10) {
		ackN = "00000000"+aux
	} else {
		ackN = "000000000"+aux
	}

	aux = strconv.FormatInt(int64(pdu.Window), 10)
	if(pdu.Window > 1000) {
		window = aux
	} else if(pdu.Window > 100) {
		window = "0"+aux
	} else if(pdu.Window > 10) {
		window = "00"+aux
	} else {
		window = "000"+aux
	}

	//Flags
	if pdu.Urg {
		u = "1"
	} else {
		u = "0"
	}
	if pdu.Ack {
		a = "1"
	} else {
		a = "0"
	}
	if pdu.Psh {
		p = "1"
	} else {
		p = "0"
	}
	if pdu.Rst {
		r = "1"
	} else {
		r = "0"
	}
	if pdu.Syn {
		s = "1"
	} else {
		s = "0"
	}
	if pdu.Fin {
		f = "1"
	} else {
		f = "0"
	}


	aux = strconv.FormatInt(int64(pdu.CheckSum), 10)
	if(pdu.CheckSum > 1000) {
		check = aux
	} else if(pdu.CheckSum > 100) {
		check = "0"+aux
	} else if(pdu.CheckSum > 10) {
		check = "00"+aux
	} else {
		check = "000"+aux
	}

	msgNetwork := src+dst+seqN+ackN+window+u+a+p+r+s+f+check+pdu.Msg

	return msgNetwork
}


func log(pdu TcpPdu, msg string) {
	fmt.Println(time.Now())
	fmt.Print("\tPDU\t")
	fmt.Println(pdu)
	fmt.Println("\tMSG\t"+msg)
}

var applicationPathWR string
var applicationPathRD string
var networkPathWR string
var networkPathRD string

const app_tra = "./app_tra"
const tra_app = "./tra_app"
const tra_net = "./tra_net"
const net_tra = "./net_tra"

func main() { 
	var pdu TcpPdu
	sequence = 1
	ackowledge = 1

	//****************************************************************
	//----------Open Pipe Network Read

	networkPathRD = "./network"
	
	// Create named pipe
	syscall.Mkfifo(net_tra, 0600)
	
	func() {
		cmd1 := exec.Command(networkPathWR, tra_net)
		// Just to forward the stdout
		cmd1.Stdout = os.Stdout
		cmd1.Run()
	}()

	log(pdu, "Abrindo net_tra")
	stdoutNetworkRD, _ := os.OpenFile(net_tra, os.O_RDONLY, 0600)
	
	//****************************************************************
	//Ler Pipe Network
//------------------------ PARTE 1 DO THREE WAY HANDSHAKE ------------------------
	log(pdu, "Lendo networkPipe")
	var buffNetworkRD bytes.Buffer
	log(pdu, "Esperando alguem escrever algo")

	io.Copy(&buffNetworkRD, stdoutNetworkRD)

	msgNetwork := buffNetworkRD.String()
	log(pdu, "Mensagem recebida de Network Layer")
	log(pdu, msgNetwork)

	//****************************************************************
	//Tirar PDU
				//SrcPort - 16 bits,	DstPort - 16, 	SeqNumber - 32,		AckNumber - 32,	   	Window - 10,		Flags - 6,																												CheckSum - 16, 			Msg	
	pdu = NewTcpPdu(msgNetwork[0:5], msgNetwork[5:10], msgNetwork[10:20], msgNetwork[20:30], msgNetwork[30:34], msgNetwork[34:35], msgNetwork[35:36], msgNetwork[36:37], msgNetwork[37:38], msgNetwork[38:39], msgNetwork[39:40], msgNetwork[40:45], msgNetwork[45:])
	msgApp := msgNetwork[16:]

	sequence = pdu.SeqNumber + uint32(len(msgApp))
	ackowledge = sequence

	//****************************************************************
	//Conferir Check Sum
	check := CheckCheckSum(pdu)

	if(check) {
		log(pdu, "Mensagem recebida com sucesso")
	} else {
		log(pdu, "Mensagem recebida com erros")
		pdu.Msg = "Mensagem recebida com erros"
	}

	//****************************************************************
	//conferir se foi recebido SYN
	if((pdu.Syn == true)) {
		log(pdu, "Mensagem recebida pelo Servidor com sucesso \n\t\tParte 1 do Three-Way Handshake")
	} else {
		log(pdu, "Mensagem não foi recebida pelo Servidor")
		pdu.Msg = "Mensagem não foi recebida pelo Servidor"
	}

	//****************************************************************
	//Open Pipe App Write
	
	applicationPathRD = "./application"

	func() {
		cmd2 := exec.Command(applicationPathRD, app_tra)
		// Just to forward the stdout
		cmd2.Stdout = os.Stdout
		cmd2.Run()
	}()

	// Create named pipe
	syscall.Mkfifo(tra_app, 0600)
	
	log(pdu, "Abrindo tra_app")
	stdoutAppWR, _ := os.OpenFile(tra_app, os.O_WRONLY, 0600)

	//****************************************************************
	//Escrever Pipe App
	stdoutAppWR.Write([]byte(msgApp))
	log(pdu, "Mensagem enviada para Application Layer")	
	log(pdu, msgApp)

	//****************************************************************
	//Open Pipe App Read
	
	// Create named pipe
	syscall.Mkfifo(app_tra, 0600)

	log(pdu, "Abrindo app_tra")
	stdoutAppRD, _ := os.OpenFile(app_tra, os.O_RDONLY, 0600)

	//****************************************************************
	//Ler Pipe App
	log(pdu, "Lendo appPipe")
	var buffAppRD bytes.Buffer
	log(pdu, "Esperando alguem escrever algo")
	
	io.Copy(&buffAppRD, stdoutAppRD)

	msgApp  = buffAppRD.String()
	log(pdu, "Mensagem recebida de Application Layer")
	log(pdu, msgApp)

	//****************************************************************
	//Montar PDU de Sincornização
//------------------------ PARTE 2 DO THREE WAY HANDSHAKE ------------------------
	
	pdu.SrcPort = 6000	
	pdu.DstPort = SmtpPort

	pdu.SeqNumber = sequence
	pdu.AckNumber = ackowledge
	pdu.Window = 5

	//Flags
	pdu.Urg	= false
	pdu.Ack = true
	pdu.Psh = true
	pdu.Rst = false
	pdu.Syn = true
	pdu.Fin	= false

	pdu.Msg = msgApp

	pdu.CheckSum = CalcCheckSum(pdu)

	log(pdu, "PDU montada")

	//****************************************************************
	//Open Pipe Network Write

	// Create named pipe
	syscall.Mkfifo(tra_net, 0600)
	

	log(pdu, "Abrindo tra_net")
	stdoutNetworkWR, _ := os.OpenFile(tra_net, os.O_WRONLY, 0600)

	//****************************************************************
	//Escrever Pipe Network
	msgNetwork = MakeMsgNetwork(pdu)
	stdoutNetworkWR.Write([]byte(msgNetwork))
	log(pdu, "Mensagem enviada para Network Layer \n\t\tParte 2 do Three-Way Handshake")
	log(pdu, msgNetwork)

	//****************************************************************
	//Ler Pipe Network
	log(pdu, "Lendo networkPipe")
	//var buff bytes.Buffer
	log(pdu, "Esperando alguem escrever algo")

	io.Copy(&buffNetworkRD, stdoutNetworkRD)

	msgNetwork = buffNetworkRD.String()
	log(pdu, "Mensagem recebida de Network Layer")
	log(pdu, msgNetwork)

	//****************************************************************
	//Tirar PDU
	pdu = NewTcpPdu(msgNetwork[0:5], msgNetwork[5:10], msgNetwork[10:20], msgNetwork[20:30], msgNetwork[30:34], msgNetwork[34:35], msgNetwork[35:36], msgNetwork[36:37], msgNetwork[37:38], msgNetwork[38:39], msgNetwork[39:40], msgNetwork[40:45], msgNetwork[45:])
	msgApp = msgNetwork[16:]

	sequence = pdu.SeqNumber + uint32(len(msgApp))
	ackowledge = sequence

	//****************************************************************
	//Conferir Check Sum
	check = CheckCheckSum(pdu)

	if(check) {
		log(pdu, "Mensagem recebida com sucesso")
	} else {
		log(pdu, "Mensagem recebida com erros")
	}

	//****************************************************************
//------------------------ PARTE 2 DO THREE WAY HANDSHAKE ------------------------
	//conferir se foi recebido ACK
	if((pdu.Ack == true)) {
		log(pdu, "Mensagem recebida pelo Servidor com sucesso \n\t\tParte 3 do Three-Way Handshake")
	} else {
		log(pdu, "Mensagem não foi recebida pelo Servidor")
	}

	//****************************************************************
	//Escrever Pipe App
	stdoutAppWR.Write([]byte(msgApp))
	log(pdu, "Mensagem enviada para Application Layer")
	log(pdu, msgApp)

//------------------------ RECEBIMENTO DO EMAIL ------------------------
	//****************************************************************
	//Ler Pipe Network
	log(pdu, "Lendo networkPipe")
	//var buff bytes.Buffer
	log(pdu, "Esperando alguem escrever algo")

	io.Copy(&buffNetworkRD, stdoutNetworkRD)

	msgNetwork = buffNetworkRD.String()
	log(pdu, "Mensagem recebida de Network Layer")
	log(pdu, msgNetwork)

	//****************************************************************
	//Tirar PDU
	pdu = NewTcpPdu(msgNetwork[0:5], msgNetwork[5:10], msgNetwork[10:20], msgNetwork[20:30], msgNetwork[30:34], msgNetwork[34:35], msgNetwork[35:36], msgNetwork[36:37], msgNetwork[37:38], msgNetwork[38:39], msgNetwork[39:40], msgNetwork[40:45], msgNetwork[45:])
	msgApp = msgNetwork[16:]

	sequence = pdu.SeqNumber + uint32(len(msgApp))
	ackowledge = sequence

	//****************************************************************
	//Conferir Check Sum
	check = CheckCheckSum(pdu)

	if(check) {
		log(pdu, "Mensagem recebida com sucesso")
	} else {
		log(pdu, "Mensagem recebida com erros")
		pdu.Msg = "Mensagem recebida com erros"
	}

	//****************************************************************
	//Escrever Pipe App
	stdoutAppWR.Write([]byte(msgApp))
	log(pdu, "Mensagem enviada para Application Layer")
	log(pdu, msgApp)

//****************************************************************
	//Ler Pipe App
	log(pdu, "Lendo appPipe")
	// var buff bytes.Buffer
	log(pdu, "Esperando alguem escrever algo")
	
	io.Copy(&buffAppRD, stdoutAppRD)

	msgApp = buffAppRD.String()
	log(pdu, "Mensagem recebida de Application Layer")
	log(pdu, msgApp)

	//****************************************************************
	//Montar PDU de Sincornização	
	pdu.SrcPort = 6000	
	pdu.DstPort = SmtpPort

	pdu.SeqNumber = sequence
	pdu.AckNumber = ackowledge
	pdu.Window = 5

	//Flags
	pdu.Urg	= false
	pdu.Ack = false
	pdu.Psh = true
	pdu.Rst = false
	pdu.Syn = false
	pdu.Fin	= false

	pdu.Msg = msgApp

	pdu.CheckSum = CalcCheckSum(pdu)

	log(pdu, "PDU montada")

	//****************************************************************
	//Escrever Pipe Network
	msgNetwork = MakeMsgNetwork(pdu)
	stdoutNetworkWR.Write([]byte(msgNetwork))
	log(pdu, "Mensagem enviada para Network Layer")
	log(pdu, msgNetwork)

//------------------------ FINALIZAR COMUNICACAO ------------------------
	//****************************************************************
	//Ler Pipe Network
	//------------------------ PARTE 1 DO THREE WAY HANDSHAKE MODIFICADO ------------------------
	log(pdu, "Lendo networkPipe")
	//var buff bytes.Buffer
	log(pdu, "Esperando alguem escrever algo")

	io.Copy(&buffNetworkRD, stdoutNetworkRD)

	msgNetwork = buffNetworkRD.String()
	log(pdu, "Mensagem recebida de Network Layer")
	log(pdu, msgNetwork)

	//****************************************************************
	//Tirar PDU
				//SrcPort - 16 bits,	DstPort - 16, 	SeqNumber - 32,		AckNumber - 32,	   	Window - 10,		Flags - 6,																												CheckSum - 16, 			Msg	
	pdu = NewTcpPdu(msgNetwork[0:5], msgNetwork[5:10], msgNetwork[10:20], msgNetwork[20:30], msgNetwork[30:34], msgNetwork[34:35], msgNetwork[35:36], msgNetwork[36:37], msgNetwork[37:38], msgNetwork[38:39], msgNetwork[39:40], msgNetwork[40:45], msgNetwork[45:])
	msgApp = msgNetwork[16:]

	sequence = pdu.SeqNumber + uint32(len(msgApp))
	ackowledge = sequence

	//****************************************************************
	//Conferir Check Sum
	check = CheckCheckSum(pdu)

	if(check) {
		log(pdu, "Mensagem recebida com sucesso")
	} else {
		log(pdu, "Mensagem recebida com erros")
		pdu.Msg = "Mensagem recebida com erros"
	}

	//****************************************************************
	//conferir se foi recebido FIN
	if((pdu.Fin == true)) {
		log(pdu, "Mensagem recebida pelo Servidor com sucesso \n\t\tParte 2 do Three-Way Handshake")
	} else {
		log(pdu, "Mensagem não foi recebida pelo Servidor")
		pdu.Msg = "Mensagem não foi recebida pelo Servidor"
	}

	//****************************************************************
	//Escrever Pipe App
	stdoutAppWR.Write([]byte(msgApp))
	log(pdu, "Mensagem enviada para Application Layer")
	log(pdu, msgApp)

	//****************************************************************
	//Ler Pipe App
	log(pdu, "Lendo appPipe")
	// var buff bytes.Buffer
	log(pdu, "Esperando alguem escrever algo")
	
	io.Copy(&buffAppRD, stdoutAppRD)

	msgApp = buffAppRD.String()
	log(pdu, "Mensagem recebida de Application Layer")
	log(pdu, msgApp)

	//****************************************************************
	//Montar PDU de Sincornização
//------------------------ PARTE 2 DO THREE WAY HANDSHAKE MODIFICADO ------------------------
	
	pdu.SrcPort = 6000	
	pdu.DstPort = SmtpPort

	pdu.SeqNumber = sequence
	pdu.AckNumber = ackowledge
	pdu.Window = 5

	//Flags
	pdu.Urg	= false
	pdu.Ack = true
	pdu.Psh = true
	pdu.Rst = false
	pdu.Syn = true
	pdu.Fin	= false

	pdu.Msg = msgApp

	pdu.CheckSum = CalcCheckSum(pdu)

	log(pdu, "PDU montada")


	//****************************************************************
	//Escrever Pipe Network
	msgNetwork = MakeMsgNetwork(pdu)
	stdoutNetworkWR.Write([]byte(msgNetwork))
	log(pdu, "Mensagem enviada para Network Layer \n\t\tParte 2 do Three-Way Handshake")
	log(pdu, msgNetwork)

	stdoutNetworkRD.Close()
	stdoutNetworkWR.Close()
	stdoutAppRD.Close()
	stdoutAppWR.Close()
}









