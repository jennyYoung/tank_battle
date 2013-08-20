/*
 * 18842 Distributed System Course Project
 * 
 * Distributed Tank Game
 * 
 * Team 22: Eason Qin (yichengq), Jin Yang (jiny), Chuyi Wei (chuyiw), Chia-Jun Lin (chiajunl)
 * Reference By the TankGame: author 韩顺平
 * 
 */
package GUI;

import java.awt.Color;
import java.awt.Font;
import java.awt.Graphics;
import java.awt.Image;
import java.awt.event.KeyEvent;
import java.awt.event.KeyListener;

import java.util.Vector;
import javax.swing.*;

import Events.Constants;
import Events.Events;
import Network.Receive;
import Network.Send;

public class Stage extends JPanel implements KeyListener  ,Runnable
{
	private static final long serialVersionUID = 1L;
	Tank OurTank=null;

	public static Vector<Events> positionQueue=new Vector<Events>();


	Receive receive; 
	Send send;

	int enSize = 4;
	public int enAll  = 6;
	public static boolean freeze=false;
	public static boolean start=false;

	public static Vector<Bomb> bomb = new Vector<Bomb>();

	JLabel be = new JLabel();

	int blinkCount=0;

	Map map = null;
	Image image1 = null;
	Image image2 = null;
	Image image3 = null;
	Image Wall = null;
	Image Water = null;
	Image Steel = null;
	Image Grass = null;
	Image Home = null;
	Image Destory = null;
	Image bigTank = null;
	Image back=null;
	public int stagenum;


	public Stage(int stage)
	{   
		
		this.stagenum = stage;
		map = new Map(stage);

		//Initial tank position: random generated.  xBound & yBound
		int x=(int) (Math.random()*Constants.xBound)%390;
		int y=(int) (Math.random()*Constants.yBound)%390;
		OurTank=new Tank(x,y,1,0,0, Constants.tankID);
		OurTank.setMap(map); 
		Tank.allTanks.addElement(OurTank);  

		Events position=new Events("set_position", 0,0,0,Constants.life, Constants.tankID, 0, x, y);

		// for test
		while(!start){ 
			try {
				Thread.sleep(500);
			} catch (InterruptedException e) {
				e.printStackTrace();
			}
			System.out.println("Waiting for start msg...");
		}

		//waiting for the position information  
		while(positionQueue.size()<TankGame.players){ 
			try {
				Thread.sleep(500);
			} catch (InterruptedException e) {
				e.printStackTrace();
			}
			System.out.println("Waiting for position information..");
			Send.sendEvent(position,Constants.port_d);
		}


		//Generate tank. Receive all the InitPositions. 
		while(positionQueue.size()!=0){
			Events tmpPos=positionQueue.remove(0);
			if(tmpPos.getTankID()==OurTank.getId()) continue;
			Tank newTank = new Tank(tmpPos.getX(),tmpPos.getY(),1,0,0, tmpPos.getTankID());
			//newTank.setColor(Color.CYAN);
			//			int color=(int) (Math.random()*4);
			//			switch(color){
			//			case 0:
			//				newTank.setColor(Color.yellow);
			//				break;
			//			case 1:
			//				newTank.setColor(Color.green);
			//				break;
			//			case 2:
			//				newTank.setColor(Color.blue);
			//				break;
			//			case 3:
			//				newTank.setColor(Color.cyan);
			//				break;
			//			case 4:
			//				newTank.setColor(Color.orange);
			//				break;
			//			}

			Tank.allTanks.add(newTank);
			newTank.setMap(map);
		}

		map.removeStuckBlock(Tank.allTanks);
		AudioPlay ad = new AudioPlay("image/start.wav");
		ad.start();
		
		ImageIcon image1Icon = new ImageIcon("image/blast3.gif");
		image1 = image1Icon.getImage();
		ImageIcon image2Icon = new ImageIcon("image/blast2.gif");
		image2 = image2Icon.getImage();
		ImageIcon image3Icon = new ImageIcon("image/blast1.gif");
		image3 = image3Icon.getImage();
		ImageIcon wallIcon = new ImageIcon("image/wall.gif");
		Wall = wallIcon.getImage();
		ImageIcon waterIcon = new ImageIcon("image/water.gif");
		Water = waterIcon.getImage();
		ImageIcon steelIcon = new ImageIcon("image/steel.gif");
		Steel = steelIcon.getImage();
		ImageIcon destoryIcon = new ImageIcon("image/destory.gif");
		Destory = destoryIcon.getImage();
		ImageIcon GrassIcon = new ImageIcon("image/grass.gif");
		Grass = GrassIcon.getImage();
		ImageIcon bigTankIcon=new ImageIcon("image/tank-05.jpg");
		bigTank= bigTankIcon.getImage();
		ImageIcon backIcon=new ImageIcon("image/back.png");
		back= backIcon.getImage();

	}



