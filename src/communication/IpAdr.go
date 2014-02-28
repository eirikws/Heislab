package communication
import(
    "net"
    "fmt"
    "strings"
    "strconv"
    "sort"
    "time"
)

func getMyIP() string{
    allIPs,err:=net.InterfaceAddrs()
    if err!=nil{
        fmt.Println("IP receiving errors!!!!!!!!\n")
        return ""
    }
    return strings.Split(allIPs[1].String(),"/")[0]
}

func getBIP(MyIP string) string{
    IP:=strings.Split(MyIP,".")
    return IP[0]+"."+IP[1]+"."+IP[2]+".255"
}

func IPsort(list []IPandTimeStamp) []IPandTimeStamp{ 
    fmt.Println("enter ipsort")
    ipbase:=strings.Split(list[0].IPadr,".")[0:3]
    var intlist []int
    var newlist []IPandTimeStamp
    var iptime IPandTimeStamp
    for i,val:=range(list){
        i,_=strconv.Atoi(strings.Split(val.IPadr,".")[3])
        intlist=append(intlist,i)
    }
    sort.Ints(intlist)
    for _,val:= range(intlist){
        iptime=IPandTimeStamp{ipbase[0]+"."+ipbase[1]+"."+ipbase[2]+"."+strconv.Itoa(val),time.Now()}
        newlist=append(newlist,iptime)
    }
    i:=0
    for i<len(list)-1{
        for j:=0; j<len(list); j++{
            if newlist[i].IPadr==list[j].IPadr{
                newlist[i].Timestamp=list[j].Timestamp
                i++
            }
        }
    }
    return newlist
}
    
func timeStampCheck(list chan []IPandTimeStamp){
    var IPlist []IPandTimeStamp
    var newlist []IPandTimeStamp
    IPlist=<-list
    
    for _,val:= range(IPlist){
        if val.Timestamp.Before(time.Now()){
            fmt.Println("delete")
            for _,bval:=range(IPlist){
                if val.IPadr!=bval.IPadr{
                    newlist=append(newlist,bval)
                
                }   
            }
               
        }
    }
    fmt.Println(newlist)
    list<-newlist
    // time.Sleep(time.Millisecond*100)
}









