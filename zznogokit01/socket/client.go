package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
	"math/rand"
)

const (
	SERVER_NETWORK = "tcp"
	SERVER_ADDRESS = "127.0.0.1:8085"
	DELIMITER      = '\t'
)

func main () {
	// network 可以为 udp udp4 udp6 ip ip4  ip6 unix unixpackage unixgram
	conn,err := net.Dial(SERVER_NETWORK, SERVER_ADDRESS)
	fmt.Printf("%s\n",err)
	defer conn.Close()

	if err != nil {
		printClientLog(1, "Dial Error: %s", err)
		return
	}
	defer conn.Close()
	printClientLog(1, "Connected to server. (remote address: %s, local address: %s)",
		conn.RemoteAddr(), conn.LocalAddr())
	time.Sleep(200 * time.Millisecond)

	//conn,err := net.DialTimeout("tcp","127.0.0.1:8822",time.Second*10)

	// 非阻塞socket的accept read write函数 没有读到或已经写满缓冲区,都不会阻塞都会返回EAGAIT 应该忽略此错误
	//conn.Write()

	conn.SetDeadline(time.Now().Add(120*time.Second))
	//writer := bufio.NewWriter(conn)

	for {
		time.Sleep(time.Second*2)
		content := RandContent()
		n,err :=  write(conn, content)
		if err!=nil{
			printClientLog(1,"Write Error: %s", err)
			continue
		}
		printClientLog(1,"Sent request (written %d bytes): %s.", n,content)
	}

	for {
		time.Sleep(time.Second*2)

		n,err :=  read(conn)
		strResp, err := read(conn)
		if err != nil {
			if err == io.EOF {
				printClientLog(1, "The connection is closed by another side.")
			} else {
				printClientLog(1, "Read Error: %s", err)
			}
			break
		}
		printClientLog(1, "read response: %s. 字符数：%d ", strResp,n)
	}

	//fmt.Println(writer)
}


func printLog(role string,sn int,format string,args ...interface{}){
	if !strings.HasSuffix(format,"\n"){
		format += "\n"
	}

	fmt.Printf("%s[%d]: %s", role, sn, fmt.Sprintf(format, args...))
}


func printClientLog(sn int, format string, args ...interface{}) {
	printLog("Client", sn, format, args...)
}


func read(conn net.Conn)(string,error) {
	var readbyte = make([]byte,1)
	var buff  bytes.Buffer

	for {
		_,err := conn.Read(readbyte)
		if err !=nil {
			printClientLog(1, " read err: %+v.", err)
		}

		readByte := readbyte[0]
		if readByte == DELIMITER {
			break
		}
		buff.WriteByte(readByte)
	}

	return buff.String(),nil

}


func write(conn net.Conn, content string) (int, error) {
	var buffer bytes.Buffer
	buffer.WriteString(content)
	buffer.WriteByte(DELIMITER)
	return conn.Write(buffer.Bytes())
}

func RandContent() (content string){
	var carr = []string{"hello","kill","曹植","刘畅","linda","huose","kilo","啦啦啦","古汉语","老虎机","罗姆尼","拉米娜","狮山村","luoken","mine","醴陵市","手动档","生产队","细胞","法国货","射雕","天然","哈巴狗","东方宾馆","任天狗","fg短短的","半乳糖","还没考","内核码","璐璐","笨蛋","一剪梅","婷婷"}


	for i:=0;i<5;i++{
		var nu=rand.Intn(30)
		content += carr[nu]
	}
	return content
}

