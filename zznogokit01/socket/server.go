package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	SERVER_NETWORK = "tcp"
	SERVER_ADDRESS = "127.0.0.1:8085"
	DELIMITER      = '\t'
)

func printLog(role string, sn int, format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Printf("%s[%d]: %s", role, sn, fmt.Sprintf(format, args...))
}

func printServerLog(format string, args ...interface{}) {
	printLog("Server", 0, format, args...)
}

func read(conn net.Conn) (string, error) {
	readBytes := make([]byte, 1)
	var buffer bytes.Buffer
	for {
		_, err := conn.Read(readBytes)
		if err != nil {
			return "", err
		}
		readByte := readBytes[0]
		if readByte == DELIMITER {
			break
		}
		buffer.WriteByte(readByte)
	}
	return buffer.String(), nil
}

func write(conn net.Conn, content string) (int, error) {
	var buffer bytes.Buffer
	buffer.WriteString(content)
	buffer.WriteByte(DELIMITER)
	return conn.Write(buffer.Bytes())
}

/***
约定边界 分割数据块
接受int32类型的数据块 不符合要求的数据块，生成错误信息返回给客户端
符合要求生成描述给客户端  鉴别闲置的链接关闭 10秒钟没有数据传递可以关闭
 */
func main(){
    listener, err := net.Listen(SERVER_NETWORK, SERVER_ADDRESS)
	fmt.Printf("read error: %s\n",err)

	if err != nil {
		printServerLog("Listen Error: %s", err)
		return
	}
	defer listener.Close()

	printServerLog("Got listener for the server. (local address: %s)", listener.Addr())
	for {
		conn, err := listener.Accept() // 阻塞直至新连接到来。
		if err != nil {
			printServerLog("Accept Error: %s", err)
			continue
		}
		printServerLog("Established a connection with a client application. (remote address: %s)",
			conn.RemoteAddr())
		go handleConn(conn)
	}

    //var  dataBuffer bytes.Buffer
    //var data = make([]byte,10)
	//
    //for {
	//	// 非阻塞socket的accept read write函数 没有读到或已经写满缓冲区,都不会阻塞都会返回EAGAIT 应该忽略此错误
	//	// 非阻塞socket 是部分读  有多少读多少
	//	n, err := conn.Read(data)
	//	if err != nil {
	//		if err== io.EOF {
	//			fmt.Println("The client connect is closed")
	//		    conn.Close()
	//		} else {
	//			fmt.Printf("read error: %s\n",err)
	//		}
	//		break
	//	}
	//	//把数据追加 到buffer
	//	dataBuffer.Write(data[:n])
	//	content := string(data[:n])
	//	fmt.Println(content)
	//
	//	// bufio = buffered I/O
	//	reader := bufio.NewReader(conn)
	//	// 协商好的消息边界
	//	reader.ReadBytes('\n')
	//}
}

func handleConn(conn net.Conn){
	defer conn.Close()

	for {
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		strReq, err := read(conn)
		if err != nil {
			if err == io.EOF {
				printServerLog("The connection is closed by another side.")
			} else {
				printServerLog("Read Error: %s", err)
			}
			break
		}
		printServerLog("Received request: %s.", strReq)
		intReq, err := strToInt32(strReq)
		if err != nil {
			n, err := write(conn, err.Error())
			printServerLog("Sent error message (written %d bytes): %s.", n, err)
			continue
		}
		floatResp := cbrt(intReq)
		respMsg := fmt.Sprintf("The cube root of %d is %f.", intReq, floatResp)
		n, err := write(conn, respMsg)
		if err != nil {
			printServerLog("Write Error: %s", err)
		}
		printServerLog("Sent response (written %d bytes): %s.", n, respMsg)
	}
}

func strToInt32(str string) (int32, error) {
	num, err := strconv.ParseInt(str, 10, 0)
	if err != nil {
		return 0, fmt.Errorf("\"%s\" is not integer", str)
	}
	if num > math.MaxInt32 || num < math.MinInt32 {
		return 0, fmt.Errorf("%d is not 32-bit integer", num)
	}
	return int32(num), nil
}

func cbrt(param int32) float64 {
	return math.Cbrt(float64(param))
}
