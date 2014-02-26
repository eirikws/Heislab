package communication
import(
    "net"
    "fmt"
    "strings"
    "strconv"
    "sort"
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

func IPsort(list []string) []string{
    ipbase:=strings.Split(list[0],".")[0:3]
    var newlist []int
    for i:=0; i<len(list); i++{
        i,_=strconv.Atoi(strings.Split(list[i],".")[3])
        newlist=append(newlist,i)
    }
    sort.Ints(newlist)
    for i,val:= range(newlist){
        list[i]=ipbase[0]+"."+ipbase[1]+"."+ipbase[2]+"."+strconv.Itoa(val)
    }
    fmt.Println(list)
    return list
}
