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

type ELEVATOR_STATE int
const (
    DRIVE_UP ELEVATOR_STATE=iota
    DRIVE_DOWN
    WAIT
    PITSTOP
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

func Elevator(sendMsgToMaster,getMsg chan ElevInfo){
	direction :=make(chan CALL_DIRECTION)
	elevInfoChan:=make(chan ElevInfo)
	var elevInfo ElevInfo

	elev_init()
    go elev_set_speed(direction)
    go set_lights(elevInfoChan)
    go check_buttons(elevInfoChan,sendMsgToMaster)
    go run_elevator(direction,elevInfoChan,sendMsgToMaster)
    
    elevator_init(direction)
    init_buttons(&elevInfo)
    elevInfoChan<-elevInfo
    for{
    	elevInfo=<-getMsg
    	elevInfoChan<-elevInfo
    }
}

func run_elevator(setElevDir chan CALL_DIRECTION,elevButtonChan,msgToMaster chan ElevInfo){
    var elevInfo ElevInfo
    var doorTime time.Time
    var setDir bool
    var State ELEVATOR_STATE
    var j int
    State=WAIT
    for {
        select{
        case elevInfo=<-elevButtonChan:
        case <-time.After(time.Millisecond*5):
        	j=elev_get_floor_sensor_signal()
        	setDir=false
        	switch State{
				case DRIVE_UP:
					if (j!=-1 && j!=N_FLOORS-1){
						if (elevInfo.Planned_stops[j]) && (elevInfo.U_buttons[j] || elevInfo.C_buttons[j]){
							State=PITSTOP
							break
						}
					}
					//checks if there are stops above
					for i:=elevInfo.Current_floor+1 ; i<N_FLOORS ; i++{
						if elevInfo.Planned_stops[i]{
							setElevDir<-CALL_UP
							elevInfo.Dir=1
							setDir=true
						}
					}
					if !setDir{
						elevInfo.Dir=0
						State=WAIT
						break
					} 
				case DRIVE_DOWN:
					if (j!=-1 && j!=0){
						if (elevInfo.Planned_stops[j]) && (elevInfo.D_buttons[j-1] || elevInfo.C_buttons[j]){
							State=PITSTOP
							break
						}
					}
					//checks if there are stops belove
					for i:=elevInfo.Current_floor-1 ;  i>-1 ; i--{
						 if elevInfo.Planned_stops[i]{
							setElevDir<-CALL_DOWN
							elevInfo.Dir=-1
							setDir=true
						}
					}
					if !setDir{
						elevInfo.Dir=0
						State=WAIT
						break
					}
				case WAIT:
					if j==-1{
						setElevDir<-CALL_UP
						elevInfo.Dir=1
						break
					} else {setElevDir<-CALL_NEUTRAL}
					if elevInfo.Planned_stops[j]{
						State=PITSTOP
						break
					} else if isAbove(elevInfo.Current_floor,elevInfo.Planned_stops)==1{
						elevInfo.Dir=1
						State=DRIVE_UP
						break
					} else if isBelove(elevInfo.Current_floor,elevInfo.Planned_stops)==1{
						elevInfo.Dir=-1
						State=DRIVE_DOWN
						break
					}
				case PITSTOP:
					if elevInfo.Door_open{
						if doorTime.Before(time.Now()) || elev_get_floor_sensor_signal()==-1{
							elevInfo.Door_open=false
							msgToMaster<-elevInfo
							if elevInfo.Dir==1{
								State=DRIVE_UP
								break
							} else if elevInfo.Dir==-1{
								State=DRIVE_DOWN
								break
							} else if elevInfo.Dir==0{
								State=WAIT
								break
							}
						}
					} else{
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
						msgToMaster<-elevInfo
					}    	
				}
				elevButtonChan<-elevInfo
			
			
        }
    }
}

//checks if there is a planned stop above the elevator
func isAbove(myFloor int,Planned_stops [N_FLOORS]bool)int{
	for i:=myFloor+1; i<N_FLOORS;i++{
		if Planned_stops[i]{
			return 1
		}
	}
	return 0
}

//checks if there is a planned stop belove the elevator
func isBelove(myFloor int,Planned_stops [N_FLOORS]bool)int{
	for i:=myFloor-1; i>-1;i--{
		if Planned_stops[i]{
			return 1
		}
	}
	return 0
}

//updates the elevInfo struct when a button is pushed. Sends a message to coms when a new button is pushed
//or over 5 seconds has passed since it sent one
func check_buttons(elevInfoChan chan ElevInfo,sendToMaster chan ElevInfo) bool{
	var elevInfo ElevInfo
	var x,i int
	var sendTime time.Time
	sendTime=time.Now().Add(5*time.Second)
	for{
		x=0
		elevInfo=<-elevInfoChan
		for i:=0; i<N_FLOORS-1; i++ {
			if elev_get_button_signal(CALL_UP, i){
			   if elevInfo.U_buttons[i]==false{
			      x=1
			   }
				elevInfo.U_buttons[i]=true
			}
			if elev_get_button_signal(CALL_DOWN, i+1){
			   if elevInfo.D_buttons[i]==false{
			      x=1
			   }
				elevInfo.D_buttons[i]=true
			}
			if elev_get_button_signal(CALL_NEUTRAL, i){
			   if elevInfo.C_buttons[i]==false{
			      x=1
			   }
				elevInfo.C_buttons[i]=true
			}
		}
		if elev_get_button_signal(CALL_NEUTRAL,N_FLOORS-1){
		   if elevInfo.C_buttons[N_FLOORS-1]==false{
			      x=1
			   }
			elevInfo.C_buttons[N_FLOORS-1]=true
		}

		if elev_get_stop_signal(){
		   if elevInfo.Stop_button==false{
			      x=1
			   }
			elevInfo.Stop_button=true
		}
		i=elev_get_floor_sensor_signal()
		if i!=-1{
		    if elevInfo.Current_floor!=i{
		       x=1
		    }
		    elevInfo.Current_floor=i
		}
		if x==1 || sendTime.Before(time.Now()){
		   sendToMaster<-elevInfo
		   sendTime=time.Now().Add(5*time.Second)
		}
		elevInfoChan<-elevInfo
		time.Sleep(50)
	}
}

func set_lights(elevInfoChan chan ElevInfo){
	var elevInfo ElevInfo
	for{
		elevInfo=<-elevInfoChan
    	for i:=0; i<N_FLOORS-1; i++{
			elev_set_button_lamp(CALL_NEUTRAL,i,elevInfo.C_buttons[i])
			elev_set_button_lamp(CALL_UP,i,elevInfo.U_buttons[i])
			elev_set_button_lamp(CALL_DOWN,i+1, elevInfo.D_buttons[i])
		}
		elev_set_button_lamp(CALL_NEUTRAL,N_FLOORS-1,elevInfo.C_buttons[N_FLOORS-1])
		elev_set_stop_lamp(elevInfo.Stop_button)
		elev_set_door_open_lamp(elevInfo.Door_open)
		elev_set_floor_indicator(elevInfo.Current_floor)
		elevInfoChan<-elevInfo
	}
}

func elev_init() bool{
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

func init_buttons(buttons *ElevInfo){
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

func elevator_init(drive chan CALL_DIRECTION){
    drive<-CALL_UP
	
    for elev_get_floor_sensor_signal()==-1{
    
    }
    drive<-CALL_NEUTRAL
}

func elev_set_speed(myDir chan CALL_DIRECTION){
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