	public void push_events(int keyCode, String type){
		java.util.Date date= new java.util.Date();	
		Events tmpEvent=new Events(type, date.getTime(), date.getTime(), OurTank.live, OurTank.score,
				Constants.tankID, keyCode, OurTank.getx(), OurTank.gety() );
		Send.sendEvent(tmpEvent,Constants.port_d);//send to the local port.		
	}


	// Major paint events
	public void paint(Graphics g)
	{
		super.paint(g);
		g.fillRect(0,0,Constants.singleWindow,Constants.singleWindow);
		g.drawImage(back,Constants.singleWindow,0,Constants.singleWindow/2,Constants.singleWindow,this);

		g.setColor(Color.gray);
		g.fill3DRect(Constants.singleWindow+15,10,90,60,false);
		
		Font myFont = new Font("TimesRoman",Font.BOLD,20);
		g.setColor(Color.red);
		g.setFont(myFont);
		g.drawString("Life: "+OurTank.live, Constants.singleWindow+20, 35);
		g.drawString("Score: "+OurTank.score, Constants.singleWindow+20, 55);

		if(freeze && blinkCount==0){
			g.drawImage(bigTank,0,0,Constants.singleWindow,Constants.singleWindow,this);
			Font gameFreezed = new Font("TimesRoman",Font.BOLD,20);
			g.setFont(gameFreezed);
			g.setColor(Color.red);
			g.drawString("Delay...",Constants.singleWindow+20,120);
			g.drawString("Delay host will be",Constants.singleWindow+20,140);
			g.drawString("kicked in 15 secs.",Constants.singleWindow+20,160);
			return;
		}
		//drawing small map.
		g.drawString("Minimap",Constants.singleWindow+20,190);
		g.setColor(Color.gray);
		g.fill3DRect(Constants.singleWindow+20,200,75,50,false);
		//g.drawString("B: Exit", Constants.singleWindow+20, 200);
		//g.drawString("N: Next stage", 400, 230);


		//DO not draw tank if outside of the vision.
		for(int i=0; i<Tank.allTanks.size(); i++){
			Tank tank=Tank.allTanks.get(i);
			if(tank==null || tank.live==0) continue; 
			if(getArea(tank.getx(), tank.gety())!=getArea(OurTank.getx(), OurTank.gety())) continue;
			int offsetX=tank.getx();
			int offsetY=tank.gety();
			while(offsetX>=Constants.singleWindow) offsetX-=Constants.singleWindow;
			while(offsetY>=Constants.singleWindow) offsetY-=Constants.singleWindow;

			this.drawTank(offsetX,offsetY, g,tank.getDirect(), tank); 

			for(int j=0; j<tank.ss.size(); j++){
				g.setColor(Color.CYAN);
				Bullet shot = tank.ss.get(j);
				if(getArea(shot.x, shot.y)!=getArea(OurTank.getx(), OurTank.gety())) continue;
				if(shot != null && shot.islive){ 
					int shotX=shot.x;
					int shotY=shot.y;
					while(shotX>=Constants.singleWindow) shotX-=Constants.singleWindow;
					while(shotY>=Constants.singleWindow) shotY-=Constants.singleWindow;
					g.fillRect(shotX,shotY,3,3);
				}
				if(shot.islive == false) tank.ss.remove(shot);
			}
		}
		for(int i=0; i<bomb.size(); i++){
			Bomb newb = bomb.get(i);
			int bX=newb.x-15;
			int bY=newb.y-15;
			if(getArea(bX, bY)!=getArea(OurTank.getx(), OurTank.gety())) continue;
			
			while(bX>=Constants.singleWindow) bX-=Constants.singleWindow;
			while(bY>=Constants.singleWindow) bY-=Constants.singleWindow;
			if(newb.live>6){
				g.drawImage(image1, bX, bY, 30, 30, this);
				newb.bao();
			}else if(newb.live>3){
				g.drawImage(image2, bX, bY, 30, 30, this);
				newb.bao();
			}
			else if(newb.live>0){
				g.drawImage(image3, bX, bY, 30, 30, this);
				newb.bao();
			}
			else {
				bomb.remove(newb);
			} 
		}


		drawMap(g);
		drawSmallMap(g, blinkCount);
		blinkCount=(blinkCount+1)%20;	//Blink

		//Check the Game End condition
		if(OurTank.live <= 0 ){
			Font gameEnd = new Font("TimesRoman",Font.BOLD,50);
			g.setFont(gameEnd);
			g.setColor(Color.white);
			g.drawString("Game Over",80,200);
			//Events result=new Events("gameresult",0,0,0,0,0,0,0, 0);
			//Send.sendEvent(result, Constants.port_d2);
		}
		else if(Tank.allTanks.size()==1 && Tank.allTanks.get(0).getId()==OurTank.getId()){
			Font gameEnd = new Font("TimesRoman",Font.BOLD,50);
			g.setFont(gameEnd);
			g.setColor(Color.white);
			g.drawString("You Win!",80,200);
			//Events result=new Events("gameresult",0,0,0,0,0,0,0, 0);
			//Send.sendEvent(result, Constants.port_d2);
		}
	}


