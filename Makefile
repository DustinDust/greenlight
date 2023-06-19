run.dev:
	ho run ./cmd/api/ 

build.bin:
	rm -rf bin
	mkdir bin && go build -o bin ./cmd/api/
