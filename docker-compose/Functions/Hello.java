import com.google.gson.JsonObject;
import java.util.Random;
import java.lang.management.ManagementFactory;
import java.lang.management.GarbageCollectorMXBean;

public class Hello {

    private static final int ARRAY_SIZE = 750000;

    public static JsonObject main(JsonObject args) {
        int seed = 42; // default seed value
        if (args.has("seed")) {
            seed = args.getAsJsonPrimitive("seed").getAsInt();
        }

        Random rand = new Random(seed);
        Integer[] arr = new Integer[ARRAY_SIZE];
        long sum = 0;

        for (int i = 0; i < arr.length; i++) {
            arr[i] = rand.nextInt(100000); // populate array with random integers between 0 and (100000 - 1)
        }

        for (int i = 0; i < arr.length; i++) {
            sum += arr[i];
        }

        JsonObject response = new JsonObject();
        response.addProperty("sum", sum);

        // Adding JVM details
        // response.addProperty("jvmLocation", System.getProperty("java.home"));
        // response.addProperty("jvmInfo", System.getProperty("java.vm.info"));
        // response.addProperty("jvmCommand", System.getProperty("sun.java.command"));

        // Garbage collector information
        long totalCollectionCount = 0;
        long totalCollectionTime = 0;
        long totalCollectionBeans = 0;
        for (GarbageCollectorMXBean gc : ManagementFactory.getGarbageCollectorMXBeans()) {
            totalCollectionBeans += 1;
            long count = gc.getCollectionCount();
            if (count != -1) { // -1 if the collection count is not available
                totalCollectionCount += count;
            }
            long timer = gc.getCollectionTime();
            if (timer != -1) { // -1 if the collection count is not available
                totalCollectionTime += timer;
            }
        }
        response.addProperty("gcTotalCollectionCount", totalCollectionCount);
        response.addProperty("gcTotalCollectionTime", totalCollectionTime);
        response.addProperty("gcTotalCollectors", totalCollectionBeans);
        
        return response;
    }
}