	public void drawTank(int x,int y,Graphics g, int direct, Tank tank)
	{	
		while(x>=Constants.singleWindow) x-=Constants.singleWindow;
		while(y>=Constants.singleWindow) y-=Constants.singleWindow; 

		if(tank.getStatus().equalsIgnoreCase("disconnect")) g.setColor(Color.red);
		if(tank.getId()==OurTank.getId()) g.setColor(Color.YELLOW);
		else g.setColor(Color.CYAN);

		switch(direct){
		case 0:
			g.fill3DRect(x-15,y-15,7,30,false);
			g.fill3DRect(x+8,y-15,7,30,false);
			g.fill3DRect(x-8,y-10,16,20,false);
			g.fill3DRect(x-2,y-15,4,11,true);
			g.fillOval(x-7,y-7,12,12);
			break;
		case 1:
			g.fill3DRect(x-15,y-15,7,30,false);
			g.fill3DRect(x+8,y-15,7,30,false);
			g.fill3DRect(x-8,y-10,16,20,false);
			g.fill3DRect(x-2,y+5,4,11,true);
			g.fillOval(x-7,y-7,12,12);
			break;
		case 2:
			g.fill3DRect(x-15,y-15,30,7,false);
			g.fill3DRect(x-15,y+8,30,7,false);
			g.fill3DRect(x-10,y-8,20,16,false);
			g.fill3DRect(x+5,y-2,11,4,true);
			g.fillOval(x-7,y-7,12,12);
			break;
		case 3:
			g.fill3DRect(x-15,y-15,30,7,false);
			g.fill3DRect(x-15,y+8,30,7,false);
			g.fill3DRect(x-10,y-8,20,16,false);
			g.fill3DRect(x-15,y-2,11,4,true);
			g.fillOval(x-7,y-7,12,12);
			break;
		}
	}


	/* Starting point.
	 * Area I:	Constants.singleWindow+20,200
	 * II:		Constants.singleWindow+45,200
	 * III:		Constants.singleWindow+70,200
	 * IV:		Constants.singleWindow+20,225
	 * V:		Constants.singleWindow+45,225
	 * VI:		Constants.singleWindow+70,225
	 */
	public void drawSmallMap(Graphics g, int count){
		if(count==0){
			g.setColor(Color.BLACK);
			g.fill3DRect(Constants.singleWindow+20,200,75,50,false);
			return;
		}
		//g.fill3DRect(Constants.singleWindow+20,200,75,50,false);
		for(int i=0; i<Tank.allTanks.size(); i++){
			Tank tank=Tank.allTanks.get(i);

			if(tank.getId()==OurTank.getId()) g.setColor(Color.BLUE);
			else  g.setColor(Color.RED); 

			int area=getArea(tank.getx(), tank.gety());
			switch (area){
			case 1:
				g.fill3DRect(Constants.singleWindow+20,200,25,25,true);
				break;
			case 2:
				g.fill3DRect(Constants.singleWindow+45,200,25,25,false);
				break;
			case 3:
				g.fill3DRect(Constants.singleWindow+70,200,25,25,false);
				break;
			case 4:
				g.fill3DRect(Constants.singleWindow+20,225,25,25,false);
				break;
			case 5:
				g.fill3DRect(Constants.singleWindow+45,225,25,25,false);
				break;
			case 6:
				g.fill3DRect(Constants.singleWindow+70,225,25,25,false);
				break;
			}			
		}
	}


	//drawing map. 
	public void drawMap(Graphics g){
		//for(int i=0; i<676; i++){
		for(int i=0; i<Constants.totalBlocks; i++){
			Block hblock = map.block.get(i);
			if(getArea(hblock.x, hblock.y)!=getArea(OurTank.x, OurTank.y)) continue; 

			int offsetX=hblock.x;
			int offsetY=hblock.y;
			while(offsetX>=Constants.singleWindow) offsetX-=Constants.singleWindow;
			while(offsetY>=Constants.singleWindow) offsetY-=Constants.singleWindow;

			switch(hblock.type){
			case 0:
				break;
			case 1: g.drawImage(Wall,offsetX,offsetY,Constants.blocksSize,Constants.blocksSize,this);
			break;
			case 2: g.drawImage(Steel,offsetX,offsetY,Constants.blocksSize,Constants.blocksSize,this);
			break;
			case 3: g.drawImage(Water,offsetX,offsetY,Constants.blocksSize,Constants.blocksSize,this);
			break;
			case 4: g.drawImage(Grass,offsetX,offsetY,Constants.blocksSize,Constants.blocksSize,this);
			}
		}
	} 


