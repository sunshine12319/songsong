package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"study/example06/protoc"
)

type Client struct{
	conn net.Conn
	buf [8192]byte
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

func IntToBytes(n int) []byte {
	tmp := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer,binary.BigEndian,tmp)
	return bytesBuffer.Bytes()
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int32

	//BigEndian大端存取，就是字节读取顺序是按照高位到低位，通常大端
	binary.Read(bytesBuffer,binary.BigEndian,&tmp)
	return int(tmp)
}

// 读包
func (p *Client) readPackage() (msg protoc.Message,err error){
	//[0:4]表示长度
	n,err := p.conn.Read(p.buf[0:4])
	if n !=4{
		err =errors.New("read header failed")
	}

	fmt.Println("read package:",p.buf[0:4])
	// 读长度
	//buf := bytes.NewBuffer(p.buf[0:4])

	var packLen uint32
	packLen = binary.BigEndian.Uint32(p.buf[0:4])
	//err = binary.Read(buf,binary.BigEndian,&packLen)
	//if err != nil {
	//	fmt.Println("read package len failed")
	//	return
	//}

	//发送数据长度太大
	if int(packLen)> len(p.buf){

	}

	fmt.Printf("receive len:%d\n",packLen)
	// 读body,因为之前读了四个字节，因此前面的四个字节没了，数据流的读取
	n, err = p.conn.Read(p.buf[:packLen])
	if n != int(packLen){
		err = errors.New("read body failed")
		return
	}

	fmt.Printf("receive data:%s\n",string(p.buf[0:packLen]))
	err = json.Unmarshal(p.buf[0:packLen],&msg)
	if err != nil {
		fmt.Println("unmarshal failed, err:",err)
	}
	return
}

// 写包
func (p *Client) writePackage(data []byte) (err error){

	packLen := uint32(len(data))
	//buffer := bytes.NewBuffer(p.buf[:4])
	// int转换byte数组
	binary.BigEndian.PutUint32(p.buf[:4],packLen)
	//err = binary.Write(buffer,binary.BigEndian,packLen)
	//if err != nil {
	//	fmt.Println("write package len failed")
	//	return
	//}

	//发送数据,发头部
	n, err := p.conn.Write(p.buf[0:4])
	if err != nil{
		fmt.Println("write data failed")
		return
	}
	//发data
	n, err = p.conn.Write(data)
	if err != nil{
		fmt.Println("write data failed")
		return
	}
	//判断数据发送完没有
	if n != int(packLen) {
		fmt.Println("write data not finished")
		err = errors.New("write data not finished")
		return
	}

	return
}

// 处理请求。读包,判断信息分类处理
func (p *Client) Process()(err error) {
	for {
		var msg protoc.Message
		msg,err = p.readPackage()
		if err != nil {
			return err
		}

		err = p.processMsg(msg)
		if err != nil {
			return
		}
	}
}

// 判断客户端发送信息是哪一种
func (p *Client) processMsg(msg protoc.Message) (err error) {

	switch msg.Cmd {
	case protoc.UserLogin:
		err = p.login(msg)
	case protoc.UserRegister:
		err = p.register(msg)
	default:
		err = errors.New("unsupported message")
		return
	}
	return
}

//返回登录信息
func (p *Client) loginResp(err error)  {
	var respMsg protoc.Message
	respMsg.Cmd = protoc.UserLoginRes

	var loginRes protoc.LoginCmdRes
	loginRes.Code=200

	if err != nil{
		loginRes.Code =500
		loginRes.Error =fmt.Sprintf("%v",err)
	}

	data, err := json.Marshal(loginRes)
	if err != nil {
		fmt.Println("marshal failed,",err)
	}

	respMsg.Data=string(data)
	data, err =json.Marshal(respMsg)
	if err != nil{
		fmt.Println("marshal failed,",err)
		return
	}
	err = p.writePackage(data)
	if err != nil{
		fmt.Println("send failed,",err)
		return
	}
}

//通过客户端发送的账号密码，实现登录注册
func (p *Client) login(msg protoc.Message)(err error)  {
	defer func() {
		p.loginResp(err)
	}()
	fmt.Printf("recv user login request,data:%v\n",msg)
	var cmd protoc.LoginCmd
	err = json.Unmarshal([]byte(msg.Data),&cmd)
	if err != nil {
		return
	}

	_,err = mgr.Login(cmd.Id,cmd.Passwd)
	if err != nil{
		return
	}
	return
}

//判断message类别后，反序列化后实现注册
func (p *Client) register(msg protoc.Message)(err error)  {
	var cmd protoc.RegisterCmd
	err = json.Unmarshal([]byte(msg.Data),&cmd)
	if err != nil {
		return
	}

	err = mgr.Register(&cmd.User)
	if err != nil {
		return
	}

	return
}