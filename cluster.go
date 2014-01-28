/*
	Author : Onkar Kore
	This is a cluster interface to create differenet 
	server and send and recive meassages between them.
*/

package cluster

import (
	zmq "github.com/pebbe/zmq4"
	"fmt"
	"os"
	"time"
	"strconv"
	"log"
	"bufio"
	"io"
	"strings"	
)

var (
	hostaddress []string 
)



const (
	BROADCAST = -1
)


/* 
   Message format 
	RPid  - id of receiver 
	MsgId - unique message id (optional)
	Msg   - Actual data
*/
type Envelope struct {
	RPid int
	MsgId int64
	Msg string
}


/*
   Server Data format
	ServerSocket - New Socket for server
	ServerID     - This server id
	ServerAdd    - This server address
	PeersId []   - This server peers
	Outboxd      - This server outbox
	Inboxd 	     - This server inbox
*/
type ServerData struct {
	ServerSocket *zmq.Socket	
	ServerID int
	ServerAdd string
	PeersId [] int
	PeersAdd [] string
	Outboxd chan *Envelope	
	Inboxd chan *Envelope	
}




type Server interface {
	Pid() int
	Peers() []int
	Outbox() chan *Envelope
	Inbox() chan *Envelope
}


/* Returns Pid of this server */
func (e ServerData) Pid() int {
	return e.ServerID
}


func PeersAddress() []string {
	return hostaddress
}

/* Returns Peers of this server */
func (e ServerData) Peers() []int {
	var av = []int{}
	
	f, err := os.OpenFile("cluster.conf", os.O_CREATE|os.O_RDONLY,0600)
        if err != nil {
		
    	    log.Fatal(err)
    	}
    	bf := bufio.NewReader(f)
    	for {
        	switch line, err := bf.ReadString('\n'); err {
        	case nil:
			line = line[:len(line)-1]
			parts := strings.Split(line, ":")
			value,_:=strconv.Atoi(parts[2])
			av = append(av,value)
			hostaddress = append(hostaddress,parts[0]+":"+parts[1])
        	case io.EOF:
        	    if line > "" {
        	        fmt.Println(line)
        	    }
        	    return av
	        default:
	            log.Fatal(err)
        	}
    	}
	return av
}


/* Returns Outbox of this server */
func (s ServerData) Outbox() chan *Envelope {	
	return s.Outboxd
}


/* Returns inbox of this server */
func (s ServerData) Inbox() chan *Envelope {
	return s.Inboxd
}


/* Returns new socket of this server */
func CreateSocket() *zmq.Socket{
	var server *zmq.Socket
	server, _ = zmq.NewSocket(zmq.REP)
	return server
}


/* Receive messages from other server and send back to sender */
func ReceiveMsg(inbox chan *Envelope,server2 ServerData){
	
	for  {			
		receivemsg, err := server2.ServerSocket.RecvMessage(0)
		
		line := receivemsg[0]
		line = line[1:len(line)-1]
		str := strings.Split(line," ")
		id,_ := strconv.Atoi(str[0])
		mid,_:=	strconv.Atoi(str[1])
		msg := ""
		for count:=2;count<len(str);count++{
			msg=msg+str[count]+" "
		}			
		
		r := Envelope{id,int64(mid),msg}
		go addinbox(server2,r)
		if err != nil {
			time.Sleep(10*time.Second)
			break 
		}
		server2.ServerSocket.SendMessage(&r)
	}
}

/* Fill inbox of this server */
func addinbox(server2 ServerData,e Envelope){
	server2.Inbox() <- &e
}


/* Send message to other servers */
func SendMsgtoServers(outbox chan *Envelope,server1 ServerData){
	sentmsg_closed := false
	for {
	        if (sentmsg_closed) { return }
		select {
        		case cakeName, strbry_ok := <-outbox:
            			if (!strbry_ok) {
			                sentmsg_closed = true
			                fmt.Println(" Outbox channel closed!")	
			        } else {
										
              				var e Envelope
					e=*cakeName
					
					if e.RPid==-1 {
						var peers = server1.PeersId

						for peers_count:=0;peers_count<len(peers);peers_count++{
							if peers[peers_count]==server1.ServerID{
								continue
							}
							var e1 Envelope
							e1.RPid = peers[peers_count]
							e1.MsgId= int64(peers[peers_count]  * 10 + 1)
							e1.Msg  = "Message form server ["+strconv.Itoa(server1.ServerID)+"] to server ["+strconv.Itoa(peers[peers_count])+"]"
							//serveraddress:="tcp://127.0.0.1:"+strconv.Itoa(peers[peers_count])
							serveraddress:=hostaddress[peers_count]+":"+strconv.Itoa(peers[peers_count])
							
							//time.Sleep(3*time.Second)
							client, _ := zmq.NewSocket(zmq.REQ)
							client.Connect(serveraddress)
							client.SendMessage(e1)
							
							//fmt.Println("Message send ---> ", e1.Msg)		
						}
		
					} else {
						if e.RPid==server1.ServerID{
						} else {
							serveraddress:="tcp://127.0.0.1:"+strconv.Itoa(e.RPid)
							client, _ := zmq.NewSocket(zmq.REQ)
							client.Connect(serveraddress)
							client.SendMessage(*cakeName)		
							//fmt.Println("Message send ---> ", *cakeName)
						}
				}				        	            	    	
			}   
		 	
		}   
    	}   				
}




