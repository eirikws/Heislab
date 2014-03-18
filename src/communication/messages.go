package communication

import (
	"strings"
	"strconv"
	gen "./../genDecl"
)
//import "fmt"
type Message struct{
	typ string
    from string
    info string
}

func msgToByte(p Message) []byte{
    return []byte(p.typ+"£"+p.from+"£"+p.info+"\x00")
}

func byteToMsg(p []byte) Message{
    msg:=strings.Split(string(p[:]),"£")
    return Message{msg[0],msg[1],msg[2]}
}

func makeMessage(typ string, IpAdrFrom string,info string) Message{
    return Message{typ,IpAdrFrom,info}
}

func elevButtonToStr(button gen.ElevButtons) string{
	return "ub:"+strconv.FormatBool(button.U_buttons[0])+"."+strconv.FormatBool(button.U_buttons[1])+"."+strconv.FormatBool(button.U_buttons[2])+";db:"+strconv.FormatBool(button.D_buttons[0])+"."+strconv.FormatBool(button.D_buttons[1])+"."+strconv.FormatBool(button.D_buttons[2])+";cb:"+strconv.FormatBool(button.C_buttons[0])+"."+strconv.FormatBool(button.C_buttons[1])+"."+strconv.FormatBool(button.C_buttons[2])+"."+strconv.FormatBool(button.C_buttons[3])+";sb:"+strconv.FormatBool(button.Stop_button)+";cf:"+strconv.Itoa(button.Current_floor)+";obs:"+strconv.FormatBool(button.Obstruction)+";do:"+strconv.FormatBool(button.Door_open)+";ps:"+strconv.FormatBool(button.Planned_stops[0])+"."+strconv.FormatBool(button.Planned_stops[1])+"."+strconv.FormatBool(button.Planned_stops[2])+"."+strconv.FormatBool(button.Planned_stops[3])+";dir:"+strconv.FormatBool(button.Dir)+";dummy:  "
}


func stringToButton(str string) gen.ElevButtons{
   var strArr,strArr2 []string
   var typ string
   var myButton gen.ElevButtons
   strArr=strings.Split(str,";")
   for _,ival:=range(strArr){
      strArr2=strings.Split(ival,":")
      typ=strArr2[0]
      
      for j,jval:=range(strings.Split(strArr2[1],".")){
         if typ=="ub"{
            myButton.U_buttons[j]=(jval=="true")
         } else if typ=="db"{
            myButton.D_buttons[j]=(jval=="true")
         } else if typ=="cb"{
            myButton.C_buttons[j]=(jval=="true")
         } else if typ=="sb"{
            myButton.Stop_button=(jval=="true")
         } else if typ=="cf"{
            myButton.Current_floor,_=strconv.Atoi(jval)
         } else if typ=="obs"{
            myButton.Obstruction=(jval=="true")
         } else if typ=="do"{
            myButton.Door_open=(jval=="true")
         } else if typ=="ps"{
            myButton.Planned_stops[j]=(jval=="true")
         } else if typ=="dir"{
            myButton.Dir=(jval=="true")
         }
      }
   }
   return myButton
}

