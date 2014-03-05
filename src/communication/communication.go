package communication

import (
    "time"
    "fmt"
)

type IPandTimeStamp struct{
   IPadr string
   Timestamp time.Time
}

const ImAlivePort="30103"
const comPORT="30102"

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
       // fmt.Println("enter coms")
        select {
        case AliveList=<-ch:
         //   fmt.Println("read CH")
            AliveList=IPsort(AliveList)
         //   fmt.Println(AliveList)
            master<-AliveList[0].IPadr
            ch<-AliveList
        //    fmt.Println("wrote CH")
        case <-time.After(time.Second*2):
        case 
        }
       // fmt.Println("exit coms")
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
      //  fmt.Println("enter loop")
        newMsg=<-alivechan
      //  fmt.Println("enter loop2")
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
      //  fmt.Println("I'm alive: wrote to ch")
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
   //         fmt.Println("sendMsg: New master")
        case msg=<-sendChan:
            fmt.Println("sending to master")
            con:=getUDPcon(mst,comPORT)
            Smsg:=makeMessage(MyIP,mst,msg)
            Bmsg:=msgToByte(Smsg)
            con.Write(Bmsg)
            fmt.Println("sendt to master")
        }
    }
}

func recieveMsg(master,getChan chan string,MyIP string){
    var Msg Message
    msg:=make(chan Message)
    go listenerCon("",comPORT,MyIP,msg)
    for{
         
        select{
  //      case mst=<-master:
  //          fmt.Println("recievemsg: new master")
        case Msg=<-msg:
            getChan<-Msg.from+Msg.to+Msg.info
        case <-time.After(time.Second*2):
        }
    }
}






