## Start Server

```
go run *.go
```

## Endpoints

# Get Current Block

```
curl http://127.0.0.1:3000/get_current_block
-> {"block_number":241116704}
```

# Subscribe

```
curl http://127.0.0.1:3000/subscribe/0x5bE36859685e4f75871EDa3dE88BB9eDeCCa8da1
-> {"subscribed":true}
```

# Get Transactions

```
curl http://127.0.0.1:3000/get_transactions/0x5bE36859685e4f75871EDa3dE88BB9eDeCCa8da1
-> {"transactions":[]}
```