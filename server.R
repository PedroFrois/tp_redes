getPackageFromFrame <- function(frame){
  package <- substr(frame, start = 6*8*2+2*8, stop = nchar(frame)) #remove header
  return(package)
}

physical <- function(){
  while(TRUE){
    writeLines('Listenning...')
    connection <- socketConnection(port = '', blocking = TRUE, #ADD PORT -------- 
      server = TRUE, open = 'r+')
    frame <- readLines(connection, 1)
    package <- getPackageFromFrame(frame)
    writeLines(package, 'file02.txt')
  }
}

server <- function(){
  physical()
}

server()