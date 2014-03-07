package communication

import "strings"
//import "strconv"
//import "fmt"
type Message struct{
    from string
    info string
}

func msgToByte(p Message) []byte{
    return []byte(p.from+"£"+p.info+"\x00")
}

func byteToMsg(p []byte) Message{
    msg:=strings.Split(string(p[:]),"£")
    return Message{msg[0],msg[1]}
}

func makeMessage(IpAdrFrom string,info string) Message{
    return Message{IpAdrFrom,info}
}


