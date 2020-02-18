import java.io.*;
import java.math.*;
import java.security.*;
import java.text.*;
import java.util.*;
import java.util.concurrent.*;
import java.util.regex.*;

public class Solution {

    private static int SHIFT = 24;
    private static int SIZE = (1 << SHIFT);
    private static long[] tree = new long[(1 << (SHIFT + 1))-1];

    private static void printTree(int levels) {
        int n = 1;
        int it = 0;
        for (int i = 0; i < levels; i++) {
            for (int j = 0; j < n; j++) {
                System.out.printf("%d ", tree[it]);
                ++it;
            }
            System.out.println();
            n *= 2;
        }
    }

    private static int left(int it) {
        return 2 * it + 1;
    }

    private static int right(int it) {
        return 2 * it + 2;
    }

    private static int up(int it) {
        return (it - 1) / 2;
    }

    private static void add(int to, int it, int k, int start, int end) {
//        System.out.printf("to: %d, it: %d, k: %d, start: %d, end: %d\n", to, it, k, start, end);
        if (to >= end - 1) {
            tree[it] += k;
            return;
        }


        int mid = (end + start) / 2;

//        System.out.println("Going left");
        add(to, left(it), k, start, mid);
        if (to >= mid) {
//            System.out.println("Going right");
            add(to, right(it), k, mid, end);
        }
    }

    private static long sumPath(int start) {
        int it = SIZE - 1 + start;
        long ans = 0;
        while (it != 0) {
            ans += tree[it];
            it = up(it);
        }
        return ans;
    }

    // Complete the arrayManipulation function below.
    static long arrayManipulation(int n, int[][] queries) {
        for (int[] query : queries) {
            int a, b, k;
            a = query[0]-1;
            b = query[1]-1;
            k = query[2];
            add(b, 0, k, 0, SIZE);
            if (a > 0) {
                add(a-1, 0, -k, 0, SIZE);
            }
        }

        long ans = 0;

        for (int i = 0; i < n; i++) {
            ans = Math.max(ans, sumPath(i));
        }
        System.out.println(ans);
        return ans;
    }

    private static final Scanner scanner = new Scanner(System.in);

    public static void main(String[] args) throws IOException {
        BufferedWriter bufferedWriter = new BufferedWriter(new FileWriter(System.getenv("OUTPUT_PATH")));

        String[] nm = scanner.nextLine().split(" ");

        int n = Integer.parseInt(nm[0]);

        int m = Integer.parseInt(nm[1]);

        int[][] queries = new int[m][3];

        for (int i = 0; i < m; i++) {
            String[] queriesRowItems = scanner.nextLine().split(" ");
            scanner.skip("(\r\n|[\n\r\u2028\u2029\u0085])?");

            for (int j = 0; j < 3; j++) {
                int queriesItem = Integer.parseInt(queriesRowItems[j]);
                queries[i][j] = queriesItem;
            }
        }

        long result = arrayManipulation(n, queries);

        bufferedWriter.write(String.valueOf(result));
        bufferedWriter.newLine();

        bufferedWriter.close();

        scanner.close();
    }
}
