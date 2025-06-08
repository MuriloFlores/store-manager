# --- Estágio 1: Build ---
# Usamos uma imagem oficial do Go como nosso "builder".
# Usar uma versão específica (ex: 1.22) é melhor que "latest" para builds consistentes.
# "alpine" é uma versão da imagem baseada em um Linux mínimo, o que acelera o download.
FROM golang:1.24-alpine AS builder

# Define o diretório de trabalho dentro do container
WORKDIR /app

# Copia os arquivos de gerenciamento de dependências primeiro.
# Isso aproveita o cache do Docker. Se esses arquivos não mudarem,
# o Docker não vai baixar as dependências novamente em builds futuros.
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o resto do código-fonte da aplicação para o container.
COPY . .

# Comando de build. Aqui está a mágica:
# - CGO_ENABLED=0: Compila um binário estático, sem dependências de bibliotecas C do sistema.
#   Isso é CRUCIAL para que o binário rode em uma imagem mínima como a do Alpine ou Scratch.
# - -ldflags="-s -w": Remove símbolos de debug, resultando em um binário final muito menor.
# - -o /app/main: Especifica o nome e local do arquivo de saída (o binário compilado).
# - ./cmd/server/main.go: O ponto de entrada da sua aplicação que será compilado.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/main ./cmd/order_manager/main.go


# --- Estágio 2: Final ---
# Começamos com uma imagem do Alpine, que é extremamente leve (cerca de 5MB).
# Ela contém um sistema operacional mínimo, o que é ótimo para segurança e tamanho.
FROM alpine:latest

# Define o diretório de trabalho
WORKDIR /app

# [BOA PRÁTICA DE SEGURANÇA]
# Cria um grupo e um usuário específicos para a aplicação, para não rodar como root.
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copia APENAS o binário compilado do estágio "builder".
# Esta é a principal vantagem do multi-stage build. Nenhum código-fonte ou ferramenta
# de build vem junto, apenas o executável final.
COPY --from=builder /app/main /app/main

# Dá permissão de execução ao binário.
RUN chmod +x /app/main

# Garante que o usuário "appuser" seja o dono dos arquivos da aplicação.
RUN chown -R appuser:appgroup /app

# Muda para o usuário não-root que criamos.
USER appuser

# Expõe a porta em que a aplicação vai rodar. Isso é para documentação e para
# que o Docker saiba qual porta o container pretende expor.
EXPOSE 8080

# Comando para iniciar a aplicação quando o container for executado.
ENTRYPOINT ["/app/main"]