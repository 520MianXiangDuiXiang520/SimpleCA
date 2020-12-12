package top.junebao;

import java.security.PrivateKey;
import java.security.PublicKey;
import java.util.Base64;

/**
 * @author JuneBao
 * @date 2020/11/25 20:27
 */
public class CodeSign {
    /**
     * 开发者对代码进行签名（开发者使用）
     *   1. 对 sourcePath 中的所有文件进行哈希，
     *   2. 通过开发者私钥进行 RSA 加密和 Base64 编码
     * @param sourcePath 要做签名的代码路径
     * @param privateKey 开发者私钥文件路径
     * @return 返回 Base64 编码后的签名结果
     */
    public static String sign(String sourcePath, String privateKey) throws Exception {
        // 对源代码做哈希
        String h = Hash.getCodeHash(sourcePath);
        System.out.println("源代码得到的哈希值：" + h);
        // 使用私钥加密
        PrivateKey pk = Key.parseThePemPrivateKey(privateKey);
        byte[] res = RSAPemCoder.encryptByPrivateKey(h.getBytes(), pk);
        return RSAPemCoder.encryptBASE64(res);
    }

    /**
     * 验证代码签名（软件使用者用）
     *   1. 从证书中解析出开发者公钥
     *   2. Base64 解码签名，并使用公钥解密出开发者生成的哈希值
     *   3. 重新对自己拿到的代码做哈希
     *   4. 对比自己生成的哈希值与解密出的开发者的哈希值，如果相同，验证通过，反之验证失败
     * @param sourcePath 要验证的代码路径
     * @param certificatePath 证书路径
     * @param signature 开发者生成的代码签名
     * @return 返回布尔值，如果验证通过，返回 true
     */
    public static boolean verification(String sourcePath, String certificatePath, String signature) throws Exception {
        // 获取证书中的公钥
        Credentials c = new Credentials(certificatePath);
        PublicKey pk = c.getPublicKey();
        // 用公钥解密签名
        byte[] sign = RSAPemCoder.decryptBASE64(signature);
        byte[] res = RSAPemCoder.decryptByPublicKey(sign, pk);
        // 重新对代码签名，判断是否一致
        String h = Hash.getCodeHash(sourcePath);
        System.out.println("解密出的签名：" + new String(res));
        System.out.println("重新哈希后的签名：" + h);
        return h.equals(new String(res));
    }

    public static void main(String[] args) {
        String sourcePath = "E:\\桌面文件\\作业\\大四上\\软件实验周\\SimpleCA\\src";
        String cerPath = "E:\\桌面文件\\作业\\大四上\\软件实验周\\SimpleCA\\src\\cers\\17_ZhangSan_251874800.cer";
        String privateKey = "E:\\桌面文件\\作业\\大四上\\软件实验周\\SimpleCA\\codeSignTool\\src\\src\\top\\junebao\\private_key.pem";
        String sign = "";
        try {
            sign = CodeSign.sign(sourcePath, privateKey);
            System.out.println("加密后的代码签名：" + sign);
        } catch (Exception e) {
            e.printStackTrace();
        }

        // -------- 验证 ---------

        String codeSign = "w7c+V0nNzPcM8xkJXvnaD1Q3wZDZj0Jc/k1lsTOPELCrHP8fcU81KxCXsANR7So6u07KuRSeoCLe\n" +
                "9GMCQGVJP9dbF78eMBOIR3zGmadldD8Z8vpYTX/dv5919MO9lEvjB7TShU0TNrFJsE5Glq/WLJtX\n" +
                "wZHRI1ncBYcL7TGa8rbjsCPSilpEBe9F/49aj4beTwtOxO0H1kAB5070SRi29vg2G694eDehUYsa\n" +
                "4OsxettJi+hnzx1awUk5iq+8/u0x54fznMnEp5vIDjogKywxjfsJ2AEt5nYswnYaY2hqwPSGW6Pw\n" +
                "dSvOl7L4Nz5KLahl+qxCMDjsBC+gbSOi2qFLbg==";
        try {
            if (verification(sourcePath, cerPath, codeSign)) {
                System.out.println("签名验证通过！");
            } else {
                System.out.println("签名验证不通过！");
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
