package main


import(
    "fmt"
    //"net"
    ."runtime"
    com "./../communication"
 //   elev "./../elevator"
//    "time"

)

func main(){
    sendMsg:=make(chan string)
    getMsg:=make(chan string)
    var msg string
    go com.Communication(sendMsg,getMsg)
 //   direction :=make(chan elev.CALL_DIRECTION)
    GOMAXPROCS(NumCPU())
    for{
        _,_=fmt.Scanf("%s",&msg)
        sendMsg<-msg
        msg=<-getMsg
        fmt.Println(msg)
    }
}
