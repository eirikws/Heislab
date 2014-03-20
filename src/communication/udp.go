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
                fmt.Println("motherfucking fail")
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
        if err != nil {
            fmt.Println("ListenerCon Fail")
            return 
        }
        psock.ReadFromUDP(buf)
        msg=byteToMsg(buf)
        ch<-msg
    }
}
