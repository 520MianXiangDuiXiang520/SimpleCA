package top.junebao;

import java.io.File;
import java.io.FileInputStream;
import java.io.InputStream;
import java.security.PublicKey;
import java.security.cert.CertificateFactory;
import java.security.cert.X509Certificate;
import java.text.SimpleDateFormat;
import java.util.Date;

/**
 * @author JuneBao
 * @date 2020/11/25 16:05
 */
public class Credentials {

    private String version;
    private String serialNumber;

    public String getVersion() {
        return version;
    }

    public String getSerialNumber() {
        return serialNumber;
    }

    public Date getNotBefore() {
        return notBefore;
    }

    public Date getNotAfter() {
        return notAfter;
    }

    public String getSubjectName() {
        return subjectName;
    }

    public String getIssuerName() {
        return issuerName;
    }

    public String getSigAlgName() {
        return sigAlgName;
    }

    public PublicKey getPublicKey() {
        return publicKey;
    }

    private Date notBefore, notAfter;
    private String subjectName, issuerName;
    private String sigAlgName;
    private PublicKey publicKey;

    public void parseTheCertificate(String filePath)
    {
        try
        {
            //读取证书文件
            File file = new File(filePath);
            InputStream inStream = new FileInputStream(file);

            //创建X509工厂类
            CertificateFactory cf = CertificateFactory.getInstance("X.509");

            //创建证书对象
            X509Certificate oCert = (X509Certificate)cf.generateCertificate(inStream);
            inStream.close();

            //获得证书版本
            this.version = String.valueOf(oCert.getVersion());

            //获得证书序列号
            this.serialNumber = oCert.getSerialNumber().toString(16);

            //获得证书有效期
            this.notBefore = oCert.getNotBefore();
            this.notAfter = oCert.getNotAfter();

            //获得证书主体信息
            this.subjectName = oCert.getSubjectDN().getName();

            //获得证书颁发者信息
            this.issuerName = oCert.getIssuerDN().getName();

            //获得证书签名算法名称
            this.sigAlgName = oCert.getSigAlgName();

            PublicKey pk;
            pk = oCert.getPublicKey();
            this.publicKey = pk;
        }
        catch (Exception e)
        {
            System.out.println("解析证书出错！");
        }
    }

    public Credentials(String path) {
        parseTheCertificate(path);
    }

    public static void main(String[] args) {
        Credentials c = new Credentials("E:\\桌面文件\\作业\\大四上\\软件实验周\\SimpleCA\\src\\cers\\17_ZhangSan_251874800.cer");
        System.out.println(c.getPublicKey());
        PublicKey pbuKey = Key.parseThePemPublicKey("E:\\桌面文件\\作业\\大四上\\软件实验周\\SimpleCA\\codeSignTool\\src\\src\\top\\junebao\\public_key.pem");
        System.out.println(pbuKey);
    }

}
