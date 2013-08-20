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
import java.util.Iterator;
import java.util.Vector;

import Events.Constants;
import Events.Events;
import Network.Send;



public class Tank  
{
	int x=0;
	int y=0;
	int speed = 5;
	int direct = 0;
	Color color;
	int score=0;
	public int bulletCounter = 0;
	int id;
	int live = Constants.life;
	Map map = null;
	Vector<Bullet> ss;
	static Vector<Tank> allTanks=new Vector<Tank>();
	
	String status;

	public void setMap(Map map) { this.map = map; }

	public Tank(int x,int y,int speed,int direct, int score, int id)
	{
		this.x=x;
		this.y=y;
		this.id=id;
		this.ss = new Vector<Bullet>();
		this.speed = speed;
		this.direct = direct;
		this.score=score;
		this.status="normal";
	}

	public void shot(){
		if(bulletCounter < 15 || this.ss.size() >=8) return;
		bulletCounter = 0;
		Bullet s= null;
		int area=Stage.getArea(x, y);
		switch(this.direct){
		case 0:
			s = new Bullet(x-1,y-14,0, area);
			break;
		case 1:
			s = new Bullet(x-1,y+14,1, area);
			break;
		case 2:
			s = new Bullet(x+14,y-1,2, area);
			break;
		case 3:
			s = new Bullet(x-14,y-1,3, area);
			break;
		}
		ss.add(s);
		//System.out.println("tank id="+id+", bullet added. ss size="+ss.size()+", ss="+ss+", s="+s);
		AudioPlay jj = new AudioPlay("image/fire.wav");
		jj.start();
	}


	public boolean isTouch(){
		for(int i=0; i<allTanks.size(); i++) {
			Tank tank = allTanks.get(i);
			if(tank.getId()==this.id) continue;	

			if(tank.live>0){
				switch(this.direct){
				case 0:
					if(this.y < tank.y+30+this.speed && this.y > tank.y +30 -this.speed
							&& this.x > tank.x-30 && this.x< tank.x+30)
					{
						return true;
					}
					break;
				case 1:
					if(this.y > tank.y-30-this.speed && this.y < tank.y -30 +this.speed 
							&& this.x > tank.x-30 && this.x< tank.x+30)		
					{
						return true;
					}
					break;
				case 2:
					if(this.y > tank.y-30 && this.y< tank.y+30 
							&& this.x > tank.x-30-this.speed && this.x< tank.x + 30 -this.speed)			
					{
						return true;
					}
					break;
				case 3:
					if(this.y > tank.y-30 && this.y< tank.y+30 
							&& this.x < tank.x+30+this.speed && this.x > tank.x -30 +this.speed)		
					{
						return true;
					}
					break;
				}
			}
		}
		return false;
	}

	/**
	 *  
	 */
	public boolean istouchblock(){
		switch(this.direct){
		case 0:for(int i=0; i<map.block.size(); i++){
			Block hb = map.block.get(i);
			if(hb.type == 0 || hb.type == 4) continue;
			else
				if(hb.x < this.x+14 && hb.x > this.x - 29 && hb.y+15 > this.y-15 -speed && this.y > hb.y)
					return true;
		}
		return false;
		case 1:for(int i=0; i<map.block.size(); i++){
			Block hb = map.block.get(i);
			if(hb.type == 0 || hb.type == 4) continue;
			else
				if(hb.x < this.x+14 && hb.x > this.x - 29 && hb.y < this.y+15 +speed && this.y < hb.y)
					return true;
		}
		return false;
		case 2:for(int i=0; i<map.block.size(); i++){
			Block hb = map.block.get(i);
			if(hb.type == 0 || hb.type == 4) continue;
			else
				if(hb.y < this.y+14 && hb.y > this.y - 29 && hb.x < this.x+15 +speed && this.x < hb.x)
					return true;
		}
		return false;
		case 3:for(int i=0; i<map.block.size(); i++){
			Block hb = map.block.get(i);
			if(hb.type == 0 || hb.type == 4) continue;
			else
				if(hb.y < this.y+14 && hb.y > this.y - 29 
						&& hb.x+15 > this.x-15 -speed && this.x > hb.x)
					return true;
		}
		return false;
		}
		return false;
	}


	protected void hitblock(){	
		for(int i=0; i<map.block.size(); i++){
			Block hb = map.block.get(i);
			
			Iterator<Bullet> iter = ss.iterator();
			//TODO bug fixed. None test yet.
			while(iter.hasNext()){
				Bullet eb = iter.next();
				if(bulletHitBlock(hb, eb)){ 
					System.out.println("block hit: blockX="+hb.x+", Y="+hb.y+",bulletX="+eb.x+",bullety="+eb.y);
					iter.remove();
				}
			}
		}
	}



