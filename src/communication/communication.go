package communication

import (
    "time"
    gen "./../genDecl"
)

const ImAlivePort="30108"
const comPORT="30107"

func Communication(sendChanMaster, sendChanElevator chan gen.ElevInfo,master chan string, getElevInfoChan chan map[string]gen.ElevInfo){
    aliveListChan:=make(chan map[string]time.Time)
    receiveChan:=make(chan Message)
    deletedIP:=make(chan string)
    elevInfoMap:= make(map[string]gen.ElevInfo)
    MyIP:=getMyIP()
    BIP:=getBIP(MyIP)
    
    var IPadr,lastMaster,from string
    var msg Message
    var elevInfo gen.ElevInfo
    var AliveList map[string]time.Time
    
    go sendImAlive(MyIP,BIP)
    go imAliveListener(MyIP,BIP,aliveListChan)
    go sendMsgToMaster(master,sendChanMaster,MyIP,receiveChan)
    go recieveMsg(receiveChan,MyIP)
    go timeStampCheck(aliveListChan,deletedIP,MyIP)
    for{
        select {
        case AliveList=<-aliveListChan:
        	master<-getSmallestIP(AliveList)
        	//if there is a new master, send the information struct to it
        	if lastMaster!=getSmallestIP(AliveList){
            	sendChanMaster<-elevInfoMap[MyIP]
            }
            lastMaster=getSmallestIP(AliveList)
            aliveListChan<-AliveList
        case IPadr=<-deletedIP:
        	delete(elevInfoMap,IPadr)
        case msg=<-receiveChan:
        	switch {
        	case msg.toModule=="To Elevator":
        		sendChanElevator<-stringToButton(msg.info)
        		elevInfoMap[MyIP]=stringToButton(msg.info)
        	case msg.toModule=="To Master":
            	from, elevInfo=msg.from, stringToButton(msg.info)
            	turnOffLights(elevInfoMap,from,elevInfo)
            	elevInfoMap[from]=elevInfo
           		getElevInfoChan<-elevInfoMap
           		elevInfoMap=<-getElevInfoChan
           		for key,val:=range(elevInfoMap){
           		    go sendMsgToThisElevator(key,val,MyIP,receiveChan)
           		}
           	
           	}
        }
    }
}

//finds up and down buttons that "from" has turned off, and sets the corresponding buttons in the other elevators false
func turnOffLights(InfoMap map[string]gen.ElevInfo, IPadrFrom string, newInfo gen.ElevInfo){
	buttonsDownDone:=[gen.N_BUTTONS]bool{false}
	buttonsUpDone:=[gen.N_BUTTONS]bool{false}
	var elevatorVar gen.ElevInfo
	//stores any newly turned off buttons in the dummylists
	for i:=0 ; i<gen.N_FLOORS-1 ; i++{
		if InfoMap[IPadrFrom].U_buttons[i] && !newInfo.U_buttons[i]{
			buttonsUpDone[i]=true
		}
		if InfoMap[IPadrFrom].D_buttons[i] && !newInfo.D_buttons[i]{
			buttonsDownDone[i]=true
		}
	}
	//sets a light to false for every elevator if 
	//the dummylists have that light as true
	for i,val:=range(buttonsUpDone){
		if val{
			for key,val:=range(InfoMap){
				elevatorVar=val
				elevatorVar.U_buttons[i]=false
				InfoMap[key]=elevatorVar
			}
		}
	}
	
	for i,val:=range(buttonsDownDone){
		if val{
			for key,val:=range(InfoMap){
				elevatorVar=val
				elevatorVar.D_buttons[i]=false
				InfoMap[key]=elevatorVar
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

func imAliveListener(MyIP, BIP string, aliveMapChan chan map[string]time.Time){
    alivechan:=make(chan Message)
    go listenerCon(BIP,ImAlivePort, MyIP, alivechan)
    var newMsg Message
    IPlist:=make(map[string]time.Time)
    IPlist[MyIP]=time.Now()
    for{
        newMsg=<-alivechan
        IPlist[newMsg.from]=time.Now().Add(1200*time.Millisecond)
        aliveMapChan<-IPlist
        IPlist=<-aliveMapChan
    }
}

func sendMsgToMaster(master chan string,sendChan chan gen.ElevInfo, MyIP string,receiveChan chan Message){
    var mst string
    var info gen.ElevInfo
    for{
        select{
        case mst=<-master:
            master<-mst
        case info=<-sendChan:
        	go sendMsg(mst,MyIP,receiveChan,info)
        }
    }
}

func sendMsg(master string,MyIP string, receiveChan chan Message,msg gen.ElevInfo){
	Smsg:=makeMessage("To Master",MyIP,elevButtonToStr(msg))
    con:=getUDPcon(master,comPORT)
    if con==nil{
    	receiveChan<-Smsg
    	return
    }
    Bmsg:=msgToByte(Smsg)
    con.Write(Bmsg)
}

func sendMsgToThisElevator(IPadrTo string,elevInfo gen.ElevInfo,MyIP string,receiveChan chan Message){
    info:=elevButtonToStr(elevInfo)
	con:=getUDPcon(IPadrTo,comPORT)
	Smsg:=makeMessage("To Elevator",MyIP,info)
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






