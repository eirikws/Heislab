
package elevator
/*
#cgo LDFLAGS: -lcomedi -lm
#include "elev.h"
#include "io.h"

#include <assert.h>
#include <stdlib.h>
*/
import "C"

func elev_init() bool{
    return C.elev_init()==1
}

func elev_set_speed(speed int){
    C.elev_set_speed(speed)
}

func elev_get_floor_sensor_signal() int{
    return C.elev_get_floor_sensor_signal()
}

func elev_get_button_signal(button elev_button_type_t,floor int) bool{
    return C.elev_get_button_signal(button,floor)==1
}

func elev_get_stop_signal()bool{
    return C.elev_get_stop_signal()==1
}

func elev_get_obstruction_signal()bool{
    return C.elev_get_obstruction_signal()==1
}

func elev_set_floor_indicator(floor int){
    C.elev_set_floor_indicator(floor)
}

func elev_set_button_lamp(button elev_button_type_t,floor, value int){
    C.elev_set_button_lamp(button,floor,value)
}

func elev_set_stop_lamp(value int){
    C.elev_set_stop_lamp(value)
}

func elev_set_door_open_lamp(value int){
    C.elev_set_door_open_lamp(value)
}
