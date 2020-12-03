package top.junebao;

import cn.hutool.crypto.digest.DigestAlgorithm;

import java.io.File;
import java.util.ArrayList;
import java.util.List;
import java.util.Objects;

/**
 * @author JuneBao
 * @date 2020/11/25 20:34
 * 用于对所有代码进行 Hash
 */
public class Hash {

    public static void findFileList(File dir, List<String> fileNames) {
        if (!dir.exists() || !dir.isDirectory()) {
            return;
        }
        String[] files = dir.list();
        for (int i = 0; i < Objects.requireNonNull(files).length; i++) {
            File file = new File(dir, files[i]);
            if (file.isFile()) {
                fileNames.add(dir + "\\" + file.getName());
            } else {
                findFileList(file, fileNames);
            }
        }
    }

    public static String getCodeHash(String sourcePath) {
        String codeHash = "";
        // 文件名列表
        List<String> files = new ArrayList<>();
        findFileList(new File(sourcePath), files);

        for (String value :  files) {
            cn.hutool.core.io.file.FileReader fileReader = new cn.hutool.core.io.file.FileReader(value);
            String result = fileReader.readString();

            cn.hutool.crypto.digest.Digester md5 = new cn.hutool.crypto.digest.Digester(DigestAlgorithm.SHA256);
            codeHash = md5.digestHex(codeHash + result);
        }
        return codeHash;
    }

    public static void main(String[] args) {
        String h = Hash.getCodeHash("E:\\桌面文件\\作业\\大四上\\软件实验周\\SimpleCA");
        System.out.println(h);
    }
}
