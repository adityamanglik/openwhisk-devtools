import com.google.gson.JsonObject;
import java.util.Random;

public class Hello {

    // Introduced a constant for the array size
    private static final int ARRAY_SIZE = 1000000;

    public static JsonObject main(JsonObject args) {
        int seed = 42; // default seed value
        if (args.has("seed")) {
            seed = args.getAsJsonPrimitive("seed").getAsInt();
        }

        Random rand = new Random(seed);
        Integer[] arr = new Integer[ARRAY_SIZE];
        long sum = 0;

        for (int i = 0; i < arr.length; i++) {
            arr[i] = rand.nextInt(100000); 
            // populate array with random integers between 0 and (100000 - 1)
        }

        for (int i = 0; i < arr.length; i++) {
            sum += arr[i];
        }

        JsonObject response = new JsonObject();
        response.addProperty("sum", sum);

        return response;
    }
}
