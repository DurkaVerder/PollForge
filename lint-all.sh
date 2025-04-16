#!/bin/bash

set -euo pipefail

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

CONFIG_PATH=".golangci.yml"

echo -e "${GREEN}🔍 Поиск всех Go-модулей...${NC}"

MODULE_DIRS=$(find . -mindepth 2 -name "go.mod" -exec dirname {} \;)

count=0
allModules=0
for dir in $MODULE_DIRS; do
  if [[ -d "$dir" ]]; then
    allModules=$((allModules + 1))
  fi
done

if [[ ! -f "$CONFIG_PATH" ]]; then
  echo -e "${RED}❌ Конфиг $CONFIG_PATH не найден!${NC}"
  exit 1
fi

# Проходимся по каждому модулю
for dir in $MODULE_DIRS; do
  echo -e "${GREEN}▶️  Линтинг модуля: $dir${NC}"
  pushd "$dir" > /dev/null
  if golangci-lint run --config="../$CONFIG_PATH" ./...; then
    echo -e "${GREEN}✅ Линтинг успешен: $dir${NC}"
    count=$((count + 1))
  else
    echo -e "${RED}❌ Проблемы найдены в: $dir${NC}"
  fi
  echo -e ""
  popd > /dev/null
done

if [[ $count -eq $allModules ]]; then
  echo -e "${GREEN}✅ Все модули прошли линтинг!${NC}"
else
  echo -e "${RED}❌ Линтинг завершен с ошибками.${NC}"
fi

echo -e ""

if [[ $count -eq $allModules ]]; then
  echo -e "${GREEN}Успешно обработано все $allModules модулей.${NC}"
  exit 1
fi

echo -e "${RED}Обработано $count из $allModules модулей с ошибками.${NC}"

