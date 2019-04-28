package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

/*
 ****************************
 进程间通信
*******************************
 */
func main(){
	var pid  = os.Getpid()
	var ppid = os.Getppid()

	fmt.Printf("pid:%d\n",pid)
	fmt.Printf("ppid:%d\n",ppid)

    var sigchan = make (chan os.Signal,1)
    // signal.Notify 不按系统默认处理
    sigs := []os.Signal{syscall.SIGINT,syscall.SIGQUIT}
    signal.Notify(sigchan,sigs...)
    for sig := range sigchan {
    	fmt.Printf("Receive a signal:%s\n",sig)
    	//恢复默认处理方法
    	signal.Stop(sigchan)
	}

    //打开进程
    attr:= &os.ProcAttr{
    	Files:[]*os.File{os.Stdin,os.Stdout,os.Stderr},
	}

    p,err := os.StartProcess("C:\\Python27\\python.exe",[]string{"C:\\Python27\\python.exe","1.py"},attr)

    if  err != nil {
    	fmt.Println(err)
	}

    fmt.Println(p)
    pro,_ := os.FindProcess(p.Pid)
	fmt.Println(pro)                //&{5488 240 0}

	//杀死进程但不释放进程相关资源
	err = p.Kill()
	fmt.Println(err)

	//释放进程相关资源，因为资源释放凋之后进程p就不能进行任何操作，此后进程Ｐ的任何操作都会被报错
	err = p.Release()
	fmt.Println(err)

	//管道
	cmd1 := exec.Command("ps","aux")
	cmd2 := exec.Command("grep","apipe")

	var outputBuf1 bytes.Buffer
	cmd1.Stdout = &outputBuf1
	if err:=cmd1.Start();err!=nil {
		fmt.Printf("Erroe:the first command can not start:% %s",err)
		return
	}

	if err:=cmd1.Wait();err!=nil {
		fmt.Printf("Erroe: could not wait for the first command:% %s",err)
		return
	}

	cmd2.Stdin = &outputBuf1
	var outputBuf2 bytes.Buffer
	cmd2.Stdout = &outputBuf2
	if err:=cmd2.Start();err!=nil {
		fmt.Printf("Erroe:the second command can not start:% %s",err)
		return
	}

	if err:=cmd2.Wait();err!=nil {
		fmt.Printf("Erroe: could not wait for the second command:% %s",err)
		return
	}

	cmda := exec.Command("echo","-n","test,teste\n")
	if err:= cmda.Start(); err!=nil{
		fmt.Printf("err:%+v",err)
		return
	}

	// stdout 是 io.ReadCloser
	stdouto,err := cmda.StdoutPipe()
    if err!=nil {
		fmt.Printf("StdoutPipe err:%+v \n %+v",stdouto,err)
		return
	}

	out := make([]byte,30)
	ck,err := stdouto.Read(out)
	if err!=nil {
		fmt.Printf("stdouto.Read err:%+v",err)
		return
	}

	fmt.Printf("StdoutPipe err:%s \n",out[:ck])

}
