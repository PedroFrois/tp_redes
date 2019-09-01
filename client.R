log <- function(description) {
  time <- strftime(Sys.time(), "%D-%H:%M:%OS4", tz = '')
  cat(time, ' ', description, '\n')
}

logPdu <- function(layer, pdu) {
  time <- strftime(Sys.time(), "%D-%H:%M:%OS4", tz = '')
  cat(time, ' ', 'Layer: ', layer,'\nPdu: ', pdu, '\n')
}

getFile <- function(file_path){
  if(missing(file_path)){
    file_path <- 'file01.txt'
  }
  file <- file(description = file_path, open = 'r', blocking = TRUE)
  data <- readLines(file, 1)
  close(file)
  return (data)
}

convertPayloadToBin <- function(payload_size, payload) {
  chars <- strsplit(payload, split= "")[[1]]
  chars_bin <- ''
  for (i in 1:payload_size) {
     char_int <- strtoi(charToRaw(chars[i]),16L)
     char_bin <- R.utils::intToBin(char_int)
     chars_bin <- paste(chars_bin, char_bin, sep='')
  }
  return(chars_bin)
}

convertMacToBin <- function(macInHex){
  macInBin <- ''
  for (i in 0:5) {
    auxHex <- substr(macInHex, start = 1+3*i, stop = 2+3*i)
    auxDec <- strtoi(auxHex, 16L)
    auxBin <- R.utils::intToBin(auxDec)
    while(nchar(auxBin) < 8){
      auxBin <- paste('0', auxBin, sep = '')
    }
    macInBin <- paste(macInBin,auxBin, sep = '')
  }
  return(macInBin)
}

getSourceMacAddress <- function(ip){
  log('Getting source Mac Address') 
  command <- paste('ifconfig | grep -A 4 ', ip, ' | grep ether | tr -s [:blank:] | cut -d" " -f3', sep = '')
  mac <- system(command, intern = TRUE)
  log(paste("Mac Address:", mac))
  return(convertMacToBin(mac))
}

getDestinationMacAddress <- function(ip){
  log('Getting destination Mac Address')
  command <- paste('arp ',ip,'| grep ether | tr -s [:blank:] | cut -d" " -f3', sep = '')
  mac <- system(command, intern = TRUE)
  log(paste("Mac Address:", mac))
  return(convertMacToBin(mac))
}

modifyPdu <- function(mac_destination, mac_source, payload){
  log("Modifying pdu")
  payload_size <- R.utils::intToBin(nchar(payload))
  while(nchar(payload_size) < 2*8){
    payload <- paste('0', payload_size, sep = '')
  }
  modified_pdu <- paste(mac_destination, mac_source, payload_size, payload, sep = '')
  logPdu('Physical', modified_pdu)
  return(modified_pdu)
}

hexToBin <- function(string_array_hex){
  return(R.utils::intToBin(strtoi(string_array_hex, base = 16L)))
}

testColision <- function(){
  if(sample(1:10,1) == 1){
    return(TRUE)
  } else {
    return(FALSE)
  }
}

getIpFromPackage <- function(package, destination){
  if(destination){
    ip_description <- 'destination ip'
    byte_blocks <- c(0:3)
  } else{
    ip_description <- 'source ip'
    byte_blocks <- c(4:7)
  }
  log(paste('Getting', ip_description,'from package'))
  ip <- ''
  for (i in byte_blocks) {
    auxBin <- substr(package, start = 1+8*i, stop = 8+8*i)  
    auxDec <- strtoi(auxBin, 2L)
    ip <- paste(ip,auxDec,'.', sep = '')
  }
  ip <- substr(ip, start = 1, stop = (nchar(ip)-1)) #remove last dot
  log(paste('Ip:',ip))
  return(ip)
}

sendToServer <- function(ip_destination, port, file){
  log('Opening Connection...')
  
  while(testColision()){
    log('Colision detected!')
    Sys.sleep(sample(1:3,1))
  }
  
  connection <- socketConnection(host = ip_destination, port = port, blocking = TRUE,
    server = FALSE, open = 'r+')
  log('Connection openned')
  writeLines(file, connection)
  close(connection)
  log('Connection closed')
}

physical <- function(package){
  log('Creating frame from physical layer')
  ip_destination <- getIpFromPackage(package, TRUE)
  ip_source <- getIpFromPackage(package, FALSE)
  mac_destination <- getDestinationMacAddress(ip_destination)
  mac_source <- getSourceMacAddress(ip_source)
  
  frame <- modifyPdu(mac_destination, mac_source, package)

  port <- '' #ADD PORT ------------------------------------------------------------------

  sendToServer(ip_destination, port, frame)
}

network <- function(){
  log('Getting package from network layer')
  #file must contain the ip (in binary, only numbers and in bytes) followed by the payload
  package <- getFile()
  logPdu('Network', package)
  return(package)
}

client <- function(){
  package <- network()
  physical(package)
}

client()