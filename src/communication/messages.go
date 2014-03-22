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
	str:="ub:"
	for i:=0; i<gen.N_FLOORS-1; i++{
		str=str+strconv.FormatBool(button.U_buttons[i])
		if i<gen.N_FLOORS-2{
			str=str+"."
		}
	}
	str=str+";db:"
	for i:=0; i<gen.N_FLOORS-1; i++{
		str=str+strconv.FormatBool(button.D_buttons[i])
		if i<gen.N_FLOORS-2{
			str=str+"."
		}
	}
	str=str+";cb:"
	for i:=0; i<gen.N_FLOORS; i++{
		str=str+strconv.FormatBool(button.C_buttons[i])
		if i<gen.N_FLOORS-1{
			str=str+"."
		}
	}
	str=str+";sb:"+strconv.FormatBool(button.Stop_button)
	str=str+";cf:"+strconv.Itoa(button.Current_floor)
	str=str+";obs:"+strconv.FormatBool(button.Obstruction)
	str=str+";do:"+strconv.FormatBool(button.Door_open)
	str=str+";ps:"
	for i:=0; i<gen.N_FLOORS; i++{
		str=str+strconv.FormatBool(button.Planned_stops[i])
		if i<gen.N_FLOORS-1{
			str=str+"."
		}
	}
	str=str+";dir:"+strconv.Itoa(button.Dir)+";dummy: "
	return str
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
            myButton.Dir,_=strconv.Atoi(jval)
         }
      }
   }
   return myButton
}

