package top.junebao;

import javax.crypto.Cipher;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.interfaces.RSAPrivateKey;
import java.security.interfaces.RSAPublicKey;
import java.util.Base64;

/**
 * @author JuneBao
 * @date 2020/11/25 16:40
 */
public class Key {
    private byte[] privateKey;

    public byte[] getPrivateKey() {
        return privateKey;
    }

    public Key(String privateKeyPath) {
        parseThePemPrivateKey(privateKeyPath);
    }

    public static PrivateKey parseThePemPrivateKey(String filePath) {
        cn.hutool.core.io.file.FileReader fileReader = new cn.hutool.core.io.file.FileReader(filePath);
        String result = (fileReader).readString();
        RSA_PEM rsa = null;
        try {
            rsa = RSA_PEM.FromPEM(result);
        } catch (Exception e) {
            e.printStackTrace();
        }
        try {
            assert rsa != null;
            RSAPrivateKey pk = rsa.getRSAPrivateKey();
            return pk;
        } catch (Exception e) {
            e.printStackTrace();
        }
        return null;
    }

    public static PublicKey parseThePemPublicKey(String filePath) {
        cn.hutool.core.io.file.FileReader fileReader = new cn.hutool.core.io.file.FileReader(filePath);
        String result = (fileReader).readString();
        RSA_PEM rsa = null;
        try {
            rsa = RSA_PEM.FromPEM(result);
        } catch (Exception e) {
            e.printStackTrace();
        }
        try {
            assert rsa != null;
            RSAPublicKey pk = rsa.getRSAPublicKey();
            return pk;
        } catch (Exception e) {
            e.printStackTrace();
        }
        return null;
    }

    public static void main(String[] args) {
        Key k = new Key("E:\\桌面文件\\作业\\大四上\\软件实验周\\SimpleCA\\codeSignTool\\src\\src\\top\\junebao\\private_key.pem");
        System.out.println(k.privateKey);
    }
}
