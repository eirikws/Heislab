package master

import(
    "fmt"
    gen "./../genDecl"
    "net"
    "strings"
)


func Master(master chan string,elevInfoChan chan map[string]gen.ElevInfo){
    var mst,lowestCostIP string
    var upButtons,downButtons [gen.N_FLOORS-1]bool 
    elevInfoMap:=make(map[string]gen.ElevInfo)
    costMap:=make(map[string]int)
    var dummycost int
    var dummyElevInfo gen.ElevInfo
    MyIP:=getMyIP()
    for{
        select{
            case mst=<-master:
            case elevInfoMap=<-elevInfoChan:
            	spreadOrders(elevInfoMap)
                if mst!=MyIP{
                	elevInfoChan<-elevInfoMap
                    continue
                }
                //the up and down buttons are the same for every elevator
                //take out from MyIP because MyIP is alway in the map
                upButtons=elevInfoMap[MyIP].U_buttons
                downButtons=elevInfoMap[MyIP].D_buttons
                
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
				for floor,order :=range(upButtons){
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
                for floor,order :=range(downButtons){
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


func costFunc(elevator gen.ElevInfo,dir int,floor int)int{
	searchFloor:=elevator.Current_floor
	searchDir:=elevator.Dir
	cost:=0

	if searchDir==0{
		if floor<searchFloor{
			searchDir=-1
		} else {searchDir=1}
	}
	for !(searchFloor==floor && (searchDir==dir || searchDir==0)){
		//check if there is anything in the direction you are searching,
		//if not, change the direction and restart the loop without iterating
		//searchFloor or increasing the cost
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
	fmt.Println(cost)
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

func spreadOrders(info map[string]gen.ElevInfo){
	u:=[gen.N_BUTTONS]bool{false}
	d:=[gen.N_BUTTONS]bool{false}
	var temp gen.ElevInfo
	for _,val:=range(info){
		
		for i:=0; i<gen.N_FLOORS-1;i++{
			if val.U_buttons[i]==true{
				u[i]=true
			}
			if val.D_buttons[i]==true{
				d[i]=true
			}
		}
	}
	for key,val:=range(info){
		temp=val
		for i:=0; i<gen.N_FLOORS-1;i++{
			if u[i]==true{
				temp.U_buttons[i]=true
			}
			if d[i]==true{
				temp.D_buttons[i]=true
			}
		}
		info[key]=temp
	}
}

func getMyIP() string{
    allIPs,err:=net.InterfaceAddrs()
    if err!=nil{
        fmt.Println("IP receiving errors!!!!!!!!\n")
        return ""
    }
    return strings.Split(allIPs[1].String(),"/")[0]
}


