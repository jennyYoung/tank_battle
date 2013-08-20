import client

def startUI():
    x= raw_input('Do you want to join in the game?(Y/N)')
    if (x=='Y' or  x == 'y' ):
        client.startClient()
        print 'ok, you can join the game'
    else:
        x = raw_input('Do you want to join in the game? (Y/N)')

startUI()