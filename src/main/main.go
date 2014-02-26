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
    var msg string
    go com.Communication(sendMsg,getMsg)
    direction :=make(chan elev.CALL_DIRECTION)
    GOMAXPROCS(NumCPU())
    for{
        _,_=fmt.Scanf("%s",&msg)
        sendMsg<-msg
        msg=<-getMsg
        fmt.Println(msg)
    }
    elev.Elev_init()
    go elev.Elev_set_speed(direction)
    elev.Elevator_init(direction)
    var elevator_buttons elev.ElevButtons
    elev.Init_buttons(&elevator_buttons)
    for{
        elev.Check_buttons(&elevator_buttons)
        elev.Set_lights(&elevator_buttons,3)
    
    }
}
