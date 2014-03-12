package communication

import (
    mst "./../master"
    "time"
    "fmt"
    gen "./../genDecl"
)

type IPandTimeStamp struct{
   IPadr string
   Timestamp time.Time
}

const ImAlivePort="30101"
const comPORT="30105"

func Communication(sendChanMaster chan string, getChan chan string){
    ch:=make(chan []IPandTimeStamp)
    receiveChan:=make(chan string)
    deletedIP:=make(chan string)
    getElevInfoChan:=make(chan map[string]gen.ElevButtons)
	orders:=make(chan string)
    //var LastMaster string
    //var IPadr string
    var IPadr,LastMaster,from,msg string
    var eleButtons gen.ElevButtons
    var AliveList []IPandTimeStamp
    elevInfo:= make(map[string]gen.ElevButtons)
    master:=make(chan string)
    MyIP:=getMyIP()
    BIP:=getBIP(MyIP)
    go sendImAlive(MyIP,BIP)
    go imAliveListener(MyIP,BIP,ch)
    go sendMsgToMaster(master,sendChanMaster,MyIP)
    go recieveMsg(master,receiveChan,MyIP)
    go timeStampCheck(ch,deletedIP)
    

	go mst.Master(master,getElevInfoChan,orders)
    for{
    	
        select {
        
        case AliveList=<-ch:
            AliveList=IPsort(AliveList)
            master<-AliveList[0].IPadr
            if LastMaster!=AliveList[0].IPadr{
            	sendChanMaster<-"U:"+gen.ElevButtonToStr(elevInfo[MyIP])
            }
            LastMaster=AliveList[0].IPadr
            ch<-AliveList
        case IPadr=<-deletedIP:
        	fmt.Println("Del IP")
        	delete(elevInfo,IPadr)
        	fmt.Println(elevInfo)
        case msg=<-receiveChan:
       		fmt.Println("got msg", msg[15:17])
        	switch {
        	case msg[15:17]=="C:":	
        		getChan<-msg[17:]
        	case msg[15:17]=="U:":
            	from,eleButtons=gen.ReadMsg(msg)
            	elevInfo[from]=eleButtons
           		spreadOrders(elevInfo)
           		getElevInfoChan<-elevInfo
           		elevInfo=<-getElevInfoChan
           		for key,val:=range(elevInfo){
           			SendMsgToThisGuy(key,"C:"+gen.ElevButtonToStr(val))
           		}
           		
           	}
        }
    }
}

func spreadOrders(info map[string]gen.ElevButtons){
	u:=[]bool{false,false,false}
	d:=[]bool{false,false,false}
	var temp gen.ElevButtons
	for _,val:=range(info){
		
		for i:=0; i<3;i++{
			if val.U_buttons[i]==true{
				u[i]=true
			}
			if val.D_buttons[i]==true{
				d[i]=true
			}
		}
	}
	for key,val:=range(info){
		temp=val
		for i:=0; i<3;i++{
			if u[i]==true{
				temp.U_buttons[i]=true
				info[key]=temp
			}
			if d[i]==true{
				temp.D_buttons[i]=true
				info[key]=temp
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
        IPlist=<-ch
    }
}

func sendMsgToMaster(master,sendChan chan string,MyIP string){
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
        }
    }
}

func SendMsgToThisGuy(IPadrTo string,msg string){
	MyIP:=getMyIP()
	con:=getUDPcon(IPadrTo,comPORT)
	Smsg:=makeMessage(MyIP,msg)
	Bmsg:=msgToByte(Smsg)
	con.Write(Bmsg)
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
            fmt.Println(Msg.from+Msg.info)
        case <-time.After(time.Second*2):
        }
    }
}






