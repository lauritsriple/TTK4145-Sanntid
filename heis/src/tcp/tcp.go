package tcp

import ("fmt"
	"net"
	"os"
)

func checkError(err error) {
	if err  != nil {
		fmt.Println("Error: " , err)
		os.Exit(0)
	}	
}

//Master: Connect to slaves.
//Slave: Reconnect to master.

func SlaveHandler
func MasterHandler

func Init(ipaddr,port,recv chan<-string,msg <-chan string){ //Write to recv, read from msg
	rAddr,err := net.ResolveTCPAddr("tcp",ipaddr+":"+port)
	checkError(err)
	conn,err := netDialTCP("tcp",nil,rAddr)
	checkError(err)
	defer conn.Close()
	laddr,err:= net.ResolveTCPAddr("tcp",":"+port)
	checkError(err)
	listen,err:=net.ListenTCP("tcp",laddr)
	checkError(err)


func tcp_receive() {
	laddr, err := net.ResolveTCPAddr("tcp",":30000")
	CheckError(err)

	ln, err := net.ListenTCP("tcp",laddr)
	CheckError(err)

	conn,err:=ln.Accept()
	CheckError(err)
	fmt.Println("here")

	buf:= make([]byte,1024)
	
	read_len,err := conn.Read(buf[:])
	if read_len > 0{
		fmt.Println("received tcp:",string(buf[0:read_len]))
		conn.Write([]byte("backloop"))
	}
}
