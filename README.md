# tp_redes
Implementação de pilha de protocolos para trabalho de Redes I

[Especifiação do Trabalho](https://docs.google.com/document/d/1O3cNM0T6gFNz9PeMYcnzbmBzEe8J7k34DaefJDSsv4A/edit)

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

*Escopo e decisões técnicas*:

