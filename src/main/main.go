package main


import(
    "fmt"
    //"net"
    ."runtime"
    com "./../communication"
    elev "./../elevator"
    "time"
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
    fmt.Println("after elev_init")
    go elev.Elev_set_speed(direction)
    elev.Elevator_init(direction)
    go elev.Set_lights(buttons)
    go elev.Check_buttons(buttons,msgbuttons)
    go gen.MakeInfoStr(sendMsgToMaster,msgbuttons)
    go elev.Run_elevator(direction,buttons)
    buttons<-button
    for{

    	msg=<-getMsg
    	button=<-buttons
    	buttons<-gen.StringToButton(msg)
    	fmt.Println("wroteTOButtons")
    	time.Sleep(time.Second*0)
    }
}
