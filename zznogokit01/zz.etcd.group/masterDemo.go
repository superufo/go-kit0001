package zz_etcd_group

import (
	"log"
	"time"
	"fmt"
	dis "zz.etcd.group/discovery"
)

/***
   如果设置的监听地址是0.0.0.0那么我们无论是通过IP10.0.75.1 还是10.1.1.12都是可以访问该服务的   
   Wireless LAN adapter WLAN:
   Connection-specific DNS Suffix  . :
   Link-local IPv6 Address . . . . . : fe80::2ce7:ed9c:29e5:ade7%5
   IPv4 Address. . . . . . . . . . . : 10.0.75.1
   Subnet Mask . . . . . . . . . . . : 255.255.255.0
   Default Gateway . . . . . . . . . : 192.168.0.1

    Ethernet adapter vEthernet (DockerNAT):

   Connection-specific DNS Suffix  . :
   Link-local IPv6 Address . . . . . : fe80::2494:7cab:ba51:265a%30
   IPv4 Address. . . . . . . . . . . : 10.0.75.1
   Subnet Mask . . . . . . . . . . . : 255.255.255.0
   Default Gateway . . . . . . . . . :
 */
func main() {

	m, err := dis.NewMaster([]string{
		"http://10.0.75.1:32792",
		"http://10.0.75.1:32796",
		"http://10.0.75.1:32794",
	}, "services/")

	if err != nil {
		log.Fatal(err)
	}

	for {
		for k, v := range  m.Nodes {
			fmt.Printf("node:%s, ip=%s\n", k, v.Info.IP)
		}
		fmt.Printf("nodes num = %d\n",len(m.Nodes))
		time.Sleep(time.Second * 5)
	}
}
