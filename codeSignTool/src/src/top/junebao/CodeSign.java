package top.junebao;

import java.security.PrivateKey;
import java.security.PublicKey;
import java.util.Base64;

/**
 * @author JuneBao
 * @date 2020/11/25 20:27
 */
public class CodeSign {
    public static String sign(String sourcePath, String privateKey) throws Exception {
        // 对源代码做哈希
        String h = Hash.getCodeHash(sourcePath);
        System.out.println(h);
        // 使用私钥加密
        PrivateKey pk = Key.parseThePemPrivateKey(privateKey);
        byte[] res = RSAPemCoder.encryptByPrivateKey(h.getBytes(), pk);
        return RSAPemCoder.encryptBASE64(res);
    }

    public static boolean verification(String sourcePath, String certificatePath, String signature) throws Exception {
        // 获取证书中的公钥
        Credentials c = new Credentials(certificatePath);
        PublicKey pk = c.getPublicKey();
        // 用公钥解密签名
        byte[] sign = RSAPemCoder.decryptBASE64(signature);
        byte[] res = RSAPemCoder.decryptByPublicKey(sign, pk);
        // 重新对代码签名，判断是否一致
        String h = Hash.getCodeHash(sourcePath);
        System.out.println(new String(res));
        System.out.println(h);
        return h.equals(new String(res));
    }

    public static void main(String[] args) {
        String sourcePath = "E:\\桌面文件\\作业\\大四上\\软件实验周\\SimpleCA";
        String cerPath = "E:\\桌面文件\\作业\\大四上\\软件实验周\\SimpleCA\\src\\cers\\17_ZhangSan_251874800.cer";
        String privateKey = "E:\\桌面文件\\作业\\大四上\\软件实验周\\SimpleCA\\codeSignTool\\src\\src\\top\\junebao\\private_key.pem";
        String sign = "";
        try {
            sign = CodeSign.sign(sourcePath, privateKey);
            System.out.println();
        } catch (Exception e) {
            e.printStackTrace();
        }

        try {
            if (CodeSign.verification(sourcePath, cerPath, sign)) {
                System.out.println("签名验证通过");
            } else {
                System.out.println("签名验证未通过！");
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
