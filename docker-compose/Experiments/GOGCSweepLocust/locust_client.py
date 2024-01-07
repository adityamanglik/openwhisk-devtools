from locust import HttpUser, task, between
import random
import locust.stats
import os

locust.stats.CSV_STATS_FLUSH_INTERVAL_SEC = 1

class ServerLoadTest(HttpUser):
    wait_time =  between(0.001, 0.001) # 1000 requests per second

    def on_start(self):
        # Marker for arraysize change
        self.arraysize = 99999
        # Read API from the environment variable
        # self.API = os.getenv("API_URL")
        self.API = "http://128.110.96.59:8180"
        if not self.API:
            print("No API URL provided in environment. Skipping requests.")
            return  # Stop executing if no API URL is set
        self.execution_times_file = open("execution_times.txt", "a")  # File to save execution times
    
    def on_stop(self):
        self.execution_times_file.close()  # Close the file when the test stops

    @task
    def send_request(self):
        if not self.API:
            print("API URL not set. Skipping task.")
            return

        random_seed = random.randint(0, 100)
        request_url = self.API + "/go?seed=" + str(random_seed) + "&arraysize=" + str(self.arraysize)

        with self.client.get(request_url, catch_response=True) as response:
            if response.status_code != 200:
                response.failure(f"Unexpected status code: {response.status_code}")
                # data = response.json()
                # execution_time = data.get("executionTime", "NA")
                # self.execution_times_file.write(str(execution_time) + "\n")
                
        # response = self.client.get()
        # print(response.json())
        # if self.is_generic_response_correct(response.json()):
        #     response.success()
        # else:
        #     response.failure("Incorrect response structure for Java API")