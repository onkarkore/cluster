/*
	Author : Onkar Kore
	This is a main main program to create multiple server and test cluster interface.
*/

package cluster

import (
	//cluster "github.com/onkarkore/cluster"
	"time"
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
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
)





func Test_main(t *testing.T) {

	sentmsgcount=0
	receivemsgcount=0
	
	/* Number of message send by each server to other peers */	
	MSG_TO_EACH_SERVER = 1

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
	var serv [30] ServerData	
	count:=0

	var s  ServerData
	var Peers = s.Peers()
	var hostaddr =  PeersAddress()

	
	/* Expected count of message should be received */
	if env.RPid == -1 {
		EXPECTED_COUNT = ((len(Peers))*(len(Peers)-1))*MSG_TO_EACH_SERVER	
	} else {
		EXPECTED_COUNT = (len(Peers)-1)*MSG_TO_EACH_SERVER
	}


	/* Create n number of server objects and initialize its properties */
	for num:=0;num<len(Peers);num++{		
		serv[count].ServerID = Peers[num]
		//serv[count].ServerAdd = "tcp://127.0.0.1:"+strconv.Itoa(Peers[num])
		serv[count].ServerAdd = hostaddr[num]+":"+strconv.Itoa(Peers[num])
		serv[count].PeersId = Peers

		serv[count].ServerSocket =  CreateSocket()
		serv[count].ServerSocket.Bind(serv[count].ServerAdd)

		serv[count].Outboxd = make(chan * Envelope)	
		serv[count].Inboxd = make(chan * Envelope)

		fmt.Println("I: echo service is ready at ", serv[count].ServerAdd)	
		count++	
	}

	
	count=0
	

	for num:=0;num<len(Peers);num++{
		go  SendMsgtoServers(serv[count].Outbox(),serv[count])	
		go  ReceiveMsg(serv[count].Inbox(),serv[count])		
		go inboxmsg(serv[count])
		go outboxmsg(serv[count])
		count++
	}

	for
	{
		select {	
			case <- time.After(5 * time.Second): 
				fmt.Println("Total msg sent : ",sentmsgcount)
				fmt.Println("Total msg received : ",receivemsgcount)
				//fmt.Println("Total msg x : ",views.x)
			        
				//fmt.Println("Total msg received : ",EXPECTED_COUNT)			
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
        		case _, strbry_ok := <- serv.Inbox():
            			if (!strbry_ok) {
					fmt.Println("Inbox channel closed!")	
			        } else {
					//fmt.Println("Message Received -----------------> ", *cakeName)
					receivemsgcount=receivemsgcount+1
					views.Add(1)
				}
			case <- time.After(10 * time.Second): 
				if receivemsgcount==EXPECTED_COUNT{
					fmt.Println("Total msg sent : ",views.x)					
					fmt.Println("Total msg received : ",receivemsgcount)
				        //fmt.Println("Total msg sent : ",sentmsgcount)
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













