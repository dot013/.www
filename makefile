all: run

dev:
	air

run: bin/www
	./bin/www

run-vercel: bin/vercel
	./bin/vercel

build-static: templ
	go run ./cmd/build/main.go

build-vercel: build-static

bin/www: main.go templ
	go build -o ./bin/www ./main.go

bin/vercel: cmd/vercel/main.go templ build-vercel
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
	rm $(TEMPL_FILES)
