class Player:

    IP = '127.0.0.1'
    Name = ''
    Score = 0
    connected = 0

    def sayName(self):
        print 'name = %s ' % self.Name

    def sayIP(self):
        print 'IP = %s ' %self.IP

    def sayBoth(self):
		self.sayName()
		self.sayIP()
		print self.connected

    def setIP(self, newIP):
        self.IP = newIP

    def setName(self, newName):
        self.Name = newName

    def setScore(self, newScore):
	self.Score = newScore

    def sayScore(self):
	print 'Score = %s ' %self.Score

    def getName(self):
        return self.Name

    def getIP(self):
	return self.IP

    def getScore(self):
	return self.Score

    def setConnected(self):
	self.connected = 1
	
    def setDisconnected(self):
	self.connected = 0
	





