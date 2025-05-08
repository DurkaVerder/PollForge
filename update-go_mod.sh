#!/bin/bash

cd services

echo -e "\n\nStart updating go.mod in all service directories\n\n"

for dir in */; do
    if [ -d "$dir" ]; then
        echo "Updating go.mod in $dir"
        cd "$dir"
        go mod tidy
        cd ..
    fi
done

echo
echo -e "\n\nUpdating go.mod in all services completed \n\n"