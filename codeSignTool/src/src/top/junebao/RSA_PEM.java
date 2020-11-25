package top.junebao;

/**
 * @author JuneBao
 * @date 2020/11/25 17:17
 */

import java.math.BigInteger;
import java.security.KeyFactory;
import java.security.interfaces.RSAPrivateKey;
import java.security.interfaces.RSAPublicKey;
import java.security.spec.RSAPrivateKeySpec;
import java.security.spec.RSAPublicKeySpec;
import java.util.Base64;
import java.util.regex.Pattern;

/**
 * RSA PEM格式秘钥对的解析和导出
 *
 * GitHub:https://github.com/xiangyuecn/RSA-java
 *
 * https://github.com/xiangyuecn/RSA-java/blob/master/RSA_PEM.java
 * 移植自：https://github.com/xiangyuecn/RSA-csharp/blob/master/RSA_PEM.cs
 */
public class RSA_PEM {
    /**modulus 模数，公钥、私钥都有**/
    public byte[] Key_Modulus;
    /**publicExponent 公钥指数，公钥、私钥都有**/
    public byte[] Key_Exponent;
    /**privateExponent 私钥指数，只有私钥的时候才有**/
    public byte[] Key_D;

    //以下参数只有私钥才有 https://docs.microsoft.com/zh-cn/dotnet/api/system.security.cryptography.rsaparameters?redirectedfrom=MSDN&view=netframework-4.8
    /**prime1**/
    public byte[] Val_P;
    /**prime2**/
    public byte[] Val_Q;
    /**exponent1**/
    public byte[] Val_DP;
    /**exponent2**/
    public byte[] Val_DQ;
    /**coefficient**/
    public byte[] Val_InverseQ;

    private RSA_PEM() {}


    /**秘钥位数**/
    public int keySize(){
        return Key_Modulus.length*8;
    }

    /**是否包含私钥**/
    public boolean hasPrivate(){
        return Key_D!=null;
    }

    /**得到公钥Java对象**/
    public RSAPublicKey getRSAPublicKey() throws Exception {
        RSAPublicKeySpec spec=new RSAPublicKeySpec(BigX(Key_Modulus), BigX(Key_Exponent));
        KeyFactory factory=KeyFactory.getInstance("RSA");
        return (RSAPublicKey)factory.generatePublic(spec);
    }

    /**得到私钥Java对象**/
    public RSAPrivateKey getRSAPrivateKey() throws Exception {
        if(Key_D==null) {
            throw new Exception("当前为公钥，无法获得私钥");
        }
        RSAPrivateKeySpec spec=new RSAPrivateKeySpec(BigX(Key_Modulus), BigX(Key_D));
        KeyFactory factory=KeyFactory.getInstance("RSA");
        return (RSAPrivateKey)factory.generatePrivate(spec);
    }

    /**转成正整数，如果是负数，需要加前导0转成正整数**/
    static public BigInteger BigX(byte[] bigb) {
        if(bigb[0]<0) {
            byte[] c=new byte[bigb.length+1];
            System.arraycopy(bigb,0,c,1,bigb.length);
            bigb=c;
        }
        return new BigInteger(bigb);
    }

    /**某些密钥参数可能会少一位（32个byte只有31个，目测是密钥生成器的问题，只在c#生成的密钥中发现这种参数，java中生成的密钥没有这种现象），直接修正一下就行；这个问题与BigB有本质区别，不能动BigB**/
    static public byte[] BigL(byte[] bytes, int keyLen) {
        if (keyLen - bytes.length == 1) {
            byte[] c = new byte[bytes.length + 1];
            System.arraycopy(bytes, 0, c, 1, bytes.length);
            bytes = c;
        }
        return bytes;
    }


