CLI_SRC = src/client/client.go 
SER_SRC = src/server/main.go src/server/user.go src/server/server.go

all: client server

client: $(CLI_SRC)
	mkdir -p target
	go build -o target/client $^

server: $(SER_SRC)
	mkdir -p target
	go build -o target/server $^

clean:
	rm target/*
