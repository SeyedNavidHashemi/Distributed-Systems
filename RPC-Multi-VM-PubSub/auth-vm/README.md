# Authentication Server (VM2)

This service provides user authentication through JSON-RPC.

Responsibilities:
- Load users from `users.json`
- Validate login requests
- Return authentication results to the Web Server

## Requirements

- Go installed
- The bycrypt library of go downloaded and available. It can be downloaded using the following commands:
    ```
    go mod init auth-vm
    go env -w GOPROXY=https://package-mirror.liara.ir/repository/go/
    go env -w GOSUMDB=off
    go get golang.org/x/crypto/bcrypt
    ```

## User Database

Users are stored in:

```text
users.json
```

Passwords are stored as bcrypt hashes.

## Run

```bash
go run main.go
```

and the auth server should start running on **192.168.1.111**.

The RPC server listens on:

```text
TCP :9000
```

## Generate Password Hashes

If new users need to be added:

```bash
go run hash.go
```

Copy the generated hash into `users.json`.