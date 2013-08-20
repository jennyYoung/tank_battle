package Events;

import java.io.BufferedReader;
import java.io.FileOutputStream;
import java.io.InputStreamReader;
import java.io.PrintStream;


public class ExtProcess implements Runnable{
	private String cmd;
	public ExtProcess(String cmd){
		this.cmd=cmd;
	}

	@Override
	public void run() {
		try {
			//System.setOut(new PrintStream(new FileOutputStream("OUT"+Constants.tankID)));
			//System.setErr(new PrintStream(new FileOutputStream("ERR"+Constants.tankID)));
			
			Process proc = Runtime.getRuntime().exec(cmd);

//			BufferedReader stdInput = new BufferedReader(new InputStreamReader(proc.getInputStream()));
//	        BufferedReader stdError = new BufferedReader(new InputStreamReader(proc.getErrorStream()));
//	        Thread.sleep(5000);
//	        String s;
	        //System.out.println("Here is the standard output of the command:\n");
//	        while (( s = stdInput.readLine()) != null) {
//	            System.out.println(s);
//	        }
	        // read any errors from the attempted command
	        //System.out.println("Here is the standard error of the command (if any):\n");
//	        while ((s = stdError.readLine()) != null) {
//	            System.out.println(s);
//	        }
	        
//	        String line;
//			BufferedReader in = new BufferedReader( new InputStreamReader(p.getInputStream()) );
//			while ((line = in.readLine()) != null) {
//				System.out.println(line);
//			}
//			in.close();
		}
		catch (Exception e) {
		}
		
	}

}