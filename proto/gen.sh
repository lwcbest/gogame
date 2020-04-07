protoc --go_out=./goproto/ ./req.proto
protoc --go_out=./goproto/ ./res.proto
protoc --go_out=./goproto/ ./push.proto

protoc --csharp_out=./csproto/ ./req.proto
protoc --csharp_out=./csproto/ ./res.proto
protoc --csharp_out=./csproto/ ./push.proto