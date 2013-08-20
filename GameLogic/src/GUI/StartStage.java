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
//import java.io.IOException;

import javax.swing.ImageIcon;
import javax.swing.JPanel;

import Events.Constants;


public class StartStage extends JPanel implements Runnable,KeyListener{

	private static final long serialVersionUID = 1L;
	int fontsize = 20;
	int i = 1;
	int stagenum = 0;
//	private int x = -15;
//	private int judge = 0;
//	private int type = -1;
	
	public void paint(Graphics g){
		super.paint(g);
		g.fillRect(0, 0, Constants.singleWindow, Constants.singleWindow);
		//g.fillRect(0, 0, Constants.xBound, Constants.yBound);
		ImageIcon backIcon=new ImageIcon("image/back.png");
		Image back= backIcon.getImage();
		g.drawImage(back,Constants.singleWindow,0,Constants.singleWindow/2,Constants.singleWindow,this);
		
		Font myFont = new Font("TimesRoman",Font.BOLD,fontsize);
		g.setFont(myFont);
		g.setColor(Color.white);
		g.drawString("Distributed TANK GAME: Team 22", 40, 150);
		
		Font seFont = new Font("TimesRoman",Font.BOLD,15);
		g.setFont(seFont);
		g.drawString("SPACE: Login Game", 120,190);
		g.drawString("Move the Tank: w, a, s, d", 120,210);
		g.drawString("Shot: j", 120,230);

}
	public void drawTank(int x,int y,Graphics g,int type)
	{
	    switch(type){
		   case -1:
			   g.setColor(Color.CYAN);
		       break;
		   case 1:
			   g.setColor(Color.yellow);
		       break;
	    }
			   g.fill3DRect(x-15,y-15,30,7,false);
	           g.fill3DRect(x-15,y+8,30,7,false);
	           g.fill3DRect(x-10,y-8,20,16,false);
	           g.fill3DRect(x+5,y-2,11,4,true);
	           g.fillOval(x-7,y-7,12,12);
	}

	public int getStagenum() {
		return stagenum;
	}
	public void setStagenum(int stagenum) {
		this.stagenum = stagenum;
	}

	public void run() {
		while(true){
			try{
				Thread.sleep(30);
			}catch ( Exception e){
				e.printStackTrace();
			}
			repaint();
		}
	}

	@Override
	public void keyPressed(KeyEvent e) {
		if(e.getKeyCode() == KeyEvent.VK_SPACE){
				stagenum = 1;	
				// for test
			try {
				Runtime server2 = Runtime.getRuntime();
				//server2.exec("python client.py " + Constants.tankID);
				server2.exec("python client.py " + TankGame.playerName);
				System.out.println("@py client start.");
			} catch (Exception e1) {
				e1.printStackTrace();
			}
		}
//		else if(e.getKeyCode() == KeyEvent.VK_H)
//			stagenum = -1;	
	}

	@Override
	public void keyReleased(KeyEvent e) {
	}

	@Override
	public void keyTyped(KeyEvent e) {
	}
}
