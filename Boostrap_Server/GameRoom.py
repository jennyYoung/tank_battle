from Player import Player

class GameRoom(object):

    name = ''
    id = 0
    players = []
    count = 0
    localName = 'Jerry'
    gameStart = 0

    
    def _init_(self):
        self.init4Players()
    
    def _repr_(self):
        return '<MyObj(%s)>' %self.s

    def sayCount(self):
        print 'there are %s players in a game room ' % self.count

    def sayId(self):
        print 'game room ID = %s ' %self.Id
    
    def sayPlayers(self):
        print 'This game room has the following players: %s ' %self.players

    
    def printPlayers(self):
		print '=====================This game room has the following players:==========='
		print 'number of players: '
		print self.count 
		print 'localname:'
		print self.localName
		for item in self.players:
			item.sayBoth()
		print '=====================This game room has the following players:==========='


    def sayAll(self):
        self.sayId()
        self.sayCount()
        self.sayPlayers()

    def setId(self, newId):
        self.id = newId

    def incrementCount(self):
        self.count = self.count+1

    def addPlayer(self, newPlayer):
        self.players.append(newPlayer)

    def setLocalName(self, newName):
        self.localName = newName

    def init4Players(self):
		tom = Player()
		tom.setName('')
		tom.setIP('127.0.0.1')
        
		jerry = Player()
		jerry.setName('')
		jerry.setIP('127.0.0.1')
		
		alice = Player()
		alice.setName('')
		alice.setIP('127.0.0.1')

		jennifer = Player()
		jennifer.setName('')
		jennifer.setIP('127.0.0.1')

		self.addPlayer(tom)
		self.addPlayer(jerry)
		self.addPlayer(alice)
		self.addPlayer(jennifer)
		self.printPlayers()

    def setIndexedIP(self,newIP):
		self.players[self.count].setIP(newIP)
    
    def setIndexedName(self,name):
		self.players[self.count].setName(name)
    
    
    def addNewPlayer(self, newName,newIP):
		self.setIndexedName(newName)
		self.setIndexedIP(newIP)
		self.incrementCount()

    def checkIfExist(self, name):
		for item in self.players:
			if name == item.Name:
				return 1
		return 0

    def setPlayerScore(self, name, score):
		for item in self.players:
			if name == item.Name:
				item.setScore(score)

    def setPlayerDisconnected(self, name):
		for item in self.players:
			if name == item.Name:
				item.setDisconnected()
		
    def setPlayerConnected(self, name):
		for item in self.players:
			if name == item.Name:
				item.setConnected()

    def removePlayer(self, name):
		for item in self.players:
			if name == item.Name:
				self.players.remove(item)

    def getConnectedPlayer(self):
		connectedNum = 0
		for item in self.players:
			if item.connected == 1:
				connectedNum = connectedNum + 1
		return connectedNum	
	
    def printPlayer(self, playerName):
		for item in self.players:
			if item.Name == playerName:
				item.sayBoth()
  
    
#room1 = GameRoom()
#room1.addNewPlayer('128.2.247.15')
#room1.addNewPlayer('128.2.247.16')
#room1.
()
