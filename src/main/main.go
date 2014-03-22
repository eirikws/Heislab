package main


import(
    ."runtime"
    "time"
    mst "./../master"
    coms "./../communication"
    elev "./../elevator"
    gen "./../genDecl"
)



func main(){
	GOMAXPROCS(NumCPU())

    sendMsgToMaster:=make(chan gen.ElevInfo)
    getMsg:=make(chan gen.ElevInfo)
    
    elevInfoChan:=make(chan map[string]gen.ElevInfo)
    master:=make(chan string)
    
    go coms.Communication(sendMsgToMaster,getMsg,master,elevInfoChan )
    go mst.Master(master,elevInfoChan)
    go elev.Elevator(sendMsgToMaster,getMsg)
    
    
    for{
    	time.Sleep(time.Second)
    }
}
