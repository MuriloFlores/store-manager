#!/bin/bash

# Este script serve como um atalho para executar comandos de migração
# do banco de dados usando o Docker Compose.

# Carrega as variáveis de ambiente do arquivo .env para que o Docker Compose as utilize.
set -a
source .env
set +a

# Define a base do comando para não precisar repetir.
# Nota: Adicionamos a flag --rm para que o container seja removido após a execução,
# mantendo o ambiente limpo.
DOCKER_CMD="docker-compose run --rm migrate -path /migrations -database postgres://${DB_USER}:${DB_PASSWORD}@postgres:5432/${DB_NAME}?sslmode=disable"

# O primeiro argumento para este script (ex: 'up', 'down', 'create').
COMMAND=$1

# Verifica se um comando foi fornecido.
if [ -z "$COMMAND" ]; then
    echo "Uso: ./migrate.sh [up|down|down-all|force|version|create]"
    echo "  - up: Aplica todas as migrações pendentes."
    echo "  - down [N]: Reverte N migrações. Ex: ./migrate.sh down 1"
    echo "  - down-all: ⚠️ CUIDADO! Reverte TODAS as migrações."
    echo "  - force [V]: Força o banco para a versão V. Ex: ./migrate.sh force 3"
    echo "  - version: Mostra a versão atual da migração."
    echo "  - create [nome]: Cria um novo arquivo de migração. Ex: ./migrate.sh create nome_da_migracao"
    exit 1
fi

# Executa o comando apropriado.
case "$COMMAND" in
    create)
        # O segundo argumento é o nome da migração.
        NAME=$2
        if [ -z "$NAME" ]; then
            echo "Erro: Forneça um nome para a migração."
            echo "Uso: ./migrate.sh create [nome_da_migracao]"
            exit 1
        fi
        # Executa o comando create DENTRO do container, que tem a CLI.
        # Mas o arquivo será criado na pasta local por causa do volume mapeado!
        docker-compose run --rm migrate create -ext sql -dir /migrations -seq "$NAME"
        echo "Arquivos de migração para '$NAME' criados em ./infrastructure/db/migrations/"
        ;;
    # --- Novo caso adicionado ---
    down-all)
        echo "⚠️  ATENÇÃO: Revertendo TODAS as migrações. Isso irá apagar o esquema do banco de dados..."
        $DOCKER_CMD down -all
        ;;
    up|down|force|version)
        # Passa todos os argumentos restantes para o comando principal.
        # Isso permite fazer './migrate.sh down 1' ou './migrate.sh force 3'.
        $DOCKER_CMD "$@"
        ;;
    *)
        echo "Comando desconhecido: $COMMAND"
        exit 1
        ;;
esac

echo "Comando '$*' executado com sucesso."