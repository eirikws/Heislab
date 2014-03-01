package main


import(
    "fmt"
    //"net"
    ."runtime"
    com "./../communication"
    elev "./../elevator"
//    "time"

)

func main(){
    sendMsg:=make(chan string)
    getMsg:=make(chan string)
    buttons:=make(chan elev.ElevButtons)
  //  var msg string
	var button elev.ElevButtons
    go com.Communication(sendMsg,getMsg)
    direction :=make(chan elev.CALL_DIRECTION)
    GOMAXPROCS(NumCPU())
	



    elev.Init_buttons(&button)

    elev.Elev_init()
    go elev.Elev_set_speed(direction)
    elev.Elevator_init(direction)
	buttons<-button
    fmt.Println("after buttonsandlights")
	go elev.Set_lights(buttons)
	go elev.Check_buttons(buttons)
	go elev.MakeInfoStr(sendMsg,buttons)
	fmt.Println("after infostr")
}
