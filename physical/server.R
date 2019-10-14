log <- function(description) {
  time <- strftime(Sys.time(), "%D-%H:%M:%OS4", tz = '')
  cat(time, ' ', description, '\n')
}

logPdu <- function(layer, pdu) {
  time <- strftime(Sys.time(), "%D-%H:%M:%OS4", tz = '')
  cat(time, ' ', 'Layer: ', layer,'\nPdu: ', pdu, '\n')
}

getPackageFromFrame <- function(frame){
  log('Removing header from frame')
  package <- substr(frame, start = 6*8*2+2*8, stop = nchar(frame)) #remove header
  log('Header removed')
  logPdu('Physical', package)
  return(package)
}

physical <- function(){
  log('Openning pipe')
  stream <- fifo(description = "phy_net", open = "w",)

  while(TRUE){
    log('Listenning...')
    connection <- socketConnection(port = '', blocking = TRUE, #ADD PORT -------- 
      server = TRUE, open = 'r+')
    log('Connection openned')
    log('Receiving frame')
    frame <- readLines(connection, 1)
    logPdu('Physical', frame)
    package <- getPackageFromFrame(frame)
    
    log(cat('Pipe is open: ', isOpen(stream)))
    log('Writing package to pipe')
    writeLines(package, stream)
  }
}

server <- function(){
  physical()
}

server()