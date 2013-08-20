import socket, traceback
from GameRoom import GameRoom
import pickle
import json
import ui
from threading import Thread
import time

room1 = GameRoom()
room1.init4Players()
room1.printPlayers()

def convert_to_builtin_type(obj):
    print 'default(', repr(obj), ')'
    # Convert objects to a dictionary of their representation
    d = { '__class__':obj.__class__.__name__, 
        '__module__':obj.__module__,
    }
    d.update(obj.__dict__)
    return d



def dict_to_object(d):
    if '__class__' in d:
        class_name = d.pop('__class__')
        module_name = d.pop('__module__')
        module = __import__(module_name)
        print 'MODULE:', module
        class_ = getattr(module, class_name)
        print 'CLASS:', class_
        #args = dict( (key.encode('ascii'), value) for key, value in d.items())
        #print 'INSTANCE ARGS:', args
    #inst = class_(**args)
        inst = class_()
    
    else:
        inst = d
    return inst




def startServer():

    host = ''                               # Bind to all interfaces
    port = 51425



    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    s.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)
    s.bind((host, port))


    CONNToGo = ('localhost', 55556)

    #def fn_clientToGo(string, *args):
    csToGo = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    csToGo.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR,1)
    #csToGo.bind(CONNToGo)
  

    while 1:
        try:
            message, address = s.recvfrom(8192)
            print address[0]
            print "Got data from", address
            if room1.count < 4:
                
                print message
                splitMsg = message.split(',')
                print splitMsg[0]
                print splitMsg[1]
                if room1.checkIfExist(splitMsg[1]) == 0:
                    room1.addNewPlayer(splitMsg[1],address[0])
                    room1.printPlayers()
                data = 'ack'+','+room1.localName
                s.sendto(data, address)
                print room1.count
                #if room1.checkIfExist(message) == 0:
                #    room1.addNewPlayer(message,address[0])
                #    room1.printPlayers()
                #    s.sendto(room1.localName, address)
            if room1.count == 1:
                print 'enter sending to go server:'
                data = json.dumps(room1, default=convert_to_builtin_type)
                myobj_instance = json.loads(data, object_hook=dict_to_object)
                print myobj_instance
                print myobj_instance.count
                print myobj_instance.localName
                print myobj_instance.count

                
                myobj_instance.printPlayers()                
                csToGo.sendto(data, CONNToGo)
#	data = s.recv(1024).strip()
#	print "{} wrote:".format(client_address[0])
            print message
            time.sleep(.5)

    
        # Acknowledge it.
#       s.sendto("allow to join in the game", address)
        except (KeyboardInterrupt, SystemExit):
            raise
        except:
                traceback.print_exc()

def startClient():
    #roomclient = GameRoom()
    
    BC_PORT = 51425
    
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    s.bind(('',0))
    s.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)
    
    #s = socket(AF_INET, SOCK_DGRAM)
    #s.setsockopt(SOL_SOCKET, SO_BROADCAST, 1)
    #data = pickle.dumps(roomclient)
    
    #s.sendto("request to join in the game", ('<broadcast>', BC_PORT))
    data = 'request'+','+room1.localName
    s.sendto(data, ('<broadcast>', BC_PORT))
    
    try:
	    # Connect to server and send data
        #	   	sock.connect((HOST, PORT))
        #	   	sock.sendall(data + "\n")
        
	    # Receive data from the server and shut down
        #              	received = sock.recv(1024)
		received, address = s.recvfrom(8192)
    
    
    
    finally:
        #		sock.close()
		s.close()
    
    #	print "Sent:     {}".format(data)
    print "Received from: {}".format(address)
    print "Message: {}".format(received)


def startUI():
    x= raw_input('Do you want to join in the game?(Y/N)')
    if (x=='Y' or  x == 'y' ):
        startClient()
        print 'ok, you can join the game'
    else:
        x = raw_input('Do you want to join in the game? (Y/N)')

#startServer()
#ui.startUI()
if __name__=='__main__':
    
    try:
        Thread(target=startServer).start()
        Thread(target=startClient).start()
    except (KeyboardInterrupt, SystemExit):
        raise
    except Exception, errtxt:
        print errtxt
    except:
        traceback.print_exc()