	protected void hitbullet(){
		{
			for(int i=0; i<allTanks.size(); i++){
				Tank tank = allTanks.get(i);
				if(tank.getId()==this.id) continue;	//skip itself
				if(tank.live<=0) continue;
				for(int j=0; j < tank.ss.size(); j++){					
					for(int n=0; n < this.ss.size(); n++){
						if(tank.ss.get(j).x <= this.ss.get(n).x+3 && tank.ss.get(j).x >= this.ss.get(n).x-3
								&& tank.ss.get(j).y <= this.ss.get(n).y+3 &&
								tank.ss.get(j).y >= this.ss.get(n).y-3){
							tank.ss.get(j).islive = false;
							this.ss.get(n).islive = false;
						}
					}
				}
			}
		}
	}

	public void judgehit(){
		for(int i=0; i<Tank.allTanks.size(); i++){
			Tank tank=Tank.allTanks.get(i);
			if(tank.getId()==this.id) continue;	//skip itself
			if(tank.live<=0 ){
				Tank.allTanks.remove(tank);
				continue;
			}

			for(int j=0; j<this.ss.size(); j++){
				if(this.ss.get(j).islive){ 
					Bullet bullet=this.ss.get(j);
					if(bullet.getX() > tank.getx()-15 && bullet.getX() < tank.getx()+15 
							&& bullet.getY() > tank.gety()-15 && bullet.getY() < tank.gety()+15){
						bullet.islive = false;
						//sync of score and live
						if(tank.live>0){ 
							tank.live-=1;
							Send.sendEvent(new Events("sync", 0, 0, tank.live, tank.score, 
									tank.getId(), 0, 0, 0), Constants.port_d);
						}
						this.score+=1;
						Send.sendEvent(new Events("sync", 0, 0, this.live, this.score, 
								this.id, 0, 0, 0),Constants.port_d);

						AudioPlay jj = new AudioPlay("image/blast.wav");
						jj.start();

						Bomb mybomb = new Bomb(tank.x,tank.y);
						Stage.bomb.add(mybomb);
					}
				}
			}
		}
	}

	/**
	 *  
	 */
	public void moveU(int newPosX, int newPosY){
		if((newPosY-this.speed-Constants.blocksSize)<=0) y=15;
		else y = newPosY- this.speed;
	}
	public void moveD(int newPosX, int newPosY){
		if((newPosY+Constants.blocksSize+this.speed)>=Constants.yBound) 
			y=Constants.yBound-Constants.blocksSize;
		else this.y =newPosY+ this.speed;
	}
	public void moveR(int newPosX, int newPosY){
		if((newPosX+Constants.blocksSize+this.speed)>=Constants.xBound) 
			x=Constants.xBound-Constants.blocksSize;
		else this.x =newPosX+ this.speed;
	}
	public void moveL(int newPosX, int newPosY){
		if((newPosX-this.speed-Constants.blocksSize)<=0) x=15;
		else this.x =newPosX- this.speed;
	}




	public int getScore() {
		return score;
	}
	public void setScore(int score) {
		this.score = score;
	}
	public Color getColor() {
		return color;
	}
	public void setColor(Color color) {
		this.color = color;
	}
	public int getDirect() {
		return direct;
	}
	public void setDirect(int direct) {
		this.direct = direct;
	}
	public int getSpeed() {
		return speed;
	}
	public void setSpeed(int speed) {
		this.speed = speed;
	}
	public int getId() {
		return id;
	}
	public void setId(int id) {
		this.id = id;
	}
	public int getx() { return x; }
	public void setx(int x) { this.x=x; }
	public int gety() { return y; }
	public void sety(int y) { this.y=y; }


	public int getLive() {
		return live;
	}

	public void setLive(int live) {
		this.live = live;
	}

	public String getStatus() {
		return status;
	}

	public void setStatus(String status) {
		this.status = status;
	}

	public void tankBulletsMove(){
		for(int i=0; i<ss.size(); i++){
			Bullet b=ss.get(i);
			if(!b.islive) ss.remove(b);  
			else b.bulletMove();
		}
	}

	public boolean bulletHitBlock(Block hb, Bullet eb){
		if(hb.type == 1){
			if(hb.x-4 <  eb.x && hb.x+14 > eb.x &&	hb.y-4 < eb.y && hb.y+14 > eb.y){
				eb.islive = false;
				hb.type = 0;
				return true;
			}
		}else if(hb.type == 2){
			if(hb.x-4 <  eb.x && hb.x+14 > eb.x &&	hb.y-4 < eb.y && hb.y+14 > eb.y){
				eb.islive = false;
				return true;
			}
		}
		return false;
	}
}

