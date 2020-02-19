import java.io.*;
import java.math.*;
import java.security.*;
import java.text.*;
import java.util.*;
import java.util.concurrent.*;
import java.util.regex.*;

public class Solution {

    private static Map<Integer, Integer> fillMap(int[] arr) {
        Map<Integer, Integer> ans = new HashMap<>();
        for (int i = 0; i < arr.length; i++) {
            if (i != arr[i] - 1)
                ans.put(i, arr[i] - 1);
        }
        return ans;
    }

    // Complete the minimumSwaps function below.
    static int minimumSwaps(int[] arr) {
        int ans = 0;
        Map<Integer, Integer> loops = fillMap(arr);
        while (!loops.isEmpty()) {
            int key = loops.keySet().iterator().next();
            int pointsTo = loops.get(key);
            int nextEdge = loops.get(pointsTo);
            ans++;
//            System.out.printf("a: %d b: %d c: %d\n", key, pointsTo, nextEdge);
            loops.remove(pointsTo);
            if (key == nextEdge) {
                loops.remove(key);
            } else {
                loops.put(key, nextEdge);
            }

        }
        return ans;
    }

    private static final Scanner scanner = new Scanner(System.in);

    public static void main(String[] args) throws IOException {
        BufferedWriter bufferedWriter = new BufferedWriter(new FileWriter(System.getenv("OUTPUT_PATH")));

        int n = scanner.nextInt();
        scanner.skip("(\r\n|[\n\r\u2028\u2029\u0085])?");

        int[] arr = new int[n];

        String[] arrItems = scanner.nextLine().split(" ");
        scanner.skip("(\r\n|[\n\r\u2028\u2029\u0085])?");

        for (int i = 0; i < n; i++) {
            int arrItem = Integer.parseInt(arrItems[i]);
            arr[i] = arrItem;
        }

        int res = minimumSwaps(arr);

        bufferedWriter.write(String.valueOf(res));
        bufferedWriter.newLine();

        bufferedWriter.close();

        scanner.close();
    }
}
