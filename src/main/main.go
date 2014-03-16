package main


import(
    //"fmt"
    //"net"
    ."runtime"
    "time"
    com "./../communication"
    elev "./../elevator"
    gen "./../genDecl"
)



func main(){
   // x:=0
    sendMsgToMaster:=make(chan string)
    getMsg:=make(chan string)
    buttons:=make(chan gen.ElevButtons)
    msgbuttons:=make (chan gen.ElevButtons)
    var msg string
    var button gen.ElevButtons
    go com.Communication(sendMsgToMaster,getMsg)
    direction :=make(chan elev.CALL_DIRECTION)
    GOMAXPROCS(NumCPU())
	
    elev.Init_buttons(&button)

    elev.Elev_init()
    go elev.Elev_set_speed(direction)
    elev.Elevator_init(direction)
    go elev.Set_lights(buttons)
    go elev.Check_buttons(buttons,msgbuttons)
    go gen.MakeInfoStr(sendMsgToMaster,msgbuttons)
    go elev.Run_elevator(direction,buttons,msgbuttons)
    buttons<-button
    for{
        
    	msg=<-getMsg
    	button=<-buttons
    	buttons<-gen.StringToButton(msg)
    	time.Sleep(time.Second*0)
    	
    }
}
