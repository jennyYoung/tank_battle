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


import java.util.Vector;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

import javax.swing.*;
import Events.Constants;
import Events.Events;
import Network.Receive;


/**
 *
 */
public class TankGame extends JFrame implements Runnable
{

	public static Vector<Events> events=new Vector<Events>();
	public static Lock eventLock=new ReentrantLock();
	public static int players;
	public static int playerName;
	private static final long serialVersionUID = 1L;

	int stage = 0;
	StartStage st = null;
	//Help help = null;

	Stage  mp = null;
	Stage  mp1 = null;

	static Process passer;
	static Process server;

	static class shutDownMessage extends Thread {
		public void run() {
			//passer.destroy();
//			server.destroy();
		}
	}

	public static void main(String[] args) 
	{
		if(args.length!=2){ 
			System.out.println("Please give the players number.");
			System.exit(0);
		}
		players=Integer.parseInt(args[1]);
		playerName=Integer.parseInt(args[0]);
		//trigger the python and GO server. 
		try {
			//			passer=Runtime.getRuntime().exec("passer.exe");
			//			System.out.println("@passer start. ");
			//server=Runtime.getRuntime().exec("python bSserver.py "+Constants.tankID+" "+Constants.players);
//			server=Runtime.getRuntime().exec("python bSserver.py "+Constants.tankID+" "+players);
			//server=Runtime.getRuntime().exec("python bSserver.py "+Constants.tankID+" "+3);
//			System.out.println("@bServer start. Sleep for 4 secs..");
			//
//			Runtime.getRuntime().addShutdownHook(new shutDownMessage());
//			Thread.sleep(4000);
			Receive receive=new Receive();
			Thread recv = new Thread(receive);
			recv.start(); 

			//while(true){
			TankGame mtg=new TankGame();
			Thread st = new Thread(mtg);
			st.start();
			//	st.join();
			//}
		} catch ( Exception e) {
			e.printStackTrace();
		}	
	}

	public TankGame()
	{
		st = new  StartStage();
		this.addKeyListener(st); 
		Thread sx = new Thread(st);
		sx.start();
		this.add(st);
		//Panel size.
		this.setSize(Constants.singleWindow+210,Constants.singleWindow+40);
		//this.setSize(Constants.xBound+160,Constants.yBound);
		this.setVisible(true);
		this.setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);
	}
	/**
	 * 
	 */
	public void whichstage(){
		switch(this.stage){
		//
		//		case -1:
		//			this.stage = help.getStagenum();
		//			System.out.println(st.stagenum);
		//			if(this.stage == 0)
		//			{
		//				this.remove(help);
		//				this.add(st); 
		//				this.addKeyListener(st);
		//				//
		//				this.setVisible(true);
		//			}
		//			break;
		case 0:this.stage = st.getStagenum();

		if(this.stage == 1){
			st.stagenum = 0;
			mp =new Stage(1);
			Thread ss = new Thread(mp);
			ss.start();

			//this.addKeyListener(mp);

			this.remove(st);
			this.add(mp); 
			this.addKeyListener(mp);
			this.setVisible(true);
		}
		break;

		}
	} 

	@Override
	public void run() {
		while(true){
			try {
				Thread.sleep(60);
			}  catch(Exception e){
				e.printStackTrace();
			}
			whichstage();
		}
	}

}


