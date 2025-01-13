from http.server import BaseHTTPRequestHandler, HTTPServer
import resource
import json
import random
import time
import psutil
from urllib.parse import urlparse, parse_qs

# Additional imports for Chameleon rendering
import six
from chameleon import PageTemplate

PORT = 9900

# Same Chameleon ZPT template that was in lambda_handler
BIGTABLE_ZPT = """\
<table xmlns="http://www.w3.org/1999/xhtml"
       xmlns:tal="http://xml.zope.org/namespaces/tal">
  <tr tal:repeat="row python: options['table']">
    <td tal:repeat="c python: row.values()">
      <span tal:define="d python: c + 1"
            tal:attributes="class python: 'column-' + %s(d)"
            tal:content="python: d" />
    </td>
  </tr>
</table>""" % six.text_type.__name__


def main_logic(seed, array_size, req_num):
    """
    Main logic function that:
      1) Uses 'seed' as num_of_cols
      2) Uses 'array_size' as num_of_rows
      3) Renders a large HTML table using Chameleon
      4) Returns the length of the rendered HTML as 'sum'
    """

    random.seed(seed)  # Ensure reproducibility with a given seed
    sum_val = 0

    # Start the timer
    start_time = time.perf_counter()

    # -----------------------------------------
    # 1) Prepare the data the same way as in lambda_handler
    num_of_rows = array_size
    num_of_cols = 10

    # Create a Chameleon PageTemplate
    tmpl = PageTemplate(BIGTABLE_ZPT)

    # Build one row's data
    row_data = {}
    for i in range(num_of_cols):
        row_data[str(i)] = i

    # Create the full table (list of identical rows, for demonstration)
    table = [row_data for _ in range(num_of_rows)]

    # Render the table
    rendered_html = tmpl.render(options={'table': table})

    # Store the length of the rendered HTML in sum_val
    sum_val = len(rendered_html)
    # -----------------------------------------

    end_time = time.perf_counter()
    duration_seconds = end_time - start_time
    duration_microseconds = duration_seconds * 1_000_000

    # Memory usage (optional; might slow things down if done frequently)
    process = psutil.Process()
    memory_info = process.memory_info()
    memory_full_info = process.memory_full_info()

    # Return response in the same format as your baseline
    response = {
        "sum": sum_val,  # Now the length of the rendered HTML
        "executionTime": duration_microseconds,
        "requestNumber": req_num,
        "arraysize": array_size,
        "usedHeapSize": memory_full_info.uss,  # Unique set size
        "totalHeapSize": memory_info.vms,
        # Optionally include some or all of the rendered HTML
        # "data": rendered_html  # Uncomment if you want the full HTML in the response
    }
    return response


class RequestHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        parsed_path = urlparse(self.path)
        path = parsed_path.path
        query_components = parse_qs(parsed_path.query)

        # Defaults
        seed = 100
        array_size = 100
        req_num = 2**53 - 1  # Just as in your original code

        if 'seed' in query_components:
            seed = int(query_components['seed'][0])
        if 'arraysize' in query_components:
            array_size = int(query_components['arraysize'][0])
        if 'requestnumber' in query_components:
            req_num = int(query_components['requestnumber'][0])

        if path.startswith("/Python"):
            response = main_logic(seed, array_size, req_num)
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            self.wfile.write(bytes(json.dumps(response), "utf8"))
        else:
            self.send_response(404)
            self.end_headers()


if __name__ == "__main__":
    # Optional: limit_memory(512 * 1024 * 1024)  # e.g., 512 MB
    server = HTTPServer(('0.0.0.0', PORT), RequestHandler)
    print("Server running on port", PORT)
    server.serve_forever()
