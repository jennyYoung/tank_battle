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

import java.io.File;
import java.io.IOException;

import javax.sound.sampled.AudioFormat;
import javax.sound.sampled.AudioInputStream;
import javax.sound.sampled.AudioSystem;
import javax.sound.sampled.DataLine;
import javax.sound.sampled.SourceDataLine;



/** 
 *
 */
public class AudioPlay extends Thread{
	private String filename ;
	public AudioPlay(String wavfile)
	{
		filename = wavfile;
		
	}
	
	public void run ()
	{
		File soundFile = new File(filename);
		AudioInputStream audioInputStream = null;
		try {
			audioInputStream = AudioSystem.getAudioInputStream(soundFile);
			
		} catch (Exception e) {
			e.printStackTrace();
			return ;
		}
		
		AudioFormat format = audioInputStream.getFormat();
		SourceDataLine auline = null;
		DataLine.Info info =new DataLine.Info(SourceDataLine.class, format);
		
		try {
			auline = (SourceDataLine) AudioSystem.getLine(info);
			auline.open(format);
		} catch (Exception e) { 
			e.printStackTrace();
			return;
		}
		
		auline.start();
		int nByteRead =0;
		byte []abData = new byte[1024];
		
		try {
			while (nByteRead!=-1)
			{
				nByteRead = audioInputStream.read(abData,0,abData.length);
				if (nByteRead>=0)
				{
					auline.write(abData, 0, nByteRead);
				}
			}
		} catch (IOException e) {
			e.printStackTrace();
			return;
		}finally
		{
			auline.drain();
			auline.close();
		}
		
	}
}
