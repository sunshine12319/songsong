package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"study/example06/protoc"
)

// 协议
type Message struct {
	Cmd string `json:"cmd"`
	Data string `json:"data"`
}

// 客户端登录信息
type LoginCmd struct{
	Id int `json:"user_id"`
	Passwd string `json:"passwd"`
}


func readPackage(conn net.Conn) (msg Message,err error){
	var buf [8192]byte
	n,err := conn.Read(buf[0:4])
	if n !=4{
		err =errors.New("read header failed")
	}

	fmt.Println("read package:",buf[0:4])
	// 读长度
	//buf := bytes.NewBuffer(p.buf[0:4])

	var packLen uint32
	packLen = binary.BigEndian.Uint32(buf[0:4])
	//err = binary.Read(buf,binary.BigEndian,&packLen)
	//if err != nil {
	//	fmt.Println("read package len failed")
	//	return
	//}

	fmt.Printf("receive len:%d\n",packLen)
	// 读body,因为之前读了四个字节，因此前面的四个字节没了，数据流的读取
	n, err = conn.Read(buf[:packLen])
	if n != int(packLen){
		err = errors.New("read body failed")
		return
	}

	fmt.Printf("receive data:%s\n",string(buf[0:packLen]))
	err = json.Unmarshal(buf[0:packLen],&msg)
	if err != nil {
		fmt.Println("unmarshal failed, err:",err)
	}
	return
}

func login (conn net.Conn)(err error){
	var msg Message
	msg.Cmd = protoc.UserLogin

	var loginCmd LoginCmd
	loginCmd.Id = 1
	loginCmd.Passwd = "123456789"

	data, err := json.Marshal(loginCmd)
	if err != nil {
		return
	}

	msg.Data=string(data)
	data,err = json.Marshal(msg)
	if err != nil {
		return
	}

	var buf [4]byte
	packLen := uint32(len(data))
	//用buffer时候会多出4个字节，因此传输过程总是读0
	//buffer := bytes.NewBuffer(buf[:])

	fmt.Println("packLen:",packLen)
	//err = binary.Write(buffer,binary.BigEndian,packLen)
	binary.BigEndian.PutUint32(buf[:],packLen)
	if err!= nil{
		fmt.Println("write package len failed")
		return
	}
	fmt.Println("packLen",buf)
	n, err := conn.Write(buf[:])
	if err != nil || n != 4{
		fmt.Println("write data failed")
		return
	}
	_, err = conn.Write(data)
	if err != nil {
		return
	}

	msg, err = readPackage(conn)
	if err != nil{
		fmt.Println("read package failed, err:",err)
		return
	}
	fmt.Println(msg)
	return
}

func main() {
	conn, err := net.Dial("tcp","localhost:10000")
	if err != nil{
		fmt.Println("Error dialing",err.Error())
		return
	}
	defer conn.Close()

	err=login(conn)
	if err != nil {
		fmt.Println("login failed,err:", err)
		return
	}
}
