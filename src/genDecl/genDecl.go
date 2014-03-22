package genDecl

const N_BUTTONS int=3
const N_FLOORS int=4

type ElevInfo struct{
    U_buttons[N_BUTTONS] bool
    D_buttons[N_BUTTONS] bool
    C_buttons[N_FLOORS] bool
    Stop_button bool
    Current_floor int
    Obstruction bool
    Door_open bool
    Planned_stops[N_FLOORS] bool
    Dir int
}




