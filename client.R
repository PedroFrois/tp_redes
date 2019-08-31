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

getFromMacAddress <- function(ip){ 
  command <- paste('ifconfig | grep -A 4 ', ip, ' | grep ether | tr -s [:blank:] | cut -d" " -f3', sep = '')
  mac <- system(command, intern = TRUE)
  cat("FROM MAC ADDRESS: ", mac , "\n")
  return(convertMacToBin(mac))
}

getToMacAddress <- function(ip){
  command <- paste('arp ',ip,'| grep ether | tr -s [:blank:] | cut -d" " -f3', sep = '')
  mac <- system(command, intern = TRUE)
  cat("TO MAC ADDRESS: ", mac , "\n")
  return(convertMacToBin(mac))
}

modifyPdu <- function(mac_to, mac_from, payload){
  payload_size <- R.utils::intToBin(nchar(payload))
  while(nchar(payload_size) < 2*8){
    payload <- paste('0', payload_size, sep = '')
  }
  return(paste(mac_to, mac_from, payload_size, payload, sep = ''))
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

getIpFromPackage <- function(package){
  ip <- ''
  for (i in 0:3) {
    auxBin <- substr(package, start = 1+8*i, stop = 8+8*i)  
    auxDec <- strtoi(auxBin, 2L)
    ip <- paste(ip,auxDec,'.', sep = '')
  }
  ip <- substr(ip, start = 1, stop = (nchar(ip)-1)) #remove last dot
  return(ip)
}

sendToServer <- function(ip_to, port, file){
  cat('File to send: ', file, '\n')

  writeLines('Opening Connection...')
  
  while(testColision()){
    writeLines('Colision detected!')
    Sys.sleep(sample(1:3,1))
  }
  
  connection <- socketConnection(host = ip_to, port = port, blocking = TRUE,
    server = FALSE, open = 'r+')
  writeLines(file, connection)
  close(connection)
}

physical <- function(package){
  ip_to <- getIpFromPackage(package)
  ip_from <- '' #ADD IP -----------------------------------------------------------------
  mac_to <- getToMacAddress(ip_to)
  mac_from <- getFromMacAddress(ip_from)
  
  frame <- modifyPdu(mac_to, mac_from, package)

  port <- '' #ADD PORT ------------------------------------------------------------------

  sendToServer(ip_to, port, frame)
}

network <- function(){
  #file must contain the ip (in binary, only numbers and in bytes) followed by the payload
  package <- getFile()
  return(package)
}

client <- function(){
  package <- network()
  physical(package)
}

client()