package master

import(
    "fmt"
    gen "./../genDecl"
)


func Master(master chan string,elevInfoChan chan map[string]gen.ElevButtons,orders chan string){
    //var mst string
    elevInfoMap:=make(map[string]gen.ElevButtons)
    var dummyElevInfo gen.ElevButtons
    for{
        select{
          //  case mst=<-master:
          //      fmt.Println(mst)
            case elevInfoMap=<-elevInfoChan:
              //  if mst!=com.GetMyIP(){
              //      continue
              //  }
                
                for key,val := range(elevInfoMap){
                    dummyElevInfo=val
                    for i,ival:=range(val.C_buttons){
                        if ival{
                            dummyElevInfo.Planned_stops[i]=true
                            fmt.Println("muha")
                        }
                    
                    }
                    elevInfoMap[key]=dummyElevInfo
                }
                
                elevInfoChan<-elevInfoMap
                
         }
            
    }
}
