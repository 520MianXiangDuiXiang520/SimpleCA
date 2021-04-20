package chubao.junebao.top;

/**
 * @author JuneBao
 * @date 2021/2/23 18:28
 */
public class RubberDuck implements BaseDuck {
    @Override
    public void showAppearance() {
        System.out.println("橡皮鸭");
    }

    @Override
    public void call() {
        System.out.println("吱吱吱");
    }
}
