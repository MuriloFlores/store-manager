#!/bin/bash

# Este script serve como um atalho para executar comandos de migração
# do banco de dados usando o Docker Compose.

# Carrega as variáveis de ambiente do arquivo .env para que o Docker Compose as utilize.
set -a
source .env
set +a

# Define a base do comando para não precisar repetir.
DOCKER_CMD="docker-compose run --rm migrate -path /migrations -database postgres://${DB_USER}:${DB_PASSWORD}@postgres:5432/${DB_NAME}?sslmode=disable"

# O primeiro argumento para este script (ex: 'up', 'down', 'create').
COMMAND=$1

# Verifica se um comando foi fornecido.
if [ -z "$COMMAND" ]; then
    echo "Uso: ./migrate.sh [up|down|down-all|force|version|create]"
    exit 1
fi

# Executa o comando apropriado.
case "$COMMAND" in
    create)
        NAME=$2
        if [ -z "$NAME" ]; then
            echo "Erro: Forneça um nome para a migração."
            echo "Uso: ./migrate.sh create [nome_da_migracao]"
            exit 1
        fi

        # Executa o container com o UID/GID do usuário host para evitar problemas de permissão.
        docker-compose run --rm --user "$(id -u):$(id -g)" migrate create -ext sql -dir /migrations -seq "$NAME"

        echo "Arquivos de migração para '$NAME' criados em ./infrastructure/db/migrations/"
        ;;
    down-all)
        echo "⚠️  ATENÇÃO: Revertendo TODAS as migrações. Isso irá apagar o esquema do banco de dados..."
        $DOCKER_CMD down -all
        ;;
    up|down|force|version)
        $DOCKER_CMD "$@"
        ;;
    *)
        echo "Comando desconhecido: $COMMAND"
        exit 1
        ;;
esac

echo "Comando '$*' executado com sucesso."