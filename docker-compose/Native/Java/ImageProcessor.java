import com.google.gson.JsonObject;
import javax.imageio.ImageIO;
import java.awt.*;
import java.awt.image.BufferedImage;
import java.io.File;
import java.io.IOException;
import java.util.Random;
import java.util.List;
import java.util.ArrayList;
import java.lang.management.ManagementFactory;
import java.lang.management.MemoryMXBean;
import java.lang.management.MemoryUsage;
import java.lang.management.GarbageCollectorMXBean;

public class ImageProcessor {

    private static final String TMP = "/tmp/";
    private static final String[] fileNames = {
        "Resources/img1.jpg",
        "Resources/img2.jpg"
    };

    public static JsonObject processImage(JsonObject args) throws IOException {
        int seed = 42; // default seed value
        int ARRAY_SIZE = 10000; // default arraysize value
        int request_number = Integer.MAX_VALUE; // default arraysize value
        if (args.has("seed")) {
            seed = args.getAsJsonPrimitive("seed").getAsInt();
        }

        if (args.has("arraysize")) {
            ARRAY_SIZE = args.getAsJsonPrimitive("arraysize").getAsInt();
        }

        if (args.has("requestnumber")) {
            request_number = args.getAsJsonPrimitive("requestnumber").getAsInt();
        }

        long startTime = System.nanoTime(); // Start time tracking
        long sum = 0;

        // Select image based on seed
        Random rand = new Random(seed);
        String selectedFile = fileNames[rand.nextInt(fileNames.length)];
        BufferedImage image = ImageIO.read(new File(selectedFile));

        // Apply transformations and sum pixels after each step
        image = resize(image, ARRAY_SIZE);
        sum += sumPixels(image);

        // Add random seed to every pixel
        for (int y = 0; y < image.getHeight(); y++) {
            for (int x = 0; x < image.getWidth(); x++) {
                Color color = new Color(image.getRGB(x, y));
                int r = clamp(color.getRed() + rand.nextInt(256));
                int g = clamp(color.getGreen() + rand.nextInt(256));
                int b = clamp(color.getBlue() + rand.nextInt(256));
                int newColor = new Color(r, g, b).getRGB();
                image.setRGB(x, y, newColor);
            }
        }
        sum += sumPixels(image);
        
        image = flipHorizontally(image);
        sum += sumPixels(image);

        // image = flipVertically(image);
        // sum += sumPixels(image);

        image = rotate(image, 90);
        sum += sumPixels(image);

        long executionTime = System.nanoTime() - startTime; // Calculate execution time
        executionTime = executionTime / 1000;

        // Create JSON response
        JsonObject response = new JsonObject();
        response.addProperty("sum", sum);
        response.addProperty("seed", seed);
        response.addProperty("arraysize", ARRAY_SIZE);
        response.addProperty("executionTime", executionTime);

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
        response.addProperty("heapInitMemory", initMemory);
        response.addProperty("heapUsedMemory", usedMemory);
        response.addProperty("heapCommittedMemory", committedMemory);
        response.addProperty("heapMaxMemory", maxMemory);
        response.addProperty("request", request_number);
        return response;
    }

    private static BufferedImage resize(BufferedImage image, int arraySize) {
        Image tmp = image.getScaledInstance(arraySize, arraySize, Image.SCALE_SMOOTH);
        BufferedImage resized = new BufferedImage(arraySize, arraySize, BufferedImage.TYPE_INT_ARGB);
        Graphics2D g2d = resized.createGraphics();
        g2d.drawImage(tmp, 0, 0, null);
        g2d.dispose();
        return resized;
    }

    private static BufferedImage flipHorizontally(BufferedImage image) {
        int width = image.getWidth();
        int height = image.getHeight();
        BufferedImage flipped = new BufferedImage(width, height, image.getType());
        Graphics2D g = flipped.createGraphics();
        g.drawImage(image, 0, 0, width, height, width, 0, 0, height, null);
        g.dispose();
        return flipped;
    }

    private static BufferedImage flipVertically(BufferedImage image) {
        int width = image.getWidth();
        int height = image.getHeight();
        BufferedImage flipped = new BufferedImage(width, height, image.getType());
        Graphics2D g = flipped.createGraphics();
        g.drawImage(image, 0, 0, width, height, 0, height, width, 0, null);
        g.dispose();
        return flipped;
    }

    private static BufferedImage rotate(BufferedImage image, int angle) {
        int width = image.getWidth();
        int height = image.getHeight();
        BufferedImage rotated = new BufferedImage(width, height, image.getType());
        Graphics2D g = rotated.createGraphics();
        g.rotate(Math.toRadians(angle), width / 2, height / 2);
        g.drawImage(image, null, 0, 0);
        g.dispose();
        return rotated;
    }

    private static long sumPixels(BufferedImage image) {
        long sum = 0;
        for (int y = 0; y < image.getHeight(); y++) {
            for (int x = 0; x < image.getWidth(); x++) {
                Color color = new Color(image.getRGB(x, y));
                sum += color.getRed() + color.getGreen() + color.getBlue();
            }
        }
        return sum;
    }

    private static int clamp(int value) {
        return Math.max(0, Math.min(255, value));
    }
}