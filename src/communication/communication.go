package communication

import (
    "time"
    "fmt"
)

const ImAlivePort="30100"
const comPORT="30101"

func Communication(sendChan chan string, getChan chan string){
    ch:=make(chan []string)
    master:=make(chan string)
    MyIP:=getMyIP()
    BIP:=getBIP(MyIP)
    go sendImAlive(MyIP,BIP)
    go imAliveListener(MyIP,BIP,ch)
    go sendMsg(master,sendChan,MyIP)
    go recieveMsg(master,getChan,MyIP)
    var AliveList []string
    for{
        AliveList=<-ch
        IPsort(AliveList)
        master<-AliveList[0]
        
    }
}

func sendImAlive(MyIP, BIP string){
    msg:=makeMessage(MyIP,"ALL","I'm Alive")
    con:=getUDPcon(BIP,ImAlivePort)
    bmsg:=msgToByte(msg)
    for {
        con.Write(bmsg)
        time.Sleep(time.Second)
    }
}

func imAliveListener(MyIP, BIP string, ch chan []string){
    alivechan:=make(chan Message)
    go listenerCon(BIP,ImAlivePort, MyIP, alivechan)
    var newMsg Message
    var IPadr string
    var IPlist=[]string{MyIP}
    x:=0
    for{
        newMsg=<-alivechan
        IPadr=newMsg.from
        for _,IP:=range IPlist{
            if IP==IPadr{
                x=1
            }
            if x==0{
                IPlist=append(IPlist,IPadr)
            }
            ch<-IPlist
        }
    }
}

func sendMsg(master,sendChan chan string,MyIP string){
    var mst,msg string   
    for{
        select{
        case mst=<-master:
        case msg=<-sendChan:
            con:=getUDPcon(mst,comPORT)
            Smsg:=makeMessage(MyIP,mst,msg)
            Bmsg:=msgToByte(Smsg)
            con.Write(Bmsg)
        }
   }
}

func recieveMsg(master,getChan chan string,MyIP string){
    var mst string
    var Msg Message
    msg:=make(chan Message)
    mst=<-master
    go listenerCon(mst,comPORT,MyIP,msg)
    for{
        Msg=<-msg
        fmt.Println("Received Message")
        getChan<-Msg.from+Msg.to+Msg.info
    }
}






