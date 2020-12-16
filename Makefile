default: build run

build:
	go build

.PHONY: run
run:
	./mandelbrot

release: build
	chmod 777 mandelbrot
	zip mandelbrot.zip mandelbrot

wasm:
	GOOS=js GOARCH=wasm go build -o www/mandelbrot.wasm github.com/keyan/mandelbrot
	cp $(go env GOROOT)/misc/wasm/wasm_exec.js .

local:
	python3 -m http.server --directory www/

.PHONY: clean
clean:
	rm mandelbrot
	find . -name *.out -or -name *.log -or -name .*.swp -or -name .*.swo -or -name .DS_Store -or -name .swp | xargs -n 1 rm
