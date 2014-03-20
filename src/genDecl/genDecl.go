package genDecl


import(
//	 "strconv"
//	 "strings"
	// "fmt"
)

const N_BUTTONS int=3
const N_FLOORS int=4

type ElevButtons struct{
    U_buttons[3] bool
    D_buttons[3] bool
    C_buttons[4] bool
    Stop_button bool
    Current_floor int
    Obstruction bool
    Door_open bool
    Planned_stops[4] bool
    Dir bool
}





