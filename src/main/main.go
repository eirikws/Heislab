package main


import(
//    "fmt"
    //"net"
    ."runtime"
    com "./../communication"
    elev "./../elevator"

)

const PORT="30001"

func main(){
    MY_IP:=com.GetMyIP()
    BCAST_IP := com.GetBIP(MY_IP)
    GOMAXPROCS(NumCPU())
    msg:=com.MakeMessage(BCAST_IP,"muhhahahaha","testtesttesttesttesttesttesttesttesttesttest")
    go com.ListenerCon(BCAST_IP,PORT,MY_IP)
    go com.SendMsgTo(BCAST_IP,PORT,msg)
    for{
    }
}
