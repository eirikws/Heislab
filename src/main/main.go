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
    go com.Communication(sendMsg,getMsg)
    direction :=make(chan elev.CALL_DIRECTION)
    GOMAXPROCS(NumCPU())
    
    elev.Elev_init()
    go elev.Elev_set_speed(direction)
    fmt.Println("after elev_set_speed")
    elev.Elevator_init(direction,buttons)
    fmt.Println("after elevator_init")
    go elev.ButtonsAndLights(buttons)
    fmt.Println("after buttonsandlights")
    for{

    }
}
