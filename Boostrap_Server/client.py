import socket
import sys,time
from socket import *
from GameRoom import GameRoom
import json

def obj_to_requestJoin(obj):
    data = { 'Name':'request'}
    return data


def startClient():
	roomclient = GameRoom()

	BC_PORT = 51425
        
	s = socket(AF_INET, SOCK_DGRAM)
	s.bind(('',0))
	s.setsockopt(SOL_SOCKET, SO_BROADCAST, 1)
    #data = pickle.dumps(roomclient)
    
        #s.sendto("request to join in the game", ('<broadcast>', BC_PORT))
	data = {'Name':'request', 'Source': sys.argv[1]}

	print "sys.argv + room1.localName:"
	print GameRoom.localName


    
	data_encoded = json.dumps(data)
	print data_encoded
	data_decoded = json.loads(data_encoded)
	print data_decoded
   
	s.sendto(data_encoded, ('<broadcast>', BC_PORT))
	#s.sendto(data_encoded, ('127.0.0.1', BC_PORT))
	s.close()
    #try:
	    # Connect to server and send data
#	   	sock.connect((HOST, PORT))
#	   	sock.sendall(data + "\n")

	    # Receive data from the server and shut down
#              	received = sock.recv(1024)
		#received, address = s.recvfrom(8192)
            

    
    #finally:
#		sock.close()


#	print "Sent:     {}".format(data)
#    print "Received from: {}".format(address)
#    print "Message: {}".format(received)
#    print received
#    rev = received.split(',')
#    print rev[0]
#    print rev[1]

startClient()
