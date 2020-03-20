package gseiri;

import java.io.File;
import java.util.ArrayList;
import java.awt.event.KeyEvent;
import javax.swing.JFrame;
import javax.swing.JFileChooser;

@SuppressWarnings("serial")
public class GseiriMain extends JFrame implements java.awt.event.KeyListener {

	public static GseiriMain theApp = null;

	public static void main(String[] args) {
		theApp = new GseiriMain();
		theApp.Run();
	}

	/** 整理対象ディレクトリ */
	File m_target_dir = null;

	/** 整理する画像ファイル一覧 */
	public ArrayList<File> m_imgfiles = null;

	/** 画像表示用パネル */
	MyPanel m_imgbox = null;

	/** 表示エリアの最大横幅 */
	int m_max_width;

	/** 表示エリアの最大縦幅 */
	int m_max_height;

	/**
	 * 
	 */
	private void Run(){
		JFileChooser dlg = new JFileChooser();
		dlg.setDialogTitle("画像ファイルが格納されているフォルダを指定してください");
		dlg.setFileSelectionMode(JFileChooser.DIRECTORIES_ONLY); // ディレクトリのみ選択可能とする
		if(dlg.showOpenDialog(null) == JFileChooser.APPROVE_OPTION) {
			// 指定されたディレクトリ配下の画像ファイル情報をファイルサイズ順にソートして保持
			getImageFilesInfo(dlg.getSelectedFile());
			// ウインドウの準備
			initWindow();
		}
	}

	/**
	 * dir配下の整理対象ファイル一覧をメンバ変数に保存する
	 */
	private void getImageFilesInfo(File dir) {
		m_target_dir = dir;
		File files[] = m_target_dir.listFiles();
		m_imgfiles = new ArrayList<File>();
		for(File file: files){
			if(isImage(file)){
				m_imgfiles.add(file);
			}
		}
		if(m_imgfiles.isEmpty()) {
			showMessage("画像ファイルが見つかりませんでした。");
			System.exit(0);
		}
		m_imgfiles.sort(new ImageSizeComparator());
	}

	/**
	 * 
	 * @param file
	 * @return
	 */
	private boolean isImage(java.io.File file){
		String imgexts[] = {".jpg", ".jpeg", ".png" };
		if(!(file.isFile() && file.canRead())){
			return false;
		}
		int extpos = file.getPath().lastIndexOf(".");
		if(extpos == -1) {
			return false;
		}
		String ext = file.getPath().toLowerCase().substring(extpos);
		for(String imgext: imgexts) {
			if(ext.compareTo(imgext) == 0) {
				return true;
			}
		}
		return false;
	}

	/**
	 * 
	 */
	private void initWindow() {
		java.awt.GraphicsEnvironment env = java.awt.GraphicsEnvironment.getLocalGraphicsEnvironment();
		java.awt.Rectangle desktopsize = env.getMaximumWindowBounds();
		double ratio = 0.7;
		m_max_width = (int)(desktopsize.getWidth() * ratio);
		m_max_height = (int)(desktopsize.getHeight() * ratio);
		setDefaultCloseOperation(DISPOSE_ON_CLOSE);
		setSize(m_max_width, m_max_height);
		m_imgbox = new MyPanel();
		m_imgbox.setMax(m_max_width, m_max_height);
		add(m_imgbox);
		setTitle("gseiri");
		addKeyListener(this);
		m_imgbox.setImage(m_imgfiles.get(0));
		paint();
		setVisible(true);
	}
	
	/**
	 * 
	 * 
	 */
	public void keyTyped(KeyEvent e){
		char k = e.getKeyChar();
		if(('a' <= k && k <= 'z') || ('A' <= k && k <= 'Z')) {
			String newdirs = m_target_dir + File.separator + k;
			File newdirf = new File(newdirs);
			if(!newdirf.exists()) {
				if(!newdirf.mkdir()) {
					showMessage("ディレクトリを作成できませんでした");
					System.exit(0);
				}
			} else if(!newdirf.isDirectory()) {
				showMessage("ディレクトリではないファイルがあります");
				System.exit(0);
			}
			String newfiles = newdirs + File.separator + m_imgfiles.get(0).getName();
			File newfilef = new File(newfiles);
			if(newfilef.exists()) {
				showMessage("ファイル「" + newfiles + "」はすでにありました");
				System.exit(0);
			}
			m_imgfiles.get(0).renameTo(newfilef);
			m_imgfiles.remove(0);
			if(m_imgfiles.isEmpty()) {
				System.exit(0);
			}
			if(!m_imgbox.setImage(m_imgfiles.get(0))) {
				showMessage("ファイル読み込みに失敗した");
				System.exit(0);
			}
			paint();
		}
	}
	
	public void paint() {
		getContentPane().setPreferredSize(new java.awt.Dimension(m_imgbox.getDspWidth(), m_imgbox.getDspHeight()));
		pack();
		m_imgbox.repaint();
	}

	//
	public void keyPressed(KeyEvent e) {}

	//
	public void keyReleased(KeyEvent e) {}

	//
	private void showMessage(String s){
		javax.swing.JOptionPane.showMessageDialog(null, s);
	}

	private class MyPanel extends javax.swing.JPanel {

		private java.awt.image.BufferedImage m_image = null;
		private int m_max_width = -1;
		private int m_max_height = -1;
		private int m_img_width = -1;
		private int m_img_height = -1;
		private int m_dsp_width = -1;
		private int m_dsp_height = -1;
		
		public void paintComponent(java.awt.Graphics g) {
			super.paintComponent(g);
			if(m_image == null) {
				return;
			}
			java.awt.Graphics2D g2d = (java.awt.Graphics2D) g;
			g2d.drawImage(m_image,
					0, 0, m_dsp_width, m_dsp_height,
					0, 0, m_img_width, m_img_height, null);
		}
		
		public void setMax(int max_width, int max_height) {
			m_max_width = max_width;
			m_max_height = max_height;
		}
		
		public boolean setImage(File f) {
			try {
				m_image = javax.imageio.ImageIO.read(f);
			} catch(java.io.IOException e) {
				return false;
			}
			// 画像の大きさと最大値の比が１を超えていたら最大値を超えているので縮小処理を行う
			m_img_width = m_image.getWidth();
			m_img_height = m_image.getHeight();
			m_dsp_width = m_img_width;
			m_dsp_height = m_img_height;
			double width_ratio = (double)m_dsp_width / (double)m_max_width;
			double height_ratio = (double)m_dsp_height / (double)m_max_height;
			double ratio = width_ratio > height_ratio ? width_ratio : height_ratio;
			if(ratio > 1.0) {
				m_dsp_width = (int)(m_dsp_width / ratio);
				m_dsp_height = (int)(m_dsp_height / ratio);
			}
			return true;
		}
		
		public int getDspWidth() {
			return m_dsp_width;
		}
		
		public int getDspHeight() {
			return m_dsp_height;
		}
				
	}
	
	/**
	 * 画像ファイル表示順序のソート用。
	 * @author oumit
	 *
	 */
	private class ImageSizeComparator implements java.util.Comparator<File> {
		public int compare(File a, File b){
			if(a.length() > b.length()) {
				return -1;
			}
			if(a.length() < b.length()) {
				return 1;
			}
			return 0;
		}
	}

	
}
