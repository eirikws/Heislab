package elevator

import(
	"fmt"
	"time"
	."./../genDecl"
)


type CALL_DIRECTION int
const (
    CALL_UP CALL_DIRECTION=iota
    CALL_DOWN
    CALL_NEUTRAL
)

const (
    DRIVE=3248
    STOPT=2048
    SLEEPTIME=(time.Millisecond * 2)
)

var lamp_channel_matrix =[N_FLOORS][N_BUTTONS]int{
    {LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
    {LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
    {LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
    {LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var button_channel_matrix =[N_FLOORS][N_BUTTONS]int{
    {FLOOR_UP1, FLOOR_DOWN1, FLOOR_COMMAND1},
    {FLOOR_UP2, FLOOR_DOWN2, FLOOR_COMMAND2},
    {FLOOR_UP3, FLOOR_DOWN3, FLOOR_COMMAND3},
    {FLOOR_UP4, FLOOR_DOWN4, FLOOR_COMMAND4},
}


func Run_elevator(setElevDir chan CALL_DIRECTION,elevButtons,msgButtons chan ElevButtons){
    var elevInfo ElevButtons
    var newStop,numPlannedStops int
    var doorTime time.Time
    var canStop,setDir bool
    for {
        select{
        case elevInfo=<-elevButtons:
        case <-time.After(time.Millisecond*10):
            numPlannedStops=0
            for _,val:= range(elevInfo.Planned_stops){
            	if val{
            		numPlannedStops++
            	}
            }
        	if elevInfo.Door_open{
            	newStop=0
        		if doorTime.Before(time.Now()) || elev_get_floor_sensor_signal()==-1{
        			elevInfo.Door_open=false
        	//		elevInfo.Planned_stops[j]=false
        		}
        	}
        	canStop=false
        	setDir=false
            j:=elev_get_floor_sensor_signal()
            if elevInfo.Dir>-1 && (j!=-1 && j!=N_FLOORS-1) {
            	if elevInfo.U_buttons[j]{
            		fmt.Println("U")
            		canStop=true
            	}
            }
            if elevInfo.Dir<1 && (j!=-1 && j!=0) {
            	
            	if elevInfo.D_buttons[j-1]{
            		fmt.Println("D")
            		canStop=true
            	}
            }
            if j!=-1{
            	if elevInfo.C_buttons[j]{
            		fmt.Println("C")
            		canStop=true
            	}
            }
            if (j!=-1 && elevInfo.Planned_stops[j]) && (canStop){
            	fmt.Println("STOP")
                newStop=1
                setElevDir<-CALL_NEUTRAL
                elevInfo.Door_open=true
                doorTime=time.Now().Add(3*time.Second)
                elevInfo.C_buttons[j]=false
                elevInfo.Planned_stops[j]=false
                if (elevInfo.Dir>-1 && j!=N_FLOORS-1) || j==0{
                    elevInfo.U_buttons[j]=false
                }
                if  (elevInfo.Dir<1 && j!=0) || j==N_FLOORS-1{
                    elevInfo.D_buttons[j-1]=false
                }
            } else if (j==-1 && elevInfo.Planned_stops[elevInfo.Current_floor]) && numPlannedStops==1{
            	setElevDir<-CALL_UP
            } else if elevInfo.Dir==1 && !elevInfo.Door_open {
            	for i:=elevInfo.Current_floor+1 ; i<N_FLOORS ; i++{
                    if elevInfo.Planned_stops[i]{
                        setElevDir<-CALL_UP
                        elevInfo.Dir=1
                        setDir=true
                    }
                }
                if !setDir{
                    elevInfo.Dir=0
                } 
            } else if elevInfo.Dir==-1 && !elevInfo.Door_open {
                for i:=elevInfo.Current_floor-1 ;  i>-1 ; i--{
                    if elevInfo.Planned_stops[i]{
                        setElevDir<-CALL_DOWN
                        elevInfo.Dir=-1
                        setDir=true
                    }
                }
                if !setDir{
                    elevInfo.Dir=0
                }
            }else if elevInfo.Dir==0{
            	if isAbove(elevInfo.Current_floor,elevInfo.Planned_stops)==1{
            		elevInfo.Dir=1
            		fmt.Println("switch to up")
            	} else if isBelove(elevInfo.Current_floor,elevInfo.Planned_stops)==1{
            		elevInfo.Dir=-1
            		fmt.Println("switch to down")
            	}
            }
            elevButtons<-elevInfo
            if newStop==1{
                msgButtons<-elevInfo
            }
        }
    }
}

func isAbove(myFloor int,Planned_stops [N_FLOORS]bool)int{
	for i:=myFloor+1; i<N_FLOORS;i++{
		if Planned_stops[i]{
			return 1
		}
	}
	return 0
}

func isBelove(myFloor int,Planned_stops [N_FLOORS]bool)int{
	for i:=myFloor-1; i>-1;i--{
		if Planned_stops[i]{
			return 1
		}
	}
	return 0
}



func Check_buttons(buttons chan ElevButtons,msgbuttons chan ElevButtons) bool{
	var elbut ElevButtons
	var x,i int
	//var i int
	for{
		x=0
		elbut=<-buttons
		for i:=0; i<N_FLOORS-1; i++ {
			if elev_get_button_signal(CALL_UP, i){
			   if elbut.U_buttons[i]==false{
			      x=1
			   }
				elbut.U_buttons[i]=true
			}
			if elev_get_button_signal(CALL_DOWN, i+1){
			   if elbut.D_buttons[i]==false{
			      x=1
			   }
				elbut.D_buttons[i]=true
			}
			if elev_get_button_signal(CALL_NEUTRAL, i){
			   if elbut.C_buttons[i]==false{
			      x=1
			   }
				elbut.C_buttons[i]=true
			}
		}
		if elev_get_button_signal(CALL_NEUTRAL,N_FLOORS-1){
		   if elbut.C_buttons[N_FLOORS-1]==false{
			      x=1
			   }
			elbut.C_buttons[N_FLOORS-1]=true
		}

		if elev_get_stop_signal(){
		   if elbut.Stop_button==false{
			      x=1
			   }
			elbut.Stop_button=true
		}
		i=elev_get_floor_sensor_signal()
		if i!=-1{
		    if elbut.Current_floor!=i{
		       x=1
		    }
		    elbut.Current_floor=i
		}
		if x==1{
		   msgbuttons<-elbut
		}
		buttons<-elbut
	}
}

/*
func MakeInfoStr(sendMsgTo chan string,msgbuttons chan ElevButtons){
	var button ElevButtons
	var str string
	for{
		button=<-msgbuttons
		str=ElevButtonToStr(button)
		sendMsgTo<-str
	}
}
*/
func Set_lights(buttons chan ElevButtons){
	var button ElevButtons
	for{
		button=<-buttons
    	for i:=0; i<N_FLOORS-1; i++{
			elev_set_button_lamp(CALL_NEUTRAL,i,button.C_buttons[i])
			elev_set_button_lamp(CALL_UP,i,button.U_buttons[i])
			elev_set_button_lamp(CALL_DOWN,i+1, button.D_buttons[i])
		}
		elev_set_button_lamp(CALL_NEUTRAL,N_FLOORS-1,button.C_buttons[N_FLOORS-1])
		elev_set_stop_lamp(button.Stop_button)
		elev_set_door_open_lamp(button.Door_open)
		elev_set_floor_indicator(button.Current_floor)
		buttons<-button
	}
}

func Elev_init() bool{
    if io_init()==0{
        return false
    }
    for i:=0;i<N_FLOORS;i++{
        if i!=0{
            elev_set_button_lamp(CALL_DOWN,i,false)
        }
        if i!=N_FLOORS-1{
            elev_set_button_lamp(CALL_UP,i,false)
        }
        elev_set_button_lamp(CALL_NEUTRAL,i,false)
    }
    elev_set_stop_lamp(false)
    elev_set_door_open_lamp(false)
    return true
}

func Init_buttons(buttons *ElevButtons){
    for i:=0;i<N_FLOORS-1;i++{
        buttons.U_buttons[i]=false
        buttons.D_buttons[i]=false
        buttons.C_buttons[i]=false
    }
    buttons.C_buttons[N_FLOORS-1]=false
    buttons.Stop_button=false
    buttons.Door_open=false
    buttons.Obstruction=false
	buttons.Current_floor=-1
	buttons.Dir=0
}

func Elevator_init(drive chan CALL_DIRECTION){
    drive<-CALL_UP
	
    for elev_get_floor_sensor_signal()==-1{
    
    }
    drive<-CALL_NEUTRAL
}

func Elev_set_speed(myDir chan CALL_DIRECTION){
    lastDir:=CALL_NEUTRAL
    nowDir:=CALL_NEUTRAL
    for{
        nowDir=<-myDir
        switch nowDir{
            case CALL_UP:
            io_clear_bit(MOTORDIR)
            io_write_analog(MOTOR,DRIVE)
            case CALL_DOWN:
            io_set_bit(MOTORDIR)
            io_write_analog(MOTOR,DRIVE)
            case CALL_NEUTRAL:
            if lastDir==CALL_UP{
                io_set_bit(MOTORDIR)
                time.Sleep(SLEEPTIME)
                io_write_analog(MOTOR,STOPT)
            }
            if lastDir==CALL_DOWN{
                io_clear_bit(MOTORDIR)
                time.Sleep(SLEEPTIME)
                io_write_analog(MOTOR,STOPT)
            }
            
        }
        lastDir=nowDir
    }
}

func elev_get_floor_sensor_signal() int{
    switch{
        case io_read_bit(SENSOR1)!=0:
        return 0
        case io_read_bit(SENSOR2)!=0:
        return 1
        case io_read_bit(SENSOR3)!=0:
        return 2
        case io_read_bit(SENSOR4)!=0:
        return 3
        default:
        return -1
    }
}

func elev_get_button_signal(button CALL_DIRECTION,floor int) bool{
    if io_read_bit(button_channel_matrix[floor][button])==1{
        return true
    }else{return false}
}



func elev_get_stop_signal()bool{
    return io_read_bit(STOP)!=0
}

func elev_get_obstruction_signal()bool{
    return io_read_bit(OBSTRUCTION) != 0
}

func elev_set_floor_indicator(floor int){
    switch floor{
        case 0:
        io_clear_bit(FLOOR_IND1)
        io_clear_bit(FLOOR_IND2)
        case 1:
        io_clear_bit(FLOOR_IND1)
        io_set_bit(FLOOR_IND2)
        case 2:
        io_set_bit(FLOOR_IND1)
        io_clear_bit(FLOOR_IND2)
        case 3:
        io_set_bit(FLOOR_IND1)
        io_set_bit(FLOOR_IND2)
    }
}

func elev_set_button_lamp(button CALL_DIRECTION,floor int, value bool){
    if value == true{
        io_set_bit(lamp_channel_matrix[floor][button])
    }
    if value == false{
        io_clear_bit(lamp_channel_matrix[floor][button])
    }
}

func elev_set_stop_lamp(value bool){
    if value==true{
        io_set_bit(LIGHT_STOP)
    }else if value==false{
        io_clear_bit(LIGHT_STOP)
    }else{
        fmt.Println("Error: set_stop_lamp")
    }
}

func elev_set_door_open_lamp(value bool){
    if value==true{
        io_set_bit(DOOR_OPEN)
    }else if value==false{
        io_clear_bit(DOOR_OPEN)
    }else{fmt.Println("Error: set_door_open")}
}
