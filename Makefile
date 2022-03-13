hello:
	echo "Hello, world"

build:
	go build -o Server main.go

cert:
	go run cert_gen.go -host="localhost" -rsabits=3072