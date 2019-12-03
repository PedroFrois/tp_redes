import os, tempfile, sys, csv, fcntl

def calculateChecksum(data):
	return bin(365)

def generateDatagram(transportPdu, dest):
	version = "0100"
	ihl = "0101"
	tos = "00100000"
	totalLength = str(bin(5+(len(transportPdu)//8)))
	identification = "0000000000000000"
	flags = "010"
	offset = "0000000000000"
	ttl = "00001000"
	protocol =  "00000010"

	datagram = version + ihl + tos
	aux = 0
	auxStr = ""
	while ((len(totalLength) + aux) < 16):
		auxStr += "0"
		aux += 1
	datagram += auxStr + totalLength[2::] + identification + flags + offset + ttl + protocol

	checksum = calculateChecksum(datagram)
	aux = 0
	auxStr = ""
	while ((len(checksum) + aux) < 16):
		auxStr += "0"
		aux += 1
	datagram += auxStr + checksum[2:]

	source = "???"
	octects = dest.split('.')
	destination = ''
	for o in octects:
		destination += bin(int(o))[2:]

	datagram += source + destination

	#print(datagram)
	return datagram

def checkRoutingTable(destination):
	with open('routing_table.csv') as csvfile:
		readCSV = csv.reader(csvfile, delimiter=';')
		header = True
		for row in readCSV:
			if header:
				header = False
				continue
			row = row[0].split(', ')
			network = row[0]
			netmask = row[1]
			gateway = row[2]
			if (network == "default"):
				return gateway
			count, ok = 0, 0
			networkF,  netmaskF, destinationF = -1, 0, 0
			while (count < 4):
				networkF = network.find('.', networkF+1)
				netmaskF = netmask.find('.', netmaskF+1)
				destinationF = destination.find('.', destinationF+1)
				if(networkF == -1):
					networkF = len(network)
				if(netmaskF == -1):
					netmaskF = len(netmask)
				if(destinationF == -1):
					destinationF = len(destination)
				networkBits = network[:networkF-1]
				netmaskBits = netmask[:netmaskF-1]
				destinationBits = gateway[:destinationF-1]
				if(networkBits and netmaskBits == destinationBits and netmaskBits):
					ok+=1
				count+=1
			if (ok == 4):
				return gateway


tra_net_path = "tra_net"
net_tra_path = "net_tra"
phy_net_path = "phy_net"
net_phy_path = "net_phy"
f_tra_net = os.open(tra_net_path, os.O_RDWR)
fl = fcntl.fcntl(f_tra_net, fcntl.F_GETFL)
fcntl.fcntl(f_tra_net, fcntl.F_SETFL, fl | os.O_NONBLOCK)
tra_net = os.fdopen(f_tra_net, 'r')
print("abri tra_net")
f_net_tra = os.open(net_tra_path, os.O_RDWR)
net_tra = os.fdopen(f_net_tra, 'w')
print("abri net_tra")
f_phy_net = os.open(phy_net_path, os.O_RDWR)
fl = fcntl.fcntl(f_phy_net, fcntl.F_GETFL)
fcntl.fcntl(f_phy_net, fcntl.F_SETFL, fl | os.O_NONBLOCK)
phy_net = os.fdopen(f_phy_net, 'r')
print("abri phy_net")
f_net_phy = os.open(net_phy_path, os.O_RDWR)
net_phy = os.fdopen(f_net_phy, 'w')
print("abri net_phy")
ipAddress = "192.16.84.17"


while True:
	line = tra_net.readline()
	if line:
		#receber pdu trans --ok
		#olhar tabela?
		rout = checkRoutingTable(ipAddress)
		#preencher header do datagrama
		datagram = generateDatagram(line, rout)
		#enviar datagrama
		net_phy.writelines(datagram)

	line = phy_net.readline()
	if line:
		#receber pdu fisica
		message = line[32:]
		net_tra.writelines(message)

tra_net.close()
net_tra.close()
phy_net.close()
net_phy.close()
