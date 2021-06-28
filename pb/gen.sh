# mac or linux run proto file
# gopb for go; cspb for csharp
protoc --go_out=./gopb/ ./req.proto
protoc --go_out=./gopb/ ./res.proto
protoc --go_out=./gopb/ ./push.proto

protoc --csharp_out=./cspb/ ./req.proto
protoc --csharp_out=./cspb/ ./res.proto
protoc --csharp_out=./cspb/ ./push.proto