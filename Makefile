default: build run

build:
	go build

.PHONY: run
run:
	./mandelbrot
	open output.png

.PHONY: clean
clean:
	rm mandelbrot
	find . -name *.out -or -name *.log -or -name .*.swp -or -name .*.swo -or -name .DS_Store -or -name .swp | xargs -n 1 rm

