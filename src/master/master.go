package master

import(
    //"fmt"
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
                			costMap[key]=costFunc(val,1,floor)
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
                			costMap[key]=costFunc(val,-1,floor+1)
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


func costFunc(elevator gen.ElevButtons,dir int,floor int)int{
	searchFloor:=elevator.Current_floor
	searchDir:=elevator.Dir
	cost:=0
	var dummy int
	for _,val:=range(elevator.Planned_stops){
		if val{
			dummy=1
		}
		if dummy==0 && (searchFloor==floor && (searchDir==dir || searchDir==0)) {
			return 0
		}
	}
	if searchDir==0{
		if floor<searchFloor{
			searchDir=-1
		} else {searchDir=1}
	}
	for !(searchFloor==floor && (searchDir==dir || searchDir==0)){
		if searchDir==1{
			if isAbove(searchFloor,floor,elevator.Planned_stops)==0{
				searchDir=-1
				continue
			}
		} else {
			if isBelove(searchFloor,floor,elevator.Planned_stops)==0{
				searchDir=1
				continue
			}
		}
		if searchDir==1 {
			searchFloor++
		} else if searchDir==-1{
			searchFloor--
		}
		cost++
	}
	return cost
}

func isAbove(myFloor,checkFloor int,Planned_stops [gen.N_FLOORS]bool)int{
	for i:=myFloor+1; i<gen.N_FLOORS;i++{
		if Planned_stops[i] || checkFloor==i{
			return 1
		}
	}
	return 0
}

func isBelove(myFloor,checkFloor int,Planned_stops [gen.N_FLOORS]bool)int{
	for i:=myFloor-1; i>-1;i--{
		if Planned_stops[i] || checkFloor==i{
			return 1
		}
	}
	return 0
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









