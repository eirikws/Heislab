package communication

import (
    "time"
    "fmt"
)

type IPandTimeStamp struct{
   IPadr string
   Timestamp time.Time
}

const ImAlivePort="30100"
const comPORT="30101"

func Communication(sendChan chan string, getChan chan string){
    ch:=make(chan []IPandTimeStamp)
    
    var AliveList []IPandTimeStamp
    master:=make(chan string)
    MyIP:=getMyIP()
    BIP:=getBIP(MyIP)
    go sendImAlive(MyIP,BIP)
    go imAliveListener(MyIP,BIP,ch)
    go sendMsg(master,sendChan,MyIP)
    go recieveMsg(master,getChan,MyIP)
    go timeStampCheck(ch)
    for{
        select {
        case AliveList=<-ch:
            //fmt.Println("read CH")
            AliveList=IPsort(AliveList)
            fmt.Println(AliveList)
            master<-AliveList[0].IPadr
            ch<-AliveList
         //   fmt.Println("wrote CH")
        case <-time.After(time.Second*2):
        }
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

func imAliveListener(MyIP, BIP string, ch chan []IPandTimeStamp){
    alivechan:=make(chan Message)
    go listenerCon(BIP,ImAlivePort, MyIP, alivechan)
    var newMsg Message
    var IPadr string
    var IPlist=[]IPandTimeStamp{{MyIP,time.Now()}}
    var iptime IPandTimeStamp
    x:=0
    for{
        x=0
        newMsg=<-alivechan
        IPadr=newMsg.from
        for i,IP:=range IPlist{
            if IP.IPadr==IPadr{
            
                x=1
                IPlist[i].Timestamp=time.Now().Add(2200*time.Millisecond)}
            }
        if x==0{
            fmt.Println("Appending")
            iptime=IPandTimeStamp{IPadr,time.Now().Add(2200*time.Millisecond)}
            IPlist=append(IPlist,iptime)
        }
        ch<-IPlist
       // fmt.Println("I'm alive: wrote to ch")
        IPlist=<-ch
       // fmt.Println("I'm Alive: Read from ch")
       // fmt.Println("I'm Alive IPlist :", IPlist)
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






