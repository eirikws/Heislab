package main


import(
    //"fmt"
    //"net"
    ."runtime"
    //"time"
    com "./../communication"
    elev "./../elevator"
    gen "./../genDecl"
)



func main(){
   // x:=0
    sendMsgToMaster:=make(chan gen.ElevButtons)
    getMsg:=make(chan gen.ElevButtons)
    buttons:=make(chan gen.ElevButtons)
    var button gen.ElevButtons
    go com.Communication(sendMsgToMaster,getMsg)
    direction :=make(chan elev.CALL_DIRECTION)
    GOMAXPROCS(NumCPU())
	
    elev.Init_buttons(&button)

    elev.Elev_init()
    go elev.Elev_set_speed(direction)
    elev.Elevator_init(direction)
    go elev.Set_lights(buttons)
    go elev.Check_buttons(buttons,sendMsgToMaster)
    go elev.Run_elevator(direction,buttons,sendMsgToMaster)
    buttons<-button
    for{
        
    	//msg=<-getMsg
    	//button=<-buttons
    	//buttons=<-getMsg
    	//time.Sleep(time.Second*0)
    	
    	button=<-getMsg
    	buttons<-button
    	
    }
}
