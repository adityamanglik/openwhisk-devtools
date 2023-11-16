import com.google.gson.JsonObject;
import java.util.ArrayList;
import java.util.List;
import java.util.Random;
import java.io.FileWriter;
import java.io.IOException;
import java.lang.management.ManagementFactory;
import java.lang.management.MemoryMXBean;
import java.lang.management.MemoryUsage;
import java.lang.management.GarbageCollectorMXBean;

public class Hello {

    private static final int ARRAY_SIZE = 0;
    private static final String FILE_NAME = "Java_execution_times.txt";
    private static final List<Long> executionTimes = new ArrayList<>();

    static {
        // Add shutdown hook
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            saveExecutionTimesToFile();
        }));
    }

    private static synchronized void saveExecutionTimeToFile(long executionTime) {
        executionTimes.add(executionTime); // Add execution time to the list
    }

    private static void saveExecutionTimesToFile() {
        try (FileWriter writer = new FileWriter(FILE_NAME, true)) { // true for append mode
            for (Long time : executionTimes) {
                writer.write(time + "\n");
            }
        } catch (IOException e) {
            System.err.println("Error writing execution times to file: " + e.getMessage());
        }
    }

    public static JsonObject main(JsonObject args) {
        int seed = 42; // default seed value
        if (args.has("seed")) {
            seed = args.getAsJsonPrimitive("seed").getAsInt();
        }
        // Start time tracking
        long startTime = System.currentTimeMillis();

        Random rand = new Random(seed);
        Integer[] arr = new Integer[ARRAY_SIZE];
        long sum = 0;

        for (int i = 0; i < arr.length; i++) {
            arr[i] = rand.nextInt(100000); // populate array with random integers between 0 and (100000 - 1)
        }

        for (int i = 0; i < arr.length; i++) {
            sum += arr[i];
        }
        long executionTime = endTime - startTime;
        saveExecutionTimeToFile(executionTime); // Save each execution time

        JsonObject response = new JsonObject();
        response.addProperty("sum", sum);

        // Garbage collector information
        int gcIndex = 0;
        for (GarbageCollectorMXBean gc : ManagementFactory.getGarbageCollectorMXBeans()) {
            gcIndex++; // To distinguish between different GC objects
            long count = gc.getCollectionCount();
            if (count != -1) { // -1 if the collection count is not available
                response.addProperty("gc" + gcIndex + "CollectionCount", count);
            }
            long timer = gc.getCollectionTime();
            if (timer != -1) { // -1 if the collection time is not available
                response.addProperty("gc" + gcIndex + "CollectionTime", timer);
            }
        }

        MemoryMXBean memoryMXBean = ManagementFactory.getMemoryMXBean();

        MemoryUsage heapMemoryUsage = memoryMXBean.getHeapMemoryUsage();

        long initMemory = heapMemoryUsage.getInit();  // Initial memory amount
        long usedMemory = heapMemoryUsage.getUsed();  // Amount of used memory
        long committedMemory = heapMemoryUsage.getCommitted(); // Amount of memory guaranteed to be available for the JVM
        long maxMemory = heapMemoryUsage.getMax();   // Maximum amount of memory (can change over time, can be undefined)
        response.addProperty("heapInitMemory: ", initMemory);
        response.addProperty("heapUsedMemory: ", usedMemory);
        response.addProperty("heapCommittedMemory: ", committedMemory);
        response.addProperty("heapMaxMemory: ", maxMemory);
        
        return response;
    }
}
