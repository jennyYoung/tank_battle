import socket, traceback
from GameRoom import GameRoom
import pickle
import json


def dict_to_object(d):
    if '__class__' in d:
        class_name = d.pop('__class__')
        module_name = d.pop('__module__')
        module = __import__(module_name)
        print 'MODULE:', module
        class_ = getattr(module, class_name)
        print 'CLASS:', class_
        args = dict( (key.encode('ascii'), value) for key, value in d.items())
        print 'INSTANCE ARGS:', args
        #inst = class_(**args)
        inst = class_()

    else:
        inst = d
    return inst


def startGoServer():
    
    CONNToGo = ('localhost', 9999)
            
  
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    #s.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)
    s.bind(CONNToGo)
    
 
    while 1:
        try:
		message, address = s.recvfrom(32768)
		print address[0]
		print "Got data from", address
		data_decoded = json.loads(message)
		name = data_decoded["Name"]
		print "Name "
		print name
		players = data_decoded["Players"]
		for obj in players:
			print "IP"
			print obj["Ip"]
			print "Name"
			print obj["Name"]
		print "Count"	
		print data_decoded["Count"]
		print "LocalName"
		print data_decoded["LocalName"]
#		print "State"
#		print data_decoded["State"]"
		

	    
            #data_decoded = jsonpickle.loads(message)
		print data_decoded
        #roomserver.printPlayers()
        
        # Acknowledge it.
        #       s.sendto("allow to join in the game", address)
        except (KeyboardInterrupt, SystemExit):
            raise
        except:
            traceback.print_exc()


startGoServer()
