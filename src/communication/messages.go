package communication

import (
	"strings"
	"strconv"
	gen "./../genDecl"
)

type Message struct{
	toModule string
    from string
    info string
}

func msgToByte(p Message) []byte{
    return []byte(p.toModule+"£"+p.from+"£"+p.info+"\x00")
}

func byteToMsg(p []byte) Message{
    msg:=strings.Split(string(p[:]),"£")
    return Message{msg[0],msg[1],msg[2]}
}

func makeMessage(toModule string, IpAdrFrom string,info string) Message{
    return Message{toModule,IpAdrFrom,info}
}

func elevButtonToStr(button gen.ElevInfo) string{
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

func stringToButton(str string) gen.ElevInfo{
   var strArr,strArr2 []string
   var toModule string
   var myButton gen.ElevInfo
   strArr=strings.Split(str,";")
   for _,ival:=range(strArr){
      strArr2=strings.Split(ival,":")
      toModule=strArr2[0]
      
      for j,jval:=range(strings.Split(strArr2[1],".")){
         if toModule=="ub"{
            myButton.U_buttons[j]=(jval=="true")
         } else if toModule=="db"{
            myButton.D_buttons[j]=(jval=="true")
         } else if toModule=="cb"{
            myButton.C_buttons[j]=(jval=="true")
         } else if toModule=="sb"{
            myButton.Stop_button=(jval=="true")
         } else if toModule=="cf"{
            myButton.Current_floor,_=strconv.Atoi(jval)
         } else if toModule=="obs"{
            myButton.Obstruction=(jval=="true")
         } else if toModule=="do"{
            myButton.Door_open=(jval=="true")
         } else if toModule=="ps"{
            myButton.Planned_stops[j]=(jval=="true")
         } else if toModule=="dir"{
            myButton.Dir,_=strconv.Atoi(jval)
         }
      }
   }
   return myButton
}

