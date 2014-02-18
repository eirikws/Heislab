
package elevator
/*
#cgo LDFLAGS: -lcomedi -lm
#include "elev.h"
#include "io.h"

#include <assert.h>
#include <stdlib.h>
*/
import "C"
import "fmt"

const N_BUTTONS int=3
const N_FLOORS int=4


type tag_elev_lamp_type int
const (
    BUTTON_CALL_UP tag_elev_lamp_type=iota
    BUTTON_CALL_DOWN
    BUTTON_COMMAND
)

type lamp_command int
const (
    LAMP_ON=iota
    LAMP_OFF
)

const lamp_channel_matrix int=[N_FLOORS][N_BUTTONS]{
    {LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
    {LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
    {LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
    {LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

const button_channel_matrix int=[N_FLOORS][N_BUTTONS]{
    {FLOOR_UP1, FLOOR_DOWN1, FLOOR_COMMAND1},
    {FLOOR_UP2, FLOOR_DOWN2, FLOOR_COMMAND2},
    {FLOOR_UP3, FLOOR_DOWN3, FLOOR_COMMAND3},
    {FLOOR_UP4, FLOOR_DOWN4, FLOOR_COMMAND4},
}
 

func elev_init() bool{
    if !C.io_init(){
        fmt.Println("Error init")
        return 0
    }

}

func elev_set_speed(speed int){
    C.elev_set_speed(speed)
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
        fmt.Println("ERROR: elev_get_floor_sensor_signal")
        return -1
    }


}

func elev_get_button_signal(button elev_button_type_t,floor int) bool{
    if io_read_bit(button_channel_matrix[floor][button]){
        return true
   }
    else{return false}
}

func elev_get_stop_signal()bool{
    return elev_get_stop_signal(STOP) !=0
}

func elev_get_obstruction_signal()bool{
    return io_read_bit(OBSTRUCTION) != 0
}

func elev_set_floor_indicator(floor int){
    if floor & 0x02{
        io_set_bit(FLOOR_IND1)
    }else{io_clear_bit(FLOOR_IND1)}
    if floor & 0x01{
        io_set_bit(FLOOR_IND2)
    }else{io_clear_bit(FLOOR_IND2)}
}

func elev_set_button_lamp(button elev_button_type_t,floor int, value lamp_command){
    if value == LAMP_ON{
        io_set_bit(lamp_channel_matrix[floor][button])
    }if value == LAMP_OFF{
        io_clear_bit(lamp_channel_matrix[floor][button])
}

func elev_set_stop_lamp(value lamp_command){
    if value==LAMP_ON{
        io_set_bit(LIGHT_STOP)
    }if value==LAMP_OFF{
        io_clear_bit(LIGHT_STOP)
    }else{fmt.Println("Error: set_stop_lamp")}
}


func elev_set_door_open_lamp(value lamp_command){
    if value==LAMP_ON{
        io_set_bit(DOOR_OPEN)
    }if value==LAMP_OFF{
        io_clear_bit(DOOR_OPEN)
    }else{fmt.Println("Error: set_door_open")}
}
