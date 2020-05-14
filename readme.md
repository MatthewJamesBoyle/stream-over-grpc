# Stream large Files Over GRPC.

I couldn't find any good code examples of how to stream large files over GRPC... so here is one :) 


```protoc -I models/ models/streamingservice.proto --go_out=plugins==grpc:./models```
