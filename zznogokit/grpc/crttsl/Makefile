
run:
    @go build -o a.out && ./a,out
    -@rm ./a.out

gen :
    protoc -I .--go_out=plugins=grpc:pb/. pb/helloworld.proto

clean:
  -rm *.pb.go