package communication

import "strings"
import "strconv"
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

func ReadMsg(msg string) (string,ElevButtons){
   return msg[0:15],StringToButton(msg[15:len(msg)-1])
}

func StringToButton(str string) ElevButtons{
   var strArr,strArr2 []string
   var typ string
   var myButton ElevButtons
   strArr=strings.Split(str,";")
   for _,ival:=range(strArr){
      strArr2=strings.Split(ival,":")
      typ=strArr2[0]
      
      for j,jval:=range(strings.Split(strArr2[1],".")){
         if typ=="ub"{
            myButton.u_buttons[j]=(jval=="true")
         } else if typ=="db"{
            myButton.d_buttons[j]=(jval=="true")
         } else if typ=="cb"{
            myButton.c_buttons[j]=(jval=="true")
         } else if typ=="sb"{
            myButton.stop_button=(jval=="true")
         } else if typ=="cf"{
              myButton.current_floor,_=strconv.Atoi(jval)
         } else if typ=="obs"{
            myButton.obstruction=(jval=="true")
         } else if typ=="do"{
            myButton.door_open=(jval=="true")
         }
      }
   }
   return myButton
}
