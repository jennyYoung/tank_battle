package Network;

import java.io.*;
import java.net.*;

import com.google.gson.Gson;
import Events.Events;


public class Send{
	public Send(){}
	public static void sendEvent(Events event, int port){
		try {
			DatagramSocket clientSocket = new DatagramSocket();
			InetAddress IPAddress = InetAddress.getByName("localhost");
			Gson gson = new Gson();
			String json = gson.toJson(event);  
			
			//System.out.println("+event Send:" + json);
			byte[] sendData = new byte[1024];
			sendData = json.getBytes();

			DatagramPacket sendPacket = new DatagramPacket(sendData, sendData.length, IPAddress, port); 
			clientSocket.send(sendPacket);	
			//System.out.println("+event Send: tankId=" + event.getTankID()+", time=" + event.getTimestamp()
			//		+", name="+event.getName());
		} catch(IOException ioe) {
			System.out.println("sendEvent(): passer connection error.");
		}
	}

}