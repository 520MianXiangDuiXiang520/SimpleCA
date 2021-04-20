package chubao.junebao.top;

/**
 * @author JuneBao
 * @date 2021/2/23 18:25
 */
public class Mallard implements BaseDuck {
    @Override
    public void showAppearance() {
        System.out.println("绿头鸭");
    }

    @Override
    public void call() {
        System.out.println("嘎嘎嘎");
    }

    public void fly() {
        System.out.println("绿头鸭会飞");
    }
}
