package elevator

import "fmt"
import "time"

const N_BUTTONS int=3
const N_FLOORS int=4


type CALL_DIRECTION int
const (
    CALL_UP CALL_DIRECTION=iota
    CALL_DOWN
    CALL_COMMAND
)

type LAMP_COMMAND int
const (
    LAMP_OFF=iota
    LAMP_ON
)

const (
    DRIVE=3248
    STOPT=2048
    SLEEPTIME=(time.Millisecond * 50)
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
 

func Elev_init() bool{
    fmt.Println(LAMP_OFF)
    if io_init()==0{
        return false
    }
    for i:=0;i<N_FLOORS;i++{
        if i!=0{
            Elev_set_button_lamp(CALL_DOWN,i,LAMP_OFF)
        }
        if i!=N_FLOORS-1{
            Elev_set_button_lamp(CALL_UP,i,LAMP_OFF)
        }
        Elev_set_button_lamp(CALL_COMMAND,i,LAMP_OFF)
    }
    Elev_set_stop_lamp(LAMP_OFF)
    Elev_set_door_open_lamp(LAMP_OFF)
    Elev_set_floor_indicator(LAMP_OFF)
    return true
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

func Elev_get_floor_sensor_signal() int{
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
        fmt.Println("ERROR: elev_get_floor_sensor_signal")
        return -1
    }
}

func Elev_get_button_signal(button CALL_DIRECTION,floor int) bool{
    if io_read_bit(button_channel_matrix[floor][button])==1{
        return true
    }else{return false}
}

func Elev_get_stop_signal()bool{
    return io_read_bit(STOP)!=0
}

func Elev_get_obstruction_signal()bool{
    return io_read_bit(OBSTRUCTION) != 0
}

func Elev_set_floor_indicator(floor int){
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

func Elev_set_button_lamp(button CALL_DIRECTION,floor int, value LAMP_COMMAND){
    if value == LAMP_ON{
        io_set_bit(lamp_channel_matrix[floor][button])
    }
    if value == LAMP_OFF{
        io_clear_bit(lamp_channel_matrix[floor][button])
    }
}

func Elev_set_stop_lamp(value LAMP_COMMAND){
    if value==LAMP_ON{
        io_set_bit(LIGHT_STOP)
    }
    if value==LAMP_OFF{
        io_clear_bit(LIGHT_STOP)
    }else{
        fmt.Println("Error: set_stop_lamp")
    }
}

func Elev_set_door_open_lamp(value LAMP_COMMAND){
    if value==LAMP_ON{
        io_set_bit(DOOR_OPEN)
    }
    if value==LAMP_OFF{
        io_clear_bit(DOOR_OPEN)
    }else{fmt.Println("Error: set_door_open")}
}
