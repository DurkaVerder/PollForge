#!/usr/bin/env bash
set -euo pipefail

echo "Остановка контейнеров и удаление связанных томов..."
docker-compose down -v

echo "Prune неиспользуемых томов..."
docker volume prune -f

ALL_VOLUMES=$(docker volume ls -q)
if [ -n "$ALL_VOLUMES" ]; then
  echo "Удаляем все тома в системе..."
  docker volume rm $ALL_VOLUMES
else
  echo "Нет томов для удаления."
fi

echo "Готово."
