Create Clusters

-----------------------------------------------------------------------------------------------------------

1. What is this?

	This project contains a library called "cluster.go" which is used to create n-number of servers 
	and send message between them.

2. Working

	This is a library is used
	   - To create n-number of servers 
	   - To queue outgoing message in the outbox of each server
	   - To queue incoming message in the inbox of each server
	   - Send and receive messages between servers

3. How to run?

		- go get github.com/onkarkore/cluster/
		- create config files for each server
			python configcreator.py
		- go test github.com/onkarkore/cluster/
		

		
4. Testing

		Set following variables in cluster_test.go 
			MSG_TO_EACH_SERVER - Number of message send by each server to other peers 
			env.RPid - If Broadcast set to -1 else set port number of receiver so that each server to to this server only



		Test file is used to check working of cluster library.In which, we call methods implemented 
		in cluster library.
		Following things are to be done for testing 
		 - Create new socket for each server
		 - Set all properties related to each server
		 - Create a envelope(sending message) for each server
	 	 - Send message between servers, count number of message sent
		 - Receive message sent by other servers , count number of received message
		 - Display number of sent and received message and verify that this count is same

5. Configuration file (cluster.conf)

		This file contains list of all sever addresses.
		e.g. tcp://127.0.0.1:2001
		     tcp://127.0.0.1:2002
		

6. References 

	- http://golangtutorials.blogspot.in/2011/10/gotest-unit-testing-and-benchmarking-go.html
	- http://stackoverflow.com/questions/10728863/how-to-lock-synchronize-access-to-a-variable-in-go-during-concurrent-goroutines
	- go language tutorial















