import socket, traceback
from GameRoom import GameRoom
import json
from threading import Thread
import time
import sys




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

def obj_to_dict(obj):
    players=[]
    for x in obj.players:
        player = { 'Ip': x.IP, 'Name': x.Name }
        players.append(player)
    data = { 'Name':'game room info', 'ID': obj.id, 'Count': obj.count, 'LocalName': obj.localName, 'Players': players, 'State':'start'}
    return data

def obj_to_rejoinMsg(obj, newPlayer):
    players=[]
    for x in obj.players:
    	if x.connected == 1:
        	player = { 'Ip': x.IP, 'Name': x.Name }
        	players.append(player)
    data = { 'Name':'game room info', 'ID': obj.id, 'Count': obj.getConnectedPlayer(), 'LocalName': obj.localName, 'Players': players, 'State':'rejoin', 'NewPlayer': newPlayer}
    return data

def obj_to_dynamicJoinMsg(obj, newPlayer):
    players=[]
    for x in obj.players:
        player = { 'Ip': x.IP, 'Name': x.Name }
        players.append(player)
    data = { 'Name':'game room info', 'ID': obj.id, 'Count': obj.count, 'LocalName': obj.localName, 'Players': players, 'State':'join', 'NewPlayer': newPlayer}
    return data

def obj_to_startMsgToUI(obj):
    data = { 'Name':'start'}
    return data

def obj_to_reJoinMsgToUI(obj):
    data = { 'Name':'rejoin'}
    return data

def obj_to_closeGame(obj):
    data = { 'Name':'game room info', 'State':'close'}
    return data



