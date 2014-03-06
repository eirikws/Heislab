package communication

import (
    "time"
    "fmt"
)

type ElevButtons struct{
    u_buttons[3] bool
    d_buttons[3] bool
    c_buttons[4] bool
    stop_button bool
    current_floor int
    obstruction bool
    door_open bool
}


type IPandTimeStamp struct{
   IPadr string
   Timestamp time.Time
}

const ImAlivePort="30103"
const comPORT="30102"

func Communication(sendChan chan string, getChan chan string){
    ch:=make(chan []IPandTimeStamp)
    deletedIP:=make(chan string)
    var IPadr string
    var from,msg string
    var eleButtons ElevButtons
    var AliveList []IPandTimeStamp
    elevInfo:= make(map[string]ElevButtons)
    master:=make(chan string)
    MyIP:=getMyIP()
    BIP:=getBIP(MyIP)
    go sendImAlive(MyIP,BIP)
    go imAliveListener(MyIP,BIP,ch)
    go sendMsg(master,sendChan,MyIP)
    go recieveMsg(master,getChan,MyIP)
    go timeStampCheck(ch,deletedIP)
    for{
        select {
        case AliveList=<-ch:
            AliveList=IPsort(AliveList)
            master<-AliveList[0].IPadr
            ch<-AliveList
        case IPadr=<-deletedIP:
        		delete(elevInfo,IPadr)
        		fmt.Println(elevInfo)
        case <-time.After(time.Second*2):
        case msg=<-getChan:
            from,eleButtons=ReadMsg(msg)
            elevInfo[from]=eleButtons
           	spreadOrders(elevInfo)
           	for key,val:=range(elevInfo){
           		sendMsg("",sendChan,MyIP)
           	}
            fmt.Println(elevInfo)
        }
    }
}

func spreadOrders(info map[string]ElevButtons){
	u:=[]bool{false,false,false}
	d:=[]bool{false,false,false}
	var temp ElevButtons
	for _,val:=range(info){
		
		for i:=0; i<3;i++{
			if val.u_buttons[i]==true{
				fmt.Println("u", i, "true")
				u[i]=true
			}
			if val.d_buttons[i]==true{
				fmt.Println("d", i, "true")
				d[i]=true
			}
		}
			
	}
	for key,val:=range(info){
		temp=val
		for i:=0; i<3;i++{
			if u[i]==true{
				temp.u_buttons[i]=true
				info[key]=temp
				fmt.Println(key,"u2",i,"true")
			}
			if d[i]==true{
				temp.d_buttons[i]=true
				info[key]=temp
				fmt.Println(key,"d2",i,"true")
			}
		}
	}
}

func sendImAlive(MyIP, BIP string){
    msg:=makeMessage(MyIP,"I'm Alive")
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
            Smsg:=makeMessage(MyIP,msg)
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
            getChan<-Msg.from+Msg.info
        case <-time.After(time.Second*2):
        }
    }
}






