# Comece a partir da imagem de base golang
FROM golang:alpine as builder

# ENV GO111MODULE=on

# Adiciona o manutedor do template 
LABEL maintainer="Joao Soares <joaosoaresa.alm@gmail.com>"

# Instale o git.
# Git é necessário para buscar as dependências.
RUN apk update && apk add --no-cache git

# Defina o diretório de trabalho atual dentro do contêiner
WORKDIR /app

# Copiar arquivos go mod e sum
COPY go.mod go.sum ./



# Baixe todas as dependências. As dependências serão armazenadas em cache se os arquivos go.mod e go.sum não forem alterados
RUN go mod download 

# Copie a fonte do diretório atual para o diretório de trabalho dentro do contêiner

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copie o arquivo binário pré-construído do estágio anterior
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Exponha a porta 8080 
EXPOSE 8080

#Comando para rodar o executável
CMD ["./main"]
