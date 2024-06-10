PORT?=8080

all: run

dev:
	$(MAKE) -j 2 dev-unocss dev-air

dev-unocss:
	pnpm unocss -w

dev-air:
	air -build.pre_cmd 'make templ' \
		-build.include_ext 'templ,css,go' \
		-build.exclude_dir 'node_modules' \
		-proxy.enabled true \
		-proxy.app_port $(PORT) \
		-proxy.proxy_port $$(($(PORT) + 1)) \
		-- -p $(PORT) -d

dev-vercel:
	air -build.pre_cmd 'make build-vercel' \
		-build.include_ext 'templ,css' \
		-build.exclude_dir 'node_modules' \
		-build.cmd 'make build-vercel' \
		-build.bin './bin/vercel' \
		-proxy.enabled true \
		-proxy.app_port $(PORT) \
		-proxy.proxy_port $$(($(PORT) + 1)) \
		-- -p $(PORT)

run: bin/www
	./bin/www

run-vercel: bin/vercel
	./bin/vercel

build-static: templ
	pnpm unocss
	go run ./cmd/build/main.go

build-vercel: bin/vercel build-static

bin/www: main.go templ
	go build -o ./bin/www ./main.go

bin/vercel: cmd/vercel/main.go templ
	go build -o ./bin/vercel ./cmd/vercel/main.go

# For some reason "templ generate" does not detect the files, so this is a
# workaround.
TEMPL_FILES=$(patsubst %.templ, %_templ.go, $(wildcard **/*.templ))
templ: $(TEMPL_FILES)
	@echo Generating templ files
%_templ.go: %.templ
	templ generate -f $^ > /dev/null

clean:
	if [[ -d "dist" ]]; then rm -r ./dist; fi
	if [[ -d "tmp" ]]; then rm -r ./tmp; fi
	if [[ -d "bin" ]]; then rm -r ./bin; fi
	rm ./static/uno.css
	rm $(TEMPL_FILES)
