package communication

import (
        "net"
        "fmt"
)
func SendMsgTo(ipAdr string, port string,message Message){
        serverAddr, err := net.ResolveUDPAddr("udp",ipAdr+":"+port)
        con, err := net.DialUDP("udp", nil, serverAddr)
        if err != nil {
                fmt.Println("fail")
        }
                
        Bmessage:=msgToByte(message)
        for i:=0;i<5;i++{
                con.Write(Bmessage)
        }
}
func ListenerCon(ipAdr string, port string){
    serverAddr, err := net.ResolveUDPAddr("udp",ipAdr+":"+port)
    psock, err := net.ListenUDP("udp4", serverAddr)
    if err != nil { return }
    buf := make([]byte,1024)
    for {                 
        if err != nil { return }
        _, remoteAddr, _ := psock.ReadFromUDP(buf)
        if remoteAddr.IP.String() != MY_IP {
            fmt.Printf("%s\n",buf)
            }
        fmt.Printf("%s\n",buf)
    }              
}
