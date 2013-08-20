package Events;

import java.awt.event.KeyEvent;

import GUI.Tank;


public class Events{

	String Name;
	long Time;
	long TimeFrame;
	int Tank_life;
	int Tank_score;
	int Tank_id;
	int Key_code;
	int X;
	int Y;
	public boolean Forwarded;

	public Events(String name, long time, long timeFrame, int tank_life, int tank_score, int tank_id, int key_code, int position_x, int position_y) {
		super();
		this.Name = name;
		this.Forwarded = true;
		this.TimeFrame=timeFrame;
		this.Time = time;
		this.Tank_life=tank_life;
		this.Tank_score=tank_score;
		this.Tank_id = tank_id;
		this.Key_code = key_code;
		this.X = position_x;
		this.Y = position_y;
	}

	public int getX() {
		return X;
	}

	public void setX(int x) {
		this.X = x;
	}

	public int getY() {
		return Y;
	}

	public void setY(int y) {
		this.Y = y;
	}

	public String getName() {
		return Name;
	}
	public void setName(String name) {
		this.Name = name;
	}
	public long getTime() {
		return Time;
	}
	public void setTime(long time) {
		this.Time = time;
	}
	public int getTankID() {
		return Tank_id;
	}
	public void setTankID(int tankID) {
		this.Tank_id = tankID;
	}
	public int getKeyCode() {
		return Key_code;
	}
	public void setKeyCode(int keyCode) {
		this.Key_code = keyCode;
	}


	public long getTimeFrame() {
		return TimeFrame;
	}

	public void setTimeFrame(long timeFrame) {
		TimeFrame = timeFrame;
	}

	public boolean isForwarded() {
		return Forwarded;
	}

	public void setForwarded(boolean forwarded) {
		Forwarded = forwarded;
	}

	public void EventAction(Tank tank){

		if(Name.equalsIgnoreCase("sync")){
			//life SYNC!!
			tank.setLive(this.Tank_life);
			tank.setScore(this.Tank_score);
		}
		
		else if(Name.equalsIgnoreCase("keyPressed")){
			tank.setStatus("normal");
			
			switch (Key_code){
			case KeyEvent.VK_S:
				tank.setDirect(1);
				for(int i=0;i<8;i++){
					if(tank.istouchblock() == false && !tank.isTouch())
						tank.moveD(X, Y+i);
				}
				break;
			case KeyEvent.VK_W:
				tank.setDirect(0);
				for(int i=0;i<8;i++){
					if(tank.istouchblock() == false && !tank.isTouch())
						tank.moveU(X, Y-i);
				}
				break;
			case KeyEvent.VK_D:
				tank.setDirect(2);
				for(int i=0;i<8;i++){
					if(tank.istouchblock() == false && !tank.isTouch())
						tank.moveR(X+i, Y);
				}
				break;
			case KeyEvent.VK_A:
				tank.setDirect(3);
				for(int i=0;i<8;i++){
					if(tank.istouchblock() == false && !tank.isTouch())
						tank.moveL(X-i, Y);
				}
				break;
			case KeyEvent.VK_J:
				tank.shot();
				break;
			}
		}
	}

}


