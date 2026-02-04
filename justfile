GO_DIR := "go"
SVELTE_DIR := "web"
UTILS_DIR := "utils"

default:
    @just --list

go:
    cd {{GO_DIR}} && air

web:
    cd {{SVELTE_DIR}} && bun run dev

both:
    just -j 2 serve-go dev-web

build-go:
    cd {{GO_DIR}} && go build -o ./bin ./cmd 

build-web:
    cd {{SVELTE_DIR}} && bun run build

cloc:
    {{UTILS_DIR}}/cloc.sh {{UTILS_DIR}}/cloc.log

base16:
    {{UTILS_DIR}}/base16/do.sh {{SVELTE_DIR}}/src/css/

clean-web:
    cd {{SVELTE_DIR}} && rm -rf .svelte-kit node_modules dist build

clean-go:
    go clean -cache -modcache && rm {{GO_DIR}}/bin

clean:
    just clean-web clean-go

typegen:
    cd {{GO_DIR}} && tygo generate
