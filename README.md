# Chat room over TCP

It's very naive, non-scalable implementation of chat room server.

To run server on `localhost:8080`:
```bash
go run .
```

Connect to server:
```bash
nc localhost 8080
```

Write message and press `enter`. Server broadcasts message to all connected clients.

TODO:
- todo in comments
- tests unit & e2e
