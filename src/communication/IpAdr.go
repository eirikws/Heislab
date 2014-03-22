package communication
import(
    "net"
    "fmt"
    "strings"
    "strconv"
    "time"
)

func getMyIP() string{
    allIPs,err:=net.InterfaceAddrs()
    if err!=nil{
        fmt.Println("IP receiving errors!\n")
        return ""
    }
    return strings.Split(allIPs[1].String(),"/")[0]
}

func getBIP(MyIP string) string{
    IP:=strings.Split(MyIP,".")
    return IP[0]+"."+IP[1]+"."+IP[2]+".255"
}


func getSmallestIP(list map[string]time.Time)string{
	var smallestIP = 10000
	var ip int
	var ipbase []string
	for key,_:=range(list){
		ip,_=strconv.Atoi(strings.Split(key,".")[3])
		if ip<smallestIP{
			smallestIP=ip
		}
		ipbase=strings.Split(key,".")[0:3]
	}
	return ipbase[0]+"."+ipbase[1]+"."+ipbase[2]+"."+strconv.Itoa(smallestIP)
}

//deletes all entries which timestamp has run out except MyIP
func timeStampCheck(list chan map[string]time.Time,deletedIP chan string,MyIP string){
    var IPlist map[string]time.Time
    for{
        IPlist=<-list
        for key,val:= range(IPlist){
            if val.Before(time.Now()) && key!=MyIP{
                deletedIP<-key
                delete(IPlist,key)
                break
            }
        }
        list<-IPlist
    }

}









