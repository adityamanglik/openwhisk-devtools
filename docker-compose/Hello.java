import com.google.gson.JsonObject;
import java.util.Random;

public class Hello {
    public static JsonObject main(JsonObject args) {
        int seed = 42; // default seed value
        if (args.has("seed")) {
            seed = args.getAsJsonPrimitive("seed").getAsInt();
        }

        Random rand = new Random(seed);
        Integer[] arr = new Integer[1000000];
        long sum = 0;

        for (int i = 0; i < arr.length; i++) {
            arr[i] = rand.nextInt(100000); 
            // populate array with random integers between 0 and 99
        }

        for (int i = 0; i < arr.length; i++) {
            sum += arr[i];
        }

        JsonObject response = new JsonObject();
        response.addProperty("sum", sum);

        return response;
    }
}
