package communication

import (
        "net"
        "fmt"
     //   "time"
)
func getUDPcon(ipAdr string, port string) *net.UDPConn{
        serverAddr, err := net.ResolveUDPAddr("udp",ipAdr+":"+port)
        con, err := net.DialUDP("udp", nil, serverAddr)
        if err != nil {
                fmt.Println("fail")
        }
        return con
}

func listenerCon(ipAdr string, port string,MY_IP string,ch chan Message){
    serverAddr, err := net.ResolveUDPAddr("udp",ipAdr+":"+port)
    psock, err := net.ListenUDP("udp4", serverAddr)
    if err != nil { return }
    buf := make([]byte,1024)
    var msg Message
    for {
        if err != nil { return }
        psock.ReadFromUDP(buf)
        msg=byteToMsg(buf)
        ch<-msg
    }              
}

/*
func listenerCon2(master chan string, port string,MY_IP string,ch chan Message){
	var mst string
	mst=<-master
	master<-mst
	x:=0
	buf := make([]byte,1024)
    var msg Message
    serverAddr, _ := net.ResolveUDPAddr("udp",mst+":"+port)
	psock, _ := net.ListenUDP("udp4", serverAddr)
	select{
		case mst=<-master:
			master<-mst
			x=0
		case <-time.After(time.Millisecond*1):
			if x==0{
				serverAddr, _ := net.ResolveUDPAddr("udp",mst+":"+port)
				psock, _ := net.ListenUDP("udp4", serverAddr)
			}
			
			x=1
    	  	psock.ReadFromUDP(buf)
    	  	msg=byteToMsg(buf)
    	  	ch<-msg
    }
}*/
