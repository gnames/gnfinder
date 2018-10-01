# `Protobuf` for `bhlindex`

## Install `protoc` v3 on Ubuntu

Check the latest [releases](https://github.com/google/protobuf/releases)

```bash
# Make sure you grab the latest version, change XXXX accordingly
curl -OL https://github.com/google/protobuf/releases/download/v3.XXXX/protoc-3.XXXX-linux-x86_64.zip

# Unzip
unzip protoc-3.XXXX-linux-x86_64.zip -d protoc3

# Move protoc to /usr/local/bin/
sudo mv protoc3/bin/* /usr/local/bin/

# Move protoc3/include to /usr/local/include/
sudo mv protoc3/include/* /usr/local/include/

# Optional: change owner
sudo chown [user] /usr/local/bin/protoc
sudo chown -R [user] /usr/local/include/google
```

## Install Go protobuf etc

From gnfinder root install all dependencies:

```bash
go get -u github.com/golang/protobuf/protoc-gen-go
go get -u ./...
```

## How to create/update

```bash
cd protob
protoc -I . ./protob.proto --go_out=plugins=grpc:.
```
