# File Server (VM3)

This service provides files and images to the distributed system.

Responsibilities:
- Host static files
- Serve images over HTTP

## Requirements

- Go installed

## Files Directory

Place files inside:

```text
files/
```

Example:

```text
files/image.jpg
```

## Run

```bash
go run main.go
```

The server listens on:

```text
http://192.168.1.108:8081
```

Example image URL:

```text
http://192.168.1.108:8081/files/image.jpg
```

## Test

Open the image URL directly in a browser:

```text
http://192.168.1.108:8081/files/image.jpg
```

and the following image should appear:

![image.png](./files/image.jpg)