echo "Запуск go mod tidy для всех модулей..."

find . -mindepth 2 -name "go.mod" -exec dirname {} \; | while read dir; do
    echo "Модуль: $dir"
    
    pushd "$dir" > /dev/null

    go mod tidy
    if ! git diff --quiet go.mod go.sum; then
        echo "❌ go.mod/go.sum изменены в $dir. Запусти 'go mod tidy'."
        git diff go.mod go.sum
        exit 1
    fi

    popd > /dev/null

done