    /**
     * 用PEM格式密钥对创建RSA，支持PKCS#1、PKCS#8格式的PEM
     */
    static public RSA_PEM FromPEM(String pem) throws Exception {
        RSA_PEM param=new RSA_PEM();

        String base64 = _PEMCode.matcher(pem).replaceAll("");
        //java byte是正负数
        byte[] dataX = Base64.getDecoder().decode(base64);
        if (dataX == null) {
            throw new Exception("PEM内容无效");
        }
        //转成正整数的bytes数组，不然byte是负数难搞
        short[] data=new short[dataX.length];
        for(int i=0;i<dataX.length;i++) {
            data[i]=(short)(dataX[i]&0xff);
        }

        int[] idx = new int[] {0};


        if (pem.contains("PUBLIC KEY")) {

            //读取数据总长度
            readLen(0x30, data, idx);

            //检测PKCS8
            int[] idx2 = new int[] {idx[0]};
            if (eq(_SeqOID, data, idx)) {
                //读取1长度
                readLen(0x03, data, idx);
                //跳过0x00
                idx[0]++;
                //读取2长度
                readLen(0x30, data, idx);
            }else {
                idx = idx2;
            }

            //Modulus
            param.Key_Modulus = readBlock(data, idx);

            //Exponent
            param.Key_Exponent = readBlock(data, idx);
        } else if (pem.contains("PRIVATE KEY")) {
            //读取数据总长度
            readLen(0x30, data, idx);

            //读取版本号
            if (!eq(_Ver, data, idx)) {
                throw new Exception("PEM未知版本");
            }

            //检测PKCS8
            int[] idx2 = new int[] {idx[0]};
            if (eq(_SeqOID, data, idx)) {
                //读取1长度
                readLen(0x04, data, idx);
                //读取2长度
                readLen(0x30, data, idx);

                //读取版本号
                if (!eq(_Ver, data, idx)) {
                    throw new Exception("PEM版本无效");
                }
            } else {
                idx = idx2;
            }

            //读取数据
            param.Key_Modulus = readBlock(data, idx);
            param.Key_Exponent = readBlock(data, idx);
            int keyLen = param.Key_Modulus.length;
            param.Key_D = BigL(readBlock(data, idx), keyLen);
            keyLen = keyLen / 2;
            param.Val_P = BigL(readBlock(data, idx), keyLen);
            param.Val_Q = BigL(readBlock(data, idx), keyLen);
            param.Val_DP = BigL(readBlock(data, idx), keyLen);
            param.Val_DQ = BigL(readBlock(data, idx), keyLen);
            param.Val_InverseQ = BigL(readBlock(data, idx), keyLen);
        } else {
            throw new Exception("pem需要BEGIN END标头");
        }

        return param;
    }
    static private Pattern _PEMCode = Pattern.compile("--+.+?--+|[\\s\\r\\n]+");
    static private byte[] _SeqOID = new byte[] { 0x30, 0x0D, 0x06, 0x09, 0x2A, (byte)0x86, 0x48, (byte)0x86, (byte)0xF7, 0x0D, 0x01, 0x01, 0x01, 0x05, 0x00 };
    static private byte[] _Ver = new byte[] { 0x02, 0x01, 0x00 };

    /**从数组start开始到指定长度复制一份**/
    static private byte[] sub(short[] arr, int start, int count) {
        byte[] val = new byte[count];
        for (int i = 0; i < count; i++) {
            val[i] = (byte)arr[start + i];
        }
        return val;
    }

    /**读取长度**/
    static private int readLen(int first, short[] data, int[] idxO) throws Exception {
        int idx=idxO[0];
        try {
            if (data[idx] == first) {
                idx++;
                if (data[idx] == 0x81) {
                    idx++;
                    return data[idx++];
                } else if (data[idx] == 0x82) {
                    idx++;
                    return (((int)data[idx++]) << 8) + data[idx++];
                } else if (data[idx] < 0x80) {
                    return data[idx++];
                }
            }
            throw new Exception("PEM未能提取到数据");
        }finally {
            idxO[0]=idx;
        }
    }

    /**读取块数据**/
    static private byte[] readBlock(short[] data, int[] idxO) throws Exception {
        int idx=idxO[0];
        try {
            int len = readLen(0x02, data, idxO);
            idx=idxO[0];
            if (data[idx] == 0x00) {
                idx++;
                len--;
            }
            byte[] val = sub(data, idx, len);
            idx += len;
            return val;
        }finally {
            idxO[0]=idx;
        }
    }

    /**比较data从idx位置开始是否是byts内容**/
    static private boolean eq(byte[] byts, short[] data, int[] idxO) {
        int idx=idxO[0];
        try {
            for (int i = 0; i < byts.length; i++, idx++) {
                if (idx >= data.length) {
                    return false;
                }
                if ((byts[i]&0xff) != data[idx]) {
                    return false;
                }
            }
            return true;
        }finally {
            idxO[0]=idx;
        }
    }

