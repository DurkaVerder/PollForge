#!/bin/bash

set -euo pipefail

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${NC}^-^ Запуск go mod tidy для всех модулей...${NC}"

MODULE_DIRS=$(find . -mindepth 2 -name "go.mod" -exec dirname {} \;)

any_failed=0

for dir in $MODULE_DIRS; do
  echo -e "${NC}Модуль: $dir${NC}"
  pushd "$dir" > /dev/null


  OUTPUT=$(go mod tidy -v 2>&1)

  # Если есть сообщения о скачивании зависимостей, то фейлим
  if echo "$OUTPUT" | grep -qE "go: downloading|go: finding module for package";  then
    echo -e "${RED}Ошибка: go mod tidy попытался скачать зависимости в $dir:${NC}"
    echo "$OUTPUT"
    any_failed=1
  else
    echo -e "${GREEN}Успешно: $dir${NC}"
  fi

  popd > /dev/null
done

if [[ "$any_failed" == "1" ]]; then
  echo -e "${RED}Некоторые модули не прошли go mod tidy${NC}"
  exit 1
else
  echo -e "${GREEN}Все модули успешно обработаны go mod tidy${NC}"
fi
