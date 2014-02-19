package main


import(
    "fmt"
    //"net"
    ."runtime"
 //   com "./../communication"
    elev "./../elevator"
    "time"

)

const PORT="30001"

func main(){
    direction :=make(chan elev.CALL_DIRECTION)
    GOMAXPROCS(NumCPU())
    elev.Elev_init()
    go elev.Elev_set_speed(direction)
    direction<-elev.CALL_UP
    time.Sleep(time.Millisecond * 2000)
    direction<-elev.CALL_COMMAND
    for{
        fmt.Println("stop")
    }
}