    public static void main(String[] args) {
        String key = "-----BEGIN RSA PRIVATE KEY-----\n" +
                "MIIEogIBAAKCAQEAytNu1enUNQGlmYzlQYG/r8hWoubetxf1mazDGL9SnvGjNj7F\n" +
                "3we9lpxT8pGbYhNBh1C2SrwoEDIMy+aVKJIAD1YxkcaRSo7H8Bri9f0zo8ZEwSY2\n" +
                "lEw5n+dFjWuOyiD1yiCKHf074mlOMswcDYFWedOwKVdmspw0GiRqP/9HjIl2C0xv\n" +
                "2i6KtMgGwfKRYdEaanvFyDHxE+PdGF5m/m5+zm1I2XS0WY2RjlIgarK/1uS9Esaj\n" +
                "FfYgG5KipiY5ZW/u7dyDzAih+LlS16cTsuwudj5lb2XX9x/+poka5aAW3YtG8GlV\n" +
                "RACYv+5K9SKqUsOrifhcJxJkRSeA1FmnKRzYsQIDAQABAoIBAEHjkbv4LDHUCSHq\n" +
                "vYccSVMnd82PxoYgSG7VysM9U+/Ce8zhc5JSh2pn+nVwi9O+gakdtTpuCW3JdJLA\n" +
                "o2/8jfxtecjrgsN/wr/jXBuhV6c6f5dnfI+Me6PQk62vZUGQl4hELdo0K8IPh5HE\n" +
                "8NAKVjdZZ37mn7wiNmLPtZx4p5ulb/CYV00fwWeXLRQi9goKq443FNBcKpdtsRZM\n" +
                "ewQX5YfPdqGG2PxMMDK4GRSPyutWx69sRp5npnTMyOTElmUqIhDf2A8CclR6+xco\n" +
                "UkCvPVU5Nuy5kNvUohERdjZIYS7njUpwLhJA8GU2WPQAca1xPMtNoeERYelP07AD\n" +
                "VMlMAoECgYEA76m+9JwkhouDynNKrPKATlznk6k0sDSbikYOI8iYnhKPLO9e/xOz\n" +
                "kRlnat4xvWd2dgE60ikCsb2HnRrDv18VhXB/qzK4oUuGFovd6yIUDWbcuvMauxQu\n" +
                "+YY8ENJTLQqZWcRQwS/7BEjnCz9bL3OcBXZz/7floVrTO1085Ww0XTkCgYEA2Kbb\n" +
                "qlxHpag2o9Je8zZ6yo0mfyHW3CkVxEKUKSJqxNERFljiVWZ+7D6ouLchzm2VJvwk\n" +
                "pWOUpO3EomV9a+MXWV1nILYLeKtO/jlrZiq3A3+cZDcrzmX3OdIfsR8sZ99ku7pZ\n" +
                "7+Q5fly75t6SxMyaKqCCo66+rZNfJ0kw2h2rzzkCgYBZHnPlndJvPZ3qQGj6Wsqf\n" +
                "WSi73eW7yDQ2fMpxP/yQezJGcVSP4ZGaWSn9sVYpqjmAtABdeeaIlYPCRduYZBEq\n" +
                "p6Sx0pCZWe4ooCYLc4alXSSjWBcOjfjRzLq1PqCzVQelO70TuXXMKBfSNOMBiCny\n" +
                "VhPYeVeoYo+9uXQVk+D88QKBgA1VD1WHgj02gc5JBuDOrHXEg+b07ST1PkqqkjWJ\n" +
                "0ao092k5pQv+V7cwD+/2DRWH9tLEV3j6DM6tdxlLR5GZEvnD3rHLoh8V47GPVQWf\n" +
                "gU2sz7H3FzIHYlRjkuGyemgV/jvzNs+lashU6pdFgSCtOpt+7ysleMRzujpPrbru\n" +
                "coE5AoGAASeig26oB2hRLJoeixezxVfi5yjuKgYiXEnY+ebhM31YkPd8tRLgKyz2\n" +
                "sEfGWlUI2Zbg2oSNm5fE/Wg2hsbvNmb7F1g8l1utyefITkgSSH2V9Lu0WSVTI7Zp\n" +
                "DinIO0uJ60Bo9JYGyYbZPo8vnn4P32TfWXcg4JiOnunIcEWwBJI=\n" +
                "-----END RSA PRIVATE KEY-----";
        try {
            RSA_PEM rsa = FromPEM(key);
            RSAPrivateKey pk = rsa.getRSAPrivateKey();

        } catch (Exception e) {
            e.printStackTrace();
        }
    }

}
