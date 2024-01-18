from locust import HttpUser, task, between, events
import random
import locust.stats
import os

locust.stats.CSV_STATS_FLUSH_INTERVAL_SEC = 1

class ServerLoadTest(HttpUser):
    # wait_time = between(0.1, 1)

    def on_start(self):
        # Read API from the environment variable
        self.API = "http://128.110.96.59:8180/go"
        # if not self.API:
            # print("No API URL provided in environment. Skipping requests.")
            # return  # Stop executing if no API URL is set
        self.execution_times_file = open("execution_times.txt", "a")  # File to save execution times
    
    def on_stop(self):
        self.execution_times_file.close()  # Close the file when the test stops

    @task
    def send_request(self):
        # if not self.API:
            # print("API URL not set. Skipping task.")
            # return
        arraysize = 100000
        requestnumber = random.randint(0, 10000)
        random_seed = random.randint(0, 10000)
        request_url = self.API + "?seed=" + str(random_seed) + "&arraysize=" + str(arraysize) + "&requestnumber=" + str(requestnumber)

        with self.client.get(request_url, catch_response=True) as response:
            if response.status_code == 200:
                data = response.json()
                execution_time = data.get("executionTime", "NA")
                self.execution_times_file.write(str(execution_time) + "\n")
            else:
                response.failure(f"Unexpected status code: {response.status_code}")

            # TODO: Use time to isolate cold starts
            # elif response.elapsed.total_seconds() > 1.0:
            #     response.failure("Request took too long")

        # response = self.client.get()
        # print(response.json())
        # if self.is_generic_response_correct(response.json()):
        #     response.success()
        # else:
        #     response.failure("Incorrect response structure for Java API")