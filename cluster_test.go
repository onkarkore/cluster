/*
	Author : Onkar Kore
	This is a main program to create multiple server and test cluster interface.
*/

package cluster

import (
	"time"
	"fmt"
	"os"
	"sync"
	"testing"
	"strconv"
	"bufio"
	"strings"
)

type Counter struct {
    mu  sync.Mutex
    x   int64
}

var (
	sentmsgcount int
	receivemsgcount int
	MSG_TO_EACH_SERVER int
	EXPECTED_COUNT int
	env  Envelope
	views Counter
	confFile = "cluster.conf"
)


func Test_main(t *testing.T) {

	sentmsgcount=0
	receivemsgcount=0
	
	/* Number of message send by each server to other peers */	
	MSG_TO_EACH_SERVER = 10

	EXPECTED_COUNT = 0

	/* If Broadcast set to -1 else set port number of receiver so that each server to to this server only */
	env.RPid=-1
	
	if len(os.Args) < 1 {
		fmt.Printf("<usage> : go run server.go \n")
		return
	}
	
	CreateSeverAndOutbox()    				
}



func CreateSeverAndOutbox(){
		
	totakNumberofServers:=TotalNumberOfServers(confFile)


	/* Expected count of message should be received */
	if env.RPid == -1 {
		EXPECTED_COUNT = (totakNumberofServers)*(totakNumberofServers-1)*MSG_TO_EACH_SERVER	
	} else {
		EXPECTED_COUNT = (totakNumberofServers-1)*MSG_TO_EACH_SERVER
	}


	var serverport[] int

	/* Read config file for number of servers. */
	configfile, _ := os.OpenFile(confFile, os.O_RDWR, 0600)
	i:=0
	var e error
	buf := bufio.NewReader(configfile)
	if e != nil {
		
	} else {
		for {
			line, err := buf.ReadString('\n')
			if err != nil {
				break
			}
			line = line[:len(line)-1]
			parts := strings.Split(line, ":")
			value,_ := strconv.Atoi(parts[2])
			serverport=append(serverport,value)
			i = i+1
		}
	}
	configfile.Close()



	num:=0
	for i:=1;i<=totakNumberofServers;i++{

		s := "cluster"+strconv.Itoa(serverport[i-1])+".conf"
		serv:=CreateServer(s)
	
		go  SendMsgtoServers(serv[num].Outbox(),serv[num])	
		go  ReceiveMsg(serv[num].Inbox(),serv[num])		
		go  inboxmsg(serv[num])
		go  outboxmsg(serv[num])
				
	}


	for
	{
		select {	
			case <- time.After(5 * time.Millisecond): 
				//fmt.Println("Total msg sent : ",sentmsgcount)
				//fmt.Println("Total msg received : ",receivemsgcount)
		}
	}
}


/* Function to queue outbox  */
func outboxmsg(serv  ServerData){
	for j:=0; j < MSG_TO_EACH_SERVER; j++ {

		env.MsgId=int64(serv.ServerID)
		env.Msg="Message from server "+serv.ServerAdd
		if env.RPid ==  BROADCAST {
			sentmsgcount=sentmsgcount+(len(serv.PeersId)-1)
		} else {
			if env.RPid==serv.ServerID {	
			} else {
				sentmsgcount=sentmsgcount+1
			}
		}
		if sentmsgcount%10==0{
			time.Sleep(10*time.Millisecond)
		}
		serv.Outbox() <- &env		
	}
}


/* Function to retrive messages from inbox  */
func inboxmsg(serv  ServerData){
	sentmsg_closed := false
	for {
		if (sentmsg_closed) { return }
		select {
        		case _, ok := <- serv.Inbox():
            			if (!ok) {
					fmt.Println("Inbox channel closed!")	
			        } else {
					receivemsgcount=receivemsgcount+1
					views.Add(1)
				}
			case <- time.After(5 * time.Millisecond): 
				if receivemsgcount==EXPECTED_COUNT{
					//fmt.Println("Total msg sent : ",views.x)					
					//fmt.Println("Total msg received : ",receivemsgcount)
				        os.Exit(0)
				}
		}
	}
}

func (c *Counter) Add(x int64) {
    c.mu.Lock()
    c.x += x
    c.mu.Unlock()
}













