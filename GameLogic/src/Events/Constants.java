package Events;


public class Constants{

	public static final int tankID=(int) (Math.random()*Integer.MAX_VALUE);
	public final static int timeFrame=50;
	
	public final static int port_s=8888;
	public final static int port_d=9999;
	public final static int port_d2=51425;
	
	public static final int singleWindow=390;//26*15
	public static final int singleWindowBlocks=676;//26*26
	public static final int xBound=390*3;
	public static final int yBound=390*2;	
	public static final int xFogBound=390/2;
	public static final int yFogBound=390/2;
	public static final int totalBlocks=26*26*3*2;
	public static final int xBlocks=26*3;
	public static final int yBlocks=26*2;
	public static final int blocksSize=15;
	
	public static final int ROAD=0;
	public static final int WALL=1;
	public static final int CONCRETE=2;
	public static final int RIVER=3;
	public static final int JUNGLE=4;
	
	public static final int life=10;
	//public static final int players=3;
}