package communication

import (
    mst "./../master"
    "time"
    "fmt"
    "net"
    gen "./../genDecl"
)

type IPandTimeStamp struct{
   IPadr string
   Timestamp time.Time
}

type MsgToGuy struct{
    IPadr string
    msg gen.ElevButtons
}

const ImAlivePort="30108"
const comPORT="30107"

func Communication(sendChanMaster chan string, getChan chan string){
    ch:=make(chan []IPandTimeStamp)
    interMsg:=make(chan Message)
    msgToGuy:=make(chan MsgToGuy,10)
    receiveChan:=make(chan string)
    deletedIP:=make(chan string)
    getElevInfoChan:=make(chan map[string]gen.ElevButtons)
	orders:=make(chan string)
    var IPadr,LastMaster,from,msg string
    var eleButtons gen.ElevButtons
    var AliveList []IPandTimeStamp
    elevInfo:= make(map[string]gen.ElevButtons)
    master:=make(chan string)
    MyIP:=getMyIP()
    BIP:=getBIP(MyIP)
    go SendMsgToThisGuy(interMsg,master,msgToGuy,MyIP)
    go sendImAlive(MyIP,BIP)
    go imAliveListener(MyIP,BIP,ch)
    go sendMsgToMaster(master,sendChanMaster,interMsg,MyIP)
    go recieveMsg(receiveChan,MyIP,interMsg)
    go timeStampCheck(ch,deletedIP,MyIP)
    go mst.Master(master,getElevInfoChan,orders,MyIP)
    for{
        select {
        case AliveList=<-ch:
            AliveList=IPsort(AliveList)
            master<-AliveList[0].IPadr
            if LastMaster!=AliveList[0].IPadr{
                fmt.Println("yep")
            	sendChanMaster<-"U:"+gen.ElevButtonToStr(elevInfo[MyIP])
            }
            LastMaster=AliveList[0].IPadr
            ch<-AliveList
        case IPadr=<-deletedIP:
        	delete(elevInfo,IPadr)
        	fmt.Println("got deleted IP")
        case msg=<-receiveChan:
       		fmt.Println("got msg", msg[15:17])
        	switch {
        	case msg[15:17]=="C:":
        		getChan<-msg[17:]
        	case msg[15:17]=="U:":
            	from,eleButtons=gen.ReadMsg(msg)
            	turnOffLights(elevInfo,from,eleButtons)
            	elevInfo[from]=eleButtons
           		spreadOrders(elevInfo)
           		getElevInfoChan<-elevInfo
           		elevInfo=<-getElevInfoChan
           		for key,val:=range(elevInfo){
           		    fmt.Println("to guy")
           		    msgToGuy<-MsgToGuy{key,val}
           		    time.Sleep(time.Millisecond*5)
           		    fmt.Println("to guy2")
           		}
           	}
        }
    }
}

func turnOffLights(InfoMap map[string]gen.ElevButtons, IPadrFrom string, newInfo gen.ElevButtons){
	u:=[]bool{false,false,false}
	d:=[]bool{false,false,false}
	var dummyvar gen.ElevButtons
	for i:=0 ; i<3 ; i++{
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
    
    for {
        msg:=makeMessage(MyIP,"I'm Alive")
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

func sendMsgToMaster(master,sendChan chan string,interMsg chan Message, MyIP string){
    var mst,msg string
    for{
        select{
        case mst=<-master:
            master<-mst
        case msg=<-sendChan:
            Smsg:=makeMessage(MyIP,msg)
            
            if mst==MyIP{
                interMsg<-Smsg
                continue
            }
            con:=getUDPcon(mst,comPORT)
            Bmsg:=msgToByte(Smsg)
            
            con.Write(Bmsg)
        }
    }
}

func SendMsgToThisGuy(interMsg chan Message,master chan string,msgAndIP chan MsgToGuy,MyIP string){
    var inc MsgToGuy
    var IP,mast string
    var info gen.ElevButtons
    var Smsg Message
    var Bmsg []byte
    var con *net.UDPConn
    for{
        select{
        case mast=<-master:
            //fmt.Println("new Master in sendmsgtotheguy :",mast)
        case inc=<-msgAndIP:
            IP=inc.IPadr
            info=inc.msg
            Smsg=makeMessage(MyIP,"C:"+gen.ElevButtonToStr(info))
            fmt.Println("send msg to guy")
            if mast==MyIP{
                fmt.Println("send msg to guy2222")
                interMsg<-Smsg
                fmt.Println("send msg to guy2")
                continue
            }
            fmt.Println("send msg to guy3")
            Bmsg=msgToByte(Smsg)
            con=getUDPcon(IP,comPORT)
            con.Write(Bmsg)
        }
    }
}

func recieveMsg(getChan chan string,MyIP string,msg chan Message){
    var Msg Message
    go listenerCon("",comPORT,MyIP,msg)
    for{
        fmt.Println("got msg in receive")
        Msg=<-msg
        fmt.Println("receive : write")
        getChan<-Msg.from+Msg.info
    }
}






