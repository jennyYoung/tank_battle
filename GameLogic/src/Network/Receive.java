package Network;

import java.io.*;
import java.net.*;
import com.google.gson.Gson;
import com.google.gson.stream.JsonReader;
import Events.Constants;
import Events.Events;
import GUI.Stage;
import GUI.TankGame;


public class Receive implements Runnable{
	//private static final String event.getName() = null;

	private DatagramSocket socket_s;
	//private Vector<Events> eventQueue;
	//private Lock eventLock;
 
	private boolean firstHeart=false;
	
	
	public Receive(){
		//this.eventLock=eventLock;
		//this.eventQueue=eventQueue;
		try {
			socket_s=new DatagramSocket(Constants.port_s);
			
		} catch (SocketException e) {
			e.printStackTrace();
		}

	}

	@Override
	public void run() {
		try {
			System.out.println("listening on port "+Constants.port_s);
			while (true) {
				
				byte[] receiveData = new byte[1024];
				DatagramPacket receivePacket = new DatagramPacket(receiveData, receiveData.length);
				socket_s.receive(receivePacket);  
				String json = new String( receivePacket.getData());
				//System.out.println("Receive:"+json+", from:"+socket_s.getRemoteSocketAddress());
				JsonReader reader = new JsonReader(new StringReader(json));
				reader.setLenient(true);

				Events event = new Gson().fromJson(reader, Events.class);
				boolean found=false;
				if(event.getName().equalsIgnoreCase("set_position")){
					for(int i=0; i<Stage.positionQueue.size();i++){
						if(Stage.positionQueue.get(i).getTankID()==event.getTankID())
							found=true;;
					}
					if(!found){
						Stage.positionQueue.add(event);
						System.out.println("++event Recv: position added. Tankid=" + event.getTankID());
					}
				}
				else if(!firstHeart && event.getName().equalsIgnoreCase("heartbeat")){
					System.out.println("HeartBeat Start.");
					firstHeart=true;
				}
				else if(event.getName().equalsIgnoreCase("start")){
					System.out.println("++event Recv: START.");
					//System.out.println("->Receive:"+json+", from:"+socket_s.getRemoteSocketAddress());
					Stage.start=true;
				}
				if(event.getName().equalsIgnoreCase("unfreeze")){
					if(!event.Forwarded){
						//tank.setStatus("disconnect");
						//tank.setColor(Color.red);
					}
					Stage.freeze=false;
				}
				else if(event.getName().equalsIgnoreCase("freeze")){
					Stage.freeze=true;
				}
				else{ //normal keypress event
//					if(!event.getName().equalsIgnoreCase("heartbeat"))
//						System.out.println("++event Recv: "+json);
					TankGame.eventLock.lock();
					TankGame.events.add(event);
					TankGame.eventLock.unlock();
				}
			}

		} catch(Exception e) {
			System.out.println("Server Closed");
			e.printStackTrace();
		}


	}

}