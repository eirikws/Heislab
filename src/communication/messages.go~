package main

import "strings"

type Message struct{
    from string
    to string
    info string
}

func msgToByte(p Message) []byte{
    return []byte(p.from+"£"+p.to+"£"+p.info+"\x00")
}

func byteToMsg(p []byte) Message{
    msg:=strings.Split(string(p[:]),"£")
    return Message{msg[0],msg[1],msg[2]}
}

func MakeMessage(IpAdrFrom string,IpAdrTo string,info string) Message{
    return Message{IpAdrFrom,IpAdrTo,info}
}
