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


const ImAlivePort="30108"
const comPORT="30107"
const lostInternetIP="192.168.0.1"


func Communication(sendChanMaster chan gen.ElevButtons, getChan chan gen.ElevButtons){
    ch:=make(chan []IPandTimeStamp)
    receiveChan:=make(chan Message)
    deletedIP:=make(chan string)
    getElevInfoChan:=make(chan map[string]gen.ElevButtons)
	orders:=make(chan string)
    var IPadr,LastMaster,from string
    var msg Message
    var eleButtons gen.ElevButtons
    var AliveList []IPandTimeStamp
    elevInfo:= make(map[string]gen.ElevButtons)
    master:=make(chan string)
    MyIP:=getMyIP()
    BIP:=getBIP(MyIP)
    go sendImAlive(MyIP,BIP)
    go imAliveListener(MyIP,BIP,ch)
    go sendMsgToMaster(master,sendChanMaster,MyIP,receiveChan)
    go recieveMsg(receiveChan,MyIP)
    go timeStampCheck(ch,deletedIP,MyIP)
    go mst.Master(master,getElevInfoChan,orders,MyIP)
    fmt.Println(MyIP)
    for{
      //  fmt.Println("coms",getMyIP(),LastMaster)
        select {
        case AliveList=<-ch:
            AliveList=IPsort(AliveList)
            master<-AliveList[0].IPadr
            if LastMaster!=AliveList[0].IPadr{
            	sendChanMaster<-elevInfo[MyIP]
            }
            LastMaster=AliveList[0].IPadr
            ch<-AliveList
        case IPadr=<-deletedIP:
        	delete(elevInfo,IPadr)
        case msg=<-receiveChan:
       		fmt.Println("got msg", msg.typ)
        	switch {
        	case msg.typ=="C":
        		getChan<-stringToButton(msg.info)
        	case msg.typ=="U":
            	from,eleButtons=msg.from,stringToButton(msg.info)
            	turnOffLights(elevInfo,from,eleButtons)
            	elevInfo[from]=eleButtons
           		spreadOrders(elevInfo)
           		getElevInfoChan<-elevInfo
           		elevInfo=<-getElevInfoChan
           		for key,val:=range(elevInfo){
           		    go sendMsgToThisGuy(key,val,MyIP,receiveChan)
           		}
           	
           	}
        }
    }
}

func turnOffLights(InfoMap map[string]gen.ElevButtons, IPadrFrom string, newInfo gen.ElevButtons){
	u:=[]bool{false,false,false}
	d:=[]bool{false,false,false}
	var dummyvar gen.ElevButtons
	for i:=0 ; i<gen.N_FLOORS-1 ; i++{
		if InfoMap[IPadrFrom].U_buttons[i] && !newInfo.U_buttons[i]{
			u[i]=true
		}
		if InfoMap[IPadrFrom].D_buttons[i] && !newInfo.D_buttons[i]{
			d[i]=true
		}
	}
	
	for i,val:=range(u){
		if val{
			for key,val:=range(InfoMap){
				dummyvar=val
				dummyvar.U_buttons[i]=false
				InfoMap[key]=dummyvar
			}
		}
	}
	
	for i,val:=range(d){
		if val{
			for key,val:=range(InfoMap){
				dummyvar=val
				dummyvar.D_buttons[i]=false
				InfoMap[key]=dummyvar
			}
		}
	}
}

func spreadOrders(info map[string]gen.ElevButtons){
	u:=[]bool{false,false,false}
	d:=[]bool{false,false,false}
	var temp gen.ElevButtons
	for _,val:=range(info){
		
		for i:=0; i<gen.N_FLOORS-1;i++{
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
		for i:=0; i<gen.N_FLOORS-1;i++{
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
    for {
        msg:=makeMessage("",MyIP,"I'm Alive")
        con:=getUDPcon(BIP,ImAlivePort)
        if con!=nil{
            bmsg:=msgToByte(msg)
            con.Write(bmsg)
        }
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

func sendMsgToMaster(master chan string,sendChan chan gen.ElevButtons, MyIP string,receiveChan chan Message){
    var mst string
    var info gen.ElevButtons
//	var sendTime time.Time
    for{
        select{
        case mst=<-master:
            master<-mst
        case info=<-sendChan:
        	go sendMsg(mst,MyIP,receiveChan,info)
        }
    }
}

func sendMsg(master string,MyIP string, receiveChan chan Message,msg gen.ElevButtons){
	Smsg:=makeMessage("U",MyIP,elevButtonToStr(msg))
    con:=getUDPcon(master,comPORT)
    if con==nil{
    	receiveChan<-Smsg
    	return
    }
    Bmsg:=msgToByte(Smsg)
    con.Write(Bmsg)
}

func sendMsgToThisGuy(IPadrTo string,elevInfo gen.ElevButtons,MyIP string,receiveChan chan Message){
    info:=elevButtonToStr(elevInfo)
	con:=getUDPcon(IPadrTo,comPORT)
	Smsg:=makeMessage("C",MyIP,info)
	if con==nil{
	    if IPadrTo==MyIP{
	        receiveChan<-Smsg
	    }
	    return
	}
	Bmsg:=msgToByte(Smsg)
	con.Write(Bmsg)
}



func recieveMsg(receiveChan chan Message,MyIP string){
    msg:=make(chan Message)
    var Msg Message
    go listenerCon("",comPORT,MyIP,msg)
    for{
        Msg=<-msg
        receiveChan<-Msg
    }
}






