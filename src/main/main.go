package main


import(
    "fmt"
    //"net"
    ."runtime"
    com "./../communication"
    elev "./../elevator"
 //   "time"

)

func main(){
   // x:=0
    sendMsg:=make(chan string)
    getMsg:=make(chan string)
    buttons:=make(chan elev.ElevButtons)
    msgbuttons:=make (chan elev.ElevButtons)
    var msg string
    var button elev.ElevButtons
    go com.Communication(sendMsg,getMsg)
    direction :=make(chan elev.CALL_DIRECTION)
    GOMAXPROCS(NumCPU())
	



    elev.Init_buttons(&button)

    elev.Elev_init()
    fmt.Println("after elev_init")
    go elev.Elev_set_speed(direction)
    elev.Elevator_init(direction)
    go elev.Set_lights(buttons)
    go elev.Check_buttons(buttons,msgbuttons)
    go elev.MakeInfoStr(sendMsg,msgbuttons)
    buttons<-button
    for{
        select{
        case msg=<-getMsg:
        fmt.Println(msg)
        }
        //time.Sleep(time.Second)
    }
}
