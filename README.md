## SockComm
A simple socket implementation between client & server

### Server
```bash
go run ./cmd/main.go server 9090
```

Server sends welcome message on new client connection, then pings it every 5 secs.

### Client
```bash
go run ./cmd/main.go client 9090
```

Client reads input from STDIN and sends it to the server.