package zz_etcd_group

import (
	"fmt"
	"time"
	dis "zz.etcd.group/discovery"
	"log"
)

/***
λ docker ps
CONTAINER ID        IMAGE                 COMMAND                  CREATED             STATUS              PORTS                                              NAMES
a8169434bef3        quay.io/coreos/etcd   "etcd -name etcd2 -a…"   41 seconds ago      Up 37 seconds       0.0.0.0:32773->2379/tcp, 0.0.0.0:32772->2380/tcp   etcd2
12ca0bb9666f        quay.io/coreos/etcd   "etcd -name etcd3 -a…"   41 seconds ago      Up 38 seconds       0.0.0.0:32771->2379/tcp, 0.0.0.0:32770->2380/tcp   etcd3
c704a36cb8e6        quay.io/coreos/etcd   "etcd -name etcd1 -a…"   41 seconds ago      Up 38 seconds       0.0.0.0:32769->2379/tcp, 0.0.0.0:32768->2380/tcp   etcd1
 */
func main() {
	serviceName := "s-test111111"
	serviceInfo := dis.ServiceInfo{IP:"10.0.75.1"}

	s, err := dis.NewService(serviceName, serviceInfo,[]string {
		"http://10.0.75.1:32792",
		"http://10.0.75.1:32796",
		"http://10.0.75.1:32794",
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("name:%s, ip:%s\n", s.Name, s.Info.IP)


	go func() {
		time.Sleep(time.Second*20)
		s.Stop()
	}()

	s.Start()
}
