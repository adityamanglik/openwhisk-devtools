import com.sun.net.httpserver.HttpServer;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpExchange;
import com.google.gson.JsonObject; // import for JsonObject
import java.net.InetSocketAddress;
import java.io.OutputStream;
import java.io.IOException; // import for IOException
import java.util.Map;
import java.util.HashMap;

public class JsonServer {

    public static void main(String[] args) throws Exception {
        HttpServer server = HttpServer.create(new InetSocketAddress(8600), 0);
        server.createContext("/jsonresponse", new JsonHandler());
        server.setExecutor(null); // creates a default executor
        server.start();
    }

    static class JsonHandler implements HttpHandler {
    @Override
    public void handle(HttpExchange exchange) throws IOException {
        Map<String, String> params = queryToMap(exchange.getRequestURI().getQuery());
        JsonObject args = new JsonObject();

        if (params.containsKey("seed")) {
            args.addProperty("seed", Integer.parseInt(params.get("seed")));
        }

        if (params.containsKey("arraysize")) {
            args.addProperty("arraysize", Integer.parseInt(params.get("arraysize")));
        }

        if (params.containsKey("requestnumber")) {
            args.addProperty("requestnumber", Integer.parseInt(params.get("requestnumber")));
        }

        JsonObject response = Hello.main(args);
        String jsonResponse = response.toString();

        // Send the response
        exchange.sendResponseHeaders(200, jsonResponse.length());
        OutputStream os = exchange.getResponseBody();
        os.write(jsonResponse.getBytes());
        os.close();
        }
    }

    private static Map<String, String> queryToMap(String query) {
    Map<String, String> result = new HashMap<>();
    if (query != null) {
        for (String param : query.split("&")) {
            String[] entry = param.split("=");
            if (entry.length > 1) {
                result.put(entry[0], entry[1]);
            }
        }
    }
    return result;
    }
}
