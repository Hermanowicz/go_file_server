hello:
	echo "Hello, I will build File server for you and run it on port 9000"

build:
	go build -o Server main.go

clear:
	rm -f Server

run_linux:
	./Server -port=9000