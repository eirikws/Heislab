package elevator

import "fmt"
import "time"
import "strconv"

const N_BUTTONS int=3
const N_FLOORS int=4


type CALL_DIRECTION int
const (
    CALL_UP CALL_DIRECTION=iota
    CALL_DOWN
    CALL_COMMAND
)

const (
    DRIVE=3248
    STOPT=2048
    SLEEPTIME=(time.Millisecond * 15)
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

type ElevButtons struct{
    u_buttons[3] bool
    d_buttons[3] bool
    c_buttons[4] bool
    stop_button bool
    current_floor int
    obstruction bool
    door_open bool
}

func MakeInfoStr(sendMsg chan string,buttons chan ElevButtons){
	var button ElevButtons
	var str string
	for{
		
		button=<-buttons
		str="ub:"+strconv.FormatBool(button.u_buttons[0])+"."+strconv.FormatBool(button.u_buttons[1])+"."+strconv.FormatBool(button.u_buttons[2])+";db:"+strconv.FormatBool(button.d_buttons[0])+"."+strconv.FormatBool(button.d_buttons[1])+"."+strconv.FormatBool(button.d_buttons[2])+";cb:"+strconv.FormatBool(button.c_buttons[0])+"."+strconv.FormatBool(button.c_buttons[1])+"."+strconv.FormatBool(button.c_buttons[2])+"."+strconv.FormatBool(button.c_buttons[3])+";sb:"+strconv.FormatBool(button.stop_button)+";cf:"+strconv.Itoa(button.current_floor)+";obs:"+strconv.FormatBool(button.obstruction)+";do:"+strconv.FormatBool(button.door_open)
		buttons<-button
		sendMsg<-str
		fmt.Println("making str")
	}
}

/*
func ButtonsAndLights(buttons chan ElevButtons){
	var butt ElevButtons
	for{
		
		fmt.Println("in butsAndLightsYYEYEYE")
		fmt.Println("read buttons")
		if check_buttons(&butt){
			fmt.Println("wrote to buttons")
			buttons<-butt
			fmt.Println("muhhaha")
		}
		set_lights(&butt,butt.current_floor)
		
    }
}
*/
func Check_buttons(buttons chan ElevButtons) bool{
	var elbut ElevButtons
	for{
		elbut=<-buttons
		for i:=0; i<N_FLOORS-1; i++ {
			if elev_get_button_signal(CALL_UP, i){
				elbut.u_buttons[i]=true
			}
			if elev_get_button_signal(CALL_DOWN, i+1){
				elbut.d_buttons[i]=true
			}
			if elev_get_button_signal(CALL_COMMAND, i){
				elbut.c_buttons[i]=true
			}
		}
		if elev_get_button_signal(CALL_COMMAND,3){
			elbut.c_buttons[3]=true
		}
		if elev_get_stop_signal(){
			elbut.stop_button=true
		}
		i:=elev_get_floor_sensor_signal()
		if i!=-1{
		    elbut.current_floor=i
		}
		buttons<-elbut
	}
}

func Set_lights(buttons chan ElevButtons){
	var button ElevButtons
	for{
		button=<-buttons
    	for i:=0; i<N_FLOORS-1; i++{
			elev_set_button_lamp(CALL_COMMAND,i,button.c_buttons[i])
			elev_set_button_lamp(CALL_UP,i,button.u_buttons[i])
			elev_set_button_lamp(CALL_DOWN,i+1, button.d_buttons[i])
		}
		elev_set_button_lamp(CALL_COMMAND,3,button.c_buttons[3])
		elev_set_stop_lamp(button.stop_button)
		//set the floor_indicators
		elev_set_floor_indicator(button.current_floor)
		buttons<-button
	}
}


func Elev_init() bool{
    fmt.Println("start elev_init")
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
        elev_set_button_lamp(CALL_COMMAND,i,false)
    }
    elev_set_stop_lamp(false)
    elev_set_door_open_lamp(false)
    fmt.Println("exit elev_init")
    return true
}

func Init_buttons(buttons *ElevButtons){
    for i:=0;i<N_FLOORS-1;i++{
        buttons.u_buttons[i]=false
        buttons.d_buttons[i]=false
        buttons.c_buttons[i]=false
    }
    buttons.c_buttons[3]=false
    buttons.stop_button=false
    buttons.door_open=false
    buttons.obstruction=false
	buttons.current_floor=-1
}

func Elevator_init(drive chan CALL_DIRECTION){
    drive<-CALL_UP
	
    for elev_get_floor_sensor_signal()==-1{
    
    }
    drive<-CALL_COMMAND
}

func Elev_set_speed(myDir chan CALL_DIRECTION){
    lastDir:=CALL_COMMAND
    nowDir:=CALL_COMMAND
    for{
        nowDir=<-myDir
        switch nowDir{
            case CALL_UP:
            io_clear_bit(MOTORDIR)
            io_write_analog(MOTOR,DRIVE)
            case CALL_DOWN:
            io_set_bit(MOTORDIR)
            io_write_analog(MOTOR,DRIVE)
            case CALL_COMMAND:
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