def startServer():
    roomLimit = int(sys.argv[2])
    room1 = GameRoom()
    room1.setLocalName(sys.argv[1])
    host = ''                               # Bind to all interfaces
    port = 51425
    dataGame = ''
    data_encoded = ''
    rejoinSuccess = 0
    recoveredFlag = 0

    print room1.localName
    room1.init4Players()


    
    

    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    s.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)
    s.bind((host, port))
    
    
    CONNToGo = ('localhost', 9999)
    
    #def fn_clientToGo(string, *args):
    csToGo = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    csToGo.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR,1)
    #csToGo.bind(CONNToGo)

    CONNToUI = ('localhost', 8888)
    
    #def fn_clientToGo(string, *args):
    csToUI = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    csToUI.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR,1)
    
    
    while 1:
        try:
			message, address = s.recvfrom(8192*16)
	    #messageFromGo, addressFromGo = csToGo.recvfrom(8192)
	    #print addressFromGo[0]
	    #print "Got data from goServer:", addressFromGo
	    #print messageFromGo

			print address[0]
			print "Got data from", address
			print message
			decoded = json.loads(message)
			messageType = decoded["Name"]
			print decoded["Name"]
			if messageType == "game room info":
				if recoveredFlag == 0:
					recoveredFlag = 1
					print 'start self recovery:'
					# 1. update current room1 update info
					players =  decoded["Players"]
					for player in players:
						playerName = player["Name"]
						playerIP = player["Ip"]
						if room1.checkIfExist(playerName) == 0:
							room1.addNewPlayer(playerName, playerIP)
						# 2. send to self go server the msg
					selfRecoverMsg = obj_to_rejoinMsg(room1, room1.localName)
					data_selfRecoverMsg = json.dumps(selfRecoverMsg)
					csToGo.sendto(data_selfRecoverMsg, CONNToGo)
			# 3. send to self UI the rejoin msg
			#selfRecoverMsgToUI = obj_to_startMsgToUI(room1)
			#data_selfRecoverMsgToUI = json.dumps(selfRecoverMsgToUI)
                	#csToUI.sendto(data_selfRecoverMsgToUI, CONNToUI)


			if messageType == "request":
				source = decoded["Source"]
				print "source + room1.localName+ decoded_source:"
				print source
				print room1.localName
				print decoded["Source"]
				if room1.gameStart == 0:
					if room1.checkIfExist(source) == 0:
						room1.addNewPlayer(source,address[0])
						room1.setPlayerConnected(source)
						room1.printPlayers()
                		#data1 = 'ack'+','+ room1.localName
                		#s.sendto(data1, address)
            #if room1.checkIfExist(message) == 0:
            #    room1.addNewPlayer(message,address[0])
            #    room1.printPlayers()
            #    s.sendto(room1.localName, address)
					if room1.count == roomLimit:
						print 'enter sending to go server:'
						room1.printPlayers()
						dataGame =  obj_to_dict(room1)
						data_encoded = json.dumps(dataGame)
                #data_decoded = jsonpickle.dumps(room1)
						print data_encoded
						data_decoded = json.loads(data_encoded)
						print data_decoded
						csToGo.sendto(data_encoded, CONNToGo)
				#room1.startGame = 1

				if room1.gameStart == 1:
					if room1.checkIfExist(source) == 1:
						print 'send rejoin data to go server:'
						rejoinPlayer = source
						room1.setPlayerConnected(source)
						dataRejoin = obj_to_rejoinMsg(room1, rejoinPlayer)
						data_reJoin = json.dumps(dataRejoin)
						csToGo.sendto(data_reJoin, CONNToGo)
					# send "room info to bsServer:"
						s.sendto(data_reJoin, (address[0],port))
						print json.loads(data_reJoin)
						print 'finish sending rejoin msg to self'
					elif room1.count < roomLimit:
						print 'dynamic join: send new join to go server:'
						room1.addNewPlayer(source,address[0])
						room1.printPlayers()
						dataDynamicJoin = 'ack'+','+GameRoom.localName
						s.sendto(dataDynamicJoin, address)
						print room1.count
				#dataDynamicJoin = {'Name':'game room info', 'Id': room1.id, 'Player': {'Ip': address[0], 'Name':source}, 'Count':room1.count, 'LocalName':room1.localName, 'State':'join'}
						dataDynamicJoin = obj_to_dynamicJoinMsg (room1, source)
						data_DynamicJoin = json.dumps(dataDynamicJoin)
						csToGo.sendto(data_DynamicJoin, CONNToGo)
			
			if messageType == "connection success":
				if room1.gameStart == 0:
					room1.gameStart = 1
					selfRecoverMsgToUI = obj_to_startMsgToUI(room1)
					data_selfRecoverMsgToUI = json.dumps(selfRecoverMsgToUI)
					csToUI.sendto(data_selfRecoverMsgToUI, CONNToUI)

                #csToUI.sendto(data_encoded, CONNToUI)
			if messageType == "gameresult":
				room1.startGame = 0
		#players = decoded["Players"]
		#for player in players:
		#	playerName = player["Name"]
		#	playerScore = player["Score"]
		#	room1.setPlayerScore(playerName, playerScore)
		#room1.printPlayers()

		#for player in players:
		#	playerName = player["Name"]
		#	if playerName != room1.localName:
		#		room1.removePlayer(playerName)
				room1.printPlayers()
				dataGameClose =  obj_to_closeGame(room1)
				data_GameClose = json.dumps(dataGameClose)
                #data_decoded = jsonpickle.dumps(room1)
				print data_GameClose
				data_decoded_GameClose = json.loads(data_GameClose)
				print data_decoded_GameClose
				csToGo.sendto(data_GameClose, CONNToGo)

				time.sleep( 15 )
				if room1.count == roomLimit:
					print 'start a new game after close a game, sending to go server:'
					dataGame =  obj_to_dict(room1)
					data_encoded = json.dumps(dataGame)
                #data_decoded = jsonpickle.dumps(room1)
					print data_encoded
					data_decoded = json.loads(data_encoded)
					print data_decoded
					csToGo.sendto(data_encoded, CONNToGo)
			#room1.startGame = 1
		
			if messageType == "disconnect":
				peer = decoded["Peer"]
				room1.setPlayerDisconnected(peer)
				print '------------ shoulb be 0:-----------' 
				#peer.sayBoth()
				room1.printPlayer(peer)
				print '-------------------------------------' 

				#cancelQuality = {'Name':'disconnect', 'Tank_id': peer}
				#cancelQualityEncoded = json.dumps(cancelQuality)
				#csToUI.sendto(cancelQualityEncoded, CONNToUI)
       
        # Acknowledge it.
        #       s.sendto("allow to join in the game", address)
        except (KeyboardInterrupt, SystemExit):
            raise
        except:
            traceback.print_exc()



startServer()
