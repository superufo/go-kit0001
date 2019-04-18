- **gomod**

go mod init  
go mod download  
go mod edit   
go mod graph   
go mod init   
go mod tidy   
go mod vendor  
go mod verify  
go mod why    



- ​       **docker**                                          

docker ps       

docker exec zkitregisterserver_etcd_1 /bin/sh -c "/usr/local/bin/etcd --version"   
docker exec zkitregisterserver_etcd_1 /bin/sh -c "ETCDCTL_API=3 /usr/local/bin/etcdctl version"    
docker exec zkitregisterserver_etcd_1 /bin/sh -c "ETCDCTL_API=3 /usr/local/bin/etcdctl endpoint health"   
docker exec zkitregisterserver_etcd_1 /bin/sh -c "ETCDCTL_API=3 /usr/local/bin/etcdctl put foo bar"   
docker exec zkitregisterserver_etcd_1 /bin/sh -c "ETCDCTL_API=3 /usr/local/bin/etcdctl get foo"   
docker exec zkitregisterserver_etcd_1  /bin/sh -c "ETCDCTL_API=3 /usr/local/bin/etcdctl get /services/book "       

test:
docker exec zkitregisterserver_etcd_1     /bin/sh -c "ETCDCTL_API=3 /usr/local/bin/etcdctl  put /se/bo  weee"   
docker exec zkitregisterserver_etcd_1     /bin/sh -c "ETCDCTL_API=3 /usr/local/bin/etcdctl get  /se/bo"       

​         

- etcd2命令
  存储:curl http://47.112.111.171:4001/v2/keys/testkey -XPUT -d value='testvalue'    
          curl -s http://47.112.111.171:4001/v2/keys/message2 -XPUT -d value='hello etcd' -d ttl=5    
  获取:curl http://47.112.111.171:4001/v2/keys/testkey    
  查看版本:curl  http://47.112.111.171:4001/version   
  删除: curl -s http://47.112.111.171:4001/v2/keys/testkey -XDELETE   
  监视 窗口1：curl -s http://47.112.111.171:4001/v2/keys/message2 -XPUT -d value='hello etcd 1'   
                     curl -s http://47.112.111.171:4001/v2/keys/message2?wait=true   
  窗口2： curl -s http://47.112.111.171:4001/v2/keys/message2 -XPUT -d value='hello etcd 2'
  自动创建key:   
            curl -s http://47.112.111.171:4001/v2/keys/message3 -XPOST -d value='hello etcd 1'    
            curl -s 'http://47.112.111.171:4001/v2/keys/message3?recursive=true&sorted=true'    
  创建目录：   
             curl -s http://47.112.111.171:4001/v2/keys/message8 -XPUT -d dir=true   
  删除目录：   
            curl -s 'http://47.112.111.171:4001/v2/keys/message7?dir=true' -XDELETE   
            curl -s 'http://47.112.111.171:4001/v2/keys/message7?recursive=true' -XDELETE   
  查看所有key:   
            curl -s http://47.112.111.171:4001/v2/keys/?recursive=true   
  存储数据：   
            curl -s http://47.112.111.171:4001/v2/keys/file -XPUT --data-urlencode value@upfile                                   

使用etcdctl客户端：   
存储:etcdctl set /mike/testkey "610" --ttl '100' --swap-with-value value   
获取:etcdctl get /mike/testkey   
更新:etcdctl update /mike/testkey "world" --ttl '100'   
删除:etcdctl rm /mike/testkey   
目录管理：   
        etcdctl mk /mike/testkey "hello"    类似set,但是如果key已经存在，报错   
        etcdctl mkdir /mike   
        etcdctl setdir /mike   
        etcdctl updatedir /mike   
        etcdctl rmdir /mike   
查看:etcdctl ls --recursive   
监视:etcdctl watch mykey  --forever   +  etcdctl update mykey "hehe"                                                                                                                                                           监视目录下所有节点的改变:    
       etcdctl exec-watch --recursive /foo -- sh -c "echo hi"   
       etcdctl exec-watch mykey -- sh -c 'ls -al'    +    etcdctl update mykey "hehe"   
       etcdctl member list   