	/**
	 * 	User key pressed catch.
	 */
	@Override
	public void keyPressed(KeyEvent e) {
		if(e.getKeyCode() == KeyEvent.VK_SPACE){
			try {
				Runtime server2 = Runtime.getRuntime();
				server2.exec("python client.py " + Constants.tankID);
				System.out.println("@py client triggled.");
				return;
			} catch (Exception e1) {
				e1.printStackTrace();
			}
		}
		push_events(e.getKeyCode(), "keyPressed"); 
	}
	/**
	 * 	User key Released catch.
	 */
	@Override
	public void keyReleased(KeyEvent e) {
	}
	@Override
	public void keyTyped(KeyEvent arg0) {
	}

	@Override
	public void run() {
		while(true)
		{
			//System.out.println("Stage RUN.");
			try {
				Thread.sleep(Constants.timeFrame);
			}  catch(Exception e){
				e.printStackTrace();
			}
			if(freeze){
				this.repaint();
				continue;
			}

			//PULL event
			/* Time frame syncronization: 2 restrictions
			 * 1. event time < time Max
			 * 2. event time < time frame Max
			 */
			//Use the last element here
			if(TankGame.events.isEmpty()) continue;
			long eventTime=TankGame.events.get(TankGame.events.size()-1).getTime();
			long curTime=TankGame.events.get(TankGame.events.size()-1).getTimeFrame();
			long timeMax=(eventTime/Constants.timeFrame)*Constants.timeFrame;
			long timeFrameMax=(curTime/Constants.timeFrame-1)*Constants.timeFrame;
			long timeThreshold=Math.min(timeMax, timeFrameMax);

			int timeSlotTotalRange=0;
			for(int i=0; i<TankGame.events.size(); i++){
				if(TankGame.events.get(i).getTime()<=timeThreshold) timeSlotTotalRange=i;
				else break;
			}

			while(timeSlotTotalRange>=0 && !TankGame.events.isEmpty()){ 
				Events tmpEvent=TankGame.events.get(0);
				// --- Event time > threshold ---  No job should be done.
				long timeSlot=(tmpEvent.getTimeFrame()-tmpEvent.getTime())/Constants.timeFrame;
				int removedElements=pull_events(timeSlot); 
				//Check the tank and bullet hit.
				Stage.allHitCheck();
				this.repaint();
				timeSlotTotalRange-=removedElements;
			}
		}
	}


	/*
	 * Pull all clients' operations in a time slot.
	 * Modify the tank data by the events, and repaint the screen
	 * Sync here: perform events in same timeFrame. 
	 */
	public int pull_events(long timeSlot){
		int removedElements=0;
		TankGame.eventLock.lock();	
		while(!TankGame.events.isEmpty()){
			if((TankGame.events.get(0).getTimeFrame()-TankGame.events.get(0).getTime())/Constants.timeFrame!=timeSlot) break;

			Events tmpEvent=TankGame.events.remove(0);
			removedElements++;

			if(tmpEvent.getName().equalsIgnoreCase("heartbeat")) continue;
			for(int j=0; j<Tank.allTanks.size(); j++){
				//*** Event Action! ***
				if(Tank.allTanks.get(j).getId()==tmpEvent.getTankID()) tmpEvent.EventAction(Tank.allTanks.get(j));
			}
		}
		TankGame.eventLock.unlock();
		return removedElements;
	}


	public static void allHitCheck(){
		for(int i=0; i<Tank.allTanks.size(); i++){
			Tank tank = Tank.allTanks.get(i);
			tank.hitblock();
			tank.hitbullet();
			tank.judgehit();
			tank.tankBulletsMove();
			tank.bulletCounter++;
		}
	}



	/*	Area Map:
	 *   -------------------
	 * 	 |  1  |  2  |  3  |
	 *   -------------------
	 *   |  4  |  5  |  6  |
	 *   ------------------- 
	 */
	public static int getArea(int x, int y){
		if(y<Constants.singleWindow){	//Area 1,2,3
			if(x<Constants.singleWindow) return 1;
			else if(x>=Constants.singleWindow*2) return 3;
			else return 2;
		}
		else{
			if(x<Constants.singleWindow) return 4;
			else if(x>=Constants.singleWindow*2) return 6;
			else return 5;
		}
	}

}
