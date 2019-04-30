创建 server.key             

openssl genrsa -out server.key 2048             

生成 server.crt                

 openssl req -new -x509 -days 3650 -subj "/C=GB/L=China/O=grpc-server/CN=server.grpc.io"  -key server.key -out server.crt           





创建 client.key         

openssl genrsa -out client.key 2048                

生成 client.crt                    

openssl req -new -x509 -days 3650 -subj "/C=GB/L=China/O=grpc-client/CN=client.grpc.io"  -key client.key -out client.crt                     





