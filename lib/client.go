package lib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"io"
	"log"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"golang.org/x/crypto/ssh"
)

type Node struct {
	User     string
	Password string
}

func NewNode(user, password string) *Node {
	node := new(Node)
	node.User = user
	node.Password = password
	return node
}

func (this *Node) Conn(addr string) (*ssh.Client, error) {

	authMethods := []ssh.AuthMethod{}

	keyboardInteractiveChallenge := func(
		user,
		instruction string,
		questions []string,
		echos []bool,
	) (answers []string, err error) {
		if len(questions) == 0 {
			return []string{}, nil
		}
		return []string{this.Password}, nil
	}

	authMethods = append(authMethods, ssh.KeyboardInteractive(keyboardInteractiveChallenge))
	authMethods = append(authMethods, ssh.Password(this.Password))

	sshConfig := &ssh.ClientConfig{
		User: this.User,
		Auth: authMethods,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s", addr), sshConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// ------------------------------------ //

// 客户端结构体
type Client struct {
	color    bool
	stat     bool
	addr     string
	user     string
	pass     string
	in       io.WriteCloser
	out      *bytes.Buffer
	session  *ssh.Session
	conn     *ssh.Client
	file     *os.File
	filePath string
	timeout  int64
}

// 创建客户端对象
func NewClient() *Client {
	client := new(Client)
	client.stat = false
	client.addr = "192.168.73.134"
	client.user = "zl"
	client.pass = "123"
	client.color = true
	return client
}

func (this *Client) IsConnected() bool {
	return this.stat
}

//失败后不在连接
func (this *Client) Connect(addr, user, pass string, handler func(client *Client, err error)) {
	this.user = user
	this.addr = addr
	this.pass = pass

	node := NewNode(this.user, this.pass)
	conn, err := node.Conn(this.addr)
	if err != nil {
		log.Printf("Connect SSH exception(%s)", err.Error())
		handler(this, err)
		time.Sleep(1 * time.Second)
		return
	}
	defer conn.Close()
	this.conn = conn

	log.Printf("Connect SSH(%s) success", this.addr)

	session, err := conn.NewSession()
	if err != nil {
		log.Printf("Connect SSH exception(%s)", err.Error())
		handler(this, err)
		time.Sleep(1 * time.Second)
		return
	}
	defer session.Close()

	this.filePath = "temp/" + uuid.New()

	file, _ := os.Create(this.filePath)

	this.in, err = session.StdinPipe()
	session.Stdout = file
	this.file = file
	if err != nil {
		log.Printf("Connect SSH exception(%s)", err.Error())
		handler(this, err)
		time.Sleep(1 * time.Second)
		this.DisConnect()
		return
	}
	// Color
	if this.color {

		modes := ssh.TerminalModes{
			ssh.ECHO:          1,     // disable echoing
			ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
			ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
		}

		if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
			handler(this, err)
			panic(err)
		}
	}

	if err := session.Shell(); err != nil {
		log.Printf("Connect SSH exception(%s)", err.Error())
		handler(this, err)
		time.Sleep(1 * time.Second)
		this.DisConnect()
		return
	}

	// Start
	this.stat = true
	this.session = session
	handler(this, nil)

	if err := this.session.Wait(); err == nil {
		fmt.Println("session disconnect")
		this.DisConnect()
	}
}

func (this *Client) DisConnect() {
	this.stat = false
	if err := this.in.Close(); err != nil {
		fmt.Println("in->" + err.Error())
	}
	if err := this.session.Close(); err != nil {
		fmt.Println("session->" + err.Error())
	}
	if err := this.file.Close(); err != nil {
		fmt.Println("file->" + err.Error())
	}
	if err := this.conn.Close(); err != nil {
		fmt.Println("conn->" + err.Error())
	}

	go func(path string) {
		for {
			err := os.Remove(this.filePath)
			if err == nil {
				break
			}
		}
	}(this.filePath)
}

func (this *Client) SendCmd(cmd string) {
	this.in.Write([]byte(cmd + "\n"))
}

func (this *Client) GetOutFile() string {
	fia, _ := os.Open(this.filePath)
	fda, _ := ioutil.ReadAll(fia)
	return string(fda)
}
