第二版的Protobuf有个默认值特性，可以为字符串或数值类型的成员定义默认值,第三版的Protobuf
中不再支持默认值特性，但是我们可以通过扩展选项自己模拟默认值.
default.proto 就是扩展规则。

在开源社区中，github.com/mwitkow/go-proto-validators 已经基于Protobuf的扩展
特性实现了功能较为强大的验证器功能。要使用该验证器首先需要下载其提供的代
码生成插件：
$ go get github.com/mwitkow/go-proto-validators/protoc-gen-govalidators
然后基于go-proto-validators验证器的规则为Message成员增加验证规则

validators.proto 就是验证规则。

validator.field表示扩展是validator包中定义的名为
field扩展选项。validator.field的类型是FieldValidator结构体，在导入的
validator.proto文件中定义

protoc \
--proto_path=${GOPATH}/src \
--proto_path=${GOPATH}/src/github.com/google/protobuf/src \
--proto_path=. \
--govalidators_out=. --go_out=plugins=grpc:.\
validators.proto
windows:替换 ${GOPATH} 为 %GOPATH% 即可
以上的命令会调用protoc-gen-govalidators程序，生成一个独立的名为
validators.validator.pb.go的文件

protoc --proto_path=%GOPATH%/src --proto_path=%GOPATH%/src\google.golang.org\genproto\protobuf --proto_path=. --govalidators_out=. --go_out=plugins=grpc:. validators.proto
编译有问题：
D:\gopromod\zz.wt\zznogokit\grpc\validate\pb (master -> origin)
λ protoc --proto_path=%GOPATH%\src --proto_path=%GOPATH%\src\google.golang.org\genproto\protobuf --proto_path=. --govalidators_out=. --go_out=plugins=grpc:. validators.proto
google/protobuf/descriptor.proto: File not found.
github.com/mwitkow/go-proto-validators/validator.proto: Import "google/protobuf/descriptor.proto" was not found or had errors.
github.com/mwitkow/go-proto-validators/validator.proto:18:8: "google.protobuf.FieldOptions" is not defined.
validators.proto: Import "github.com/mwitkow/go-proto-validators/validator.proto" was not found or had errors.

