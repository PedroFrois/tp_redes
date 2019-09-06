## Physical Layer
*Preparação de ambiente*:
- Instalar R 
```
sudo apt install r-base
```
- Instalar pacote R-utils
```
R
install.packages("R-utils")
```

*Preparar execução*:
- Colocar os arquivos client.R e server.R em computadores diferentes
- No computador com o client.R escrever um arquivo "file01.txt" contendo a pdu recebida da camada de rede
  - ip do destino (em 4 bytes)
  - ip da origem (em 4 bytes)
  - tamanho do payload (em 2 bytes)
  - payload
- Escrever a porta desejada para estabelecer a conexão tanto no arquivo "client.R" e "server.R" na variável "port"

*Execução*:
- Na máquina com o arquivo "server.R" executar o comando:
```Rscript server.R```
- Na máquina com o arquivo "client.R" executar o comando:
```Rscript client.R```

*Decisões técnicas*:
- A PDU recebida pela parte física:
  - deverá estar escrita em um arquivo chamado "file01.txt" na mesma pasta do código fonte "client.R"
  - deverá conter o ip de destino(4 bytes), ip de origem(4 bytes), tamanho do payload(2 bytes) e o payload, sendo que os ips e o tamanho do payload devem estar em bits. Ex: 10101010101010101010101010101010101010101010101010101010101010101010101010101010payload

*Teste*
- ![Servidor esperando cliente](/testImagesFolder/testServerListening.png)
- ![Cliente recebendo pacote da camada de rede; montando o quadro; abrindo a conexão com o servidor; teste de colisão; envio finalizado; conexão terminada](/testImagesFolder/testClient.png)
- ![Servidor recebendo quadro; retirando cabeçalho; e salvando em arquivo](/testImagesFolder/testServer.jpeg)
