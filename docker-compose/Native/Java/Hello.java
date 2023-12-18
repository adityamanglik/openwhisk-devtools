import com.google.gson.JsonObject;
import java.util.Random;
import java.lang.management.ManagementFactory;
import java.lang.management.MemoryMXBean;
import java.lang.management.MemoryUsage;
import java.lang.management.GarbageCollectorMXBean;

public class Hello {

    private static final int ARRAY_SIZE = 10000;

    public static JsonObject main(JsonObject args) {
        long sum = 0;
        long executionTime = 0;
        JsonObject response = new JsonObject();
        response.addProperty("sum", sum);
        response.addProperty("executionTime", executionTime); // Add execution time to response
        response.addProperty("heapInitMemory: ", 0);
        response.addProperty("heapUsedMemory: ", 0);
        response.addProperty("heapCommittedMemory: ", 0);
        response.addProperty("heapMaxMemory: ", 0);
        return response;
    }

}
