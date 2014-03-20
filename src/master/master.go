package master

import(
   // "fmt"
    gen "./../genDecl"
 //   "net"
//    "strings"
)


func Master(master chan string,elevInfoChan chan map[string]gen.ElevButtons,orders chan string,MyIP string){
    var mst,lowestCostIP string
    var u,d [gen.N_FLOORS-1]bool
    elevInfoMap:=make(map[string]gen.ElevButtons)
    costMap:=make(map[string]int)
    var dummycost int
    var dummyElevInfo gen.ElevButtons
    for{
        select{
            case mst=<-master:
            case elevInfoMap=<-elevInfoChan:
                if mst!=MyIP{
                	elevInfoChan<-elevInfoMap
                    continue
                }
                u=elevInfoMap[MyIP].U_buttons
                d=elevInfoMap[MyIP].D_buttons
                
                for key,val:=range(elevInfoMap){
                	dummyElevInfo=val
                	dummyElevInfo.Planned_stops=[gen.N_FLOORS]bool{false}
                	elevInfoMap[key]=dummyElevInfo
                }
                
                
                for key,val := range(elevInfoMap){
                    dummyElevInfo=val
                    for i,ival:=range(val.C_buttons){
                        if ival{
                            dummyElevInfo.Planned_stops[i]=true
                        }
					}
					elevInfoMap[key]=dummyElevInfo
                }
				for floor,order :=range(u){
					dummycost=1000000
					if order{
						for key,val:=range(elevInfoMap){
                			costMap[key]=costFunc(val,true,floor)
                		}
                		for costkey,cost:=range(costMap){
                			if cost<dummycost{
                				dummycost=cost
                				lowestCostIP=costkey
                			}
                		}
                		dummyElevInfo=elevInfoMap[lowestCostIP]
                		dummyElevInfo.Planned_stops[floor]=true
                		elevInfoMap[lowestCostIP]=dummyElevInfo
                	}
                }
                for floor,order :=range(d){
					dummycost=1000000
					if order{
						for key,val:=range(elevInfoMap){
                			costMap[key]=costFunc(val,false,floor+1)
                		}
                		for costkey,cost:=range(costMap){
                			if cost<dummycost{
                				dummycost=cost
                				lowestCostIP=costkey
                			}
                		}
                		dummyElevInfo=elevInfoMap[lowestCostIP]
                		dummyElevInfo.Planned_stops[floor+1]=true
                		elevInfoMap[lowestCostIP]=dummyElevInfo
                	}
                }
                elevInfoChan<-elevInfoMap
         }
    }
}


func costFunc(elevator gen.ElevButtons,dir bool,floor int)int{
	searchFloor:=elevator.Current_floor
	searchDir:=elevator.Dir
	cost:=0
	dummy:=0
	var i int
	for _,val:=range(elevator.Planned_stops){
		if val{
			dummy=1
		}
		if dummy==0 && (searchFloor==floor && searchDir==dir) {
			return 0
		}
	}
	for !(searchFloor==floor && searchDir==dir){
		cost++
		if searchDir{
			searchFloor++
		} else{
			searchFloor--
		}
		i=searchFloor
		dummy=0
		for i<gen.N_FLOORS && i>-1{
			if elevator.Planned_stops[i]{
				dummy=1
			}
			if i==floor{
				dummy=1
			}	
			if searchDir{
				i++
			} else {
				i--
			}
		}
		if dummy==0 {
			searchDir= !searchDir
		}
	}
	return cost
}
/*
func getMyIP() string{
    allIPs,err:=net.InterfaceAddrs()
    if err!=nil{
        fmt.Println("IP receiving errors!!!!!!!!\n")
        return ""
    }
    return strings.Split(allIPs[1].String(),"/")[0]
}

*/









