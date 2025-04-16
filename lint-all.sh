#!/bin/bash

set -euo pipefail

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

CONFIG_PATH=".golangci.yml"

echo -e "${GREEN}Проверка кода...${NC}"

MODULE_DIRS=$(find . -mindepth 2 -name "go.mod" -exec dirname {} \;)

count=0
allModules=0
for dir in $MODULE_DIRS; do
  if [[ -d "$dir" ]]; then
    allModules=$((allModules + 1))
  fi
done

if [[ ! -f "$CONFIG_PATH" ]]; then
  echo -e "${RED}Конфиг $CONFIG_PATH не найден!${NC}"
  exit 1
fi

# Проходимся по каждому модулю
for dir in $MODULE_DIRS; do
  echo -e "${GREEN}▶Проверка модуля: $dir${NC}"
  pushd "$dir" > /dev/null
  if golangci-lint run --config="../$CONFIG_PATH" ./...; then
    echo -e "${GREEN}Проверка успешна: $dir${NC}"
    count=$((count + 1))
  else
    echo -e "${RED}Проблемы найдены в: $dir${NC}"
  fi
  echo -e ""
  popd > /dev/null
done

if [[ $count -eq $allModules ]]; then
  echo -e "${GREEN}Все модули прошли проверку!${NC}"
else
  echo -e "${RED}Проверка завершена с ошибками.${NC}"
fi

echo -e ""

if [[ $count -eq $allModules ]]; then
  echo -e "${GREEN}Успешно обработаны все модули.${NC}"
  exit 1
fi

echo -e "${RED}Обработано $count из $allModules с ошибками.${NC}"

