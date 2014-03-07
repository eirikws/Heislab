package genDecl


import(
	 "strconv"
	 "strings"
//	 "fmt"
)


type ElevButtons struct{
    U_buttons[3] bool
    D_buttons[3] bool
    C_buttons[4] bool
    Stop_button bool
    Current_floor int
    Obstruction bool
    Door_open bool
}

func MakeInfoStr(sendMsgTo chan string,msgbuttons chan ElevButtons){
	var button ElevButtons
	var str string
	for{
		button=<-msgbuttons
		str="U:"+ElevButtonToStr(button)
		sendMsgTo<-str
	}
}

func ElevButtonToStr(button ElevButtons) string{
	return "ub:"+strconv.FormatBool(button.U_buttons[0])+"."+strconv.FormatBool(button.U_buttons[1])+"."+strconv.FormatBool(button.U_buttons[2])+";db:"+strconv.FormatBool(button.D_buttons[0])+"."+strconv.FormatBool(button.D_buttons[1])+"."+strconv.FormatBool(button.D_buttons[2])+";cb:"+strconv.FormatBool(button.C_buttons[0])+"."+strconv.FormatBool(button.C_buttons[1])+"."+strconv.FormatBool(button.C_buttons[2])+"."+strconv.FormatBool(button.C_buttons[3])+";sb:"+strconv.FormatBool(button.Stop_button)+";cf:"+strconv.Itoa(button.Current_floor)+";obs:"+strconv.FormatBool(button.Obstruction)+";do:"+strconv.FormatBool(button.Door_open)

}

func ReadMsg(msg string) (string,ElevButtons){
   return msg[0:15],StringToButton(msg[17:len(msg)-1])
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
         }
      }
   }
   return myButton
}




