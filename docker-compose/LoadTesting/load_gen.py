from locust import HttpUser, task, between, events
import random
import locust.stats
import sys

locust.stats.CSV_STATS_FLUSH_INTERVAL_SEC = 1

class ServerLoadTest(HttpUser):
    # wait_time = between(0.1, 1)

    def on_start(self):
        # Read API from the first tag if available
        if self.environment.tags:
            self.API = self.environment.tags[0]
            print("Running load test on API: " + self.API)
        else:
            print("No API URL provided. Exiting.")
            sys.exit(1)

    @task
    def send_request(self):
        random_seed = random.randint(0, 10000)
        request_url = self.API + "?seed=" + str(random_seed)

        with self.client.get(request_url, catch_response=True) as response:
            if response.status_code == 404:
                response.failure("404 error code")

            # TODO: Use time to isolate cold starts
            # elif response.elapsed.total_seconds() > 1.0:
            #     response.failure("Request took too long")

        # response = self.client.get()
        # print(response.json())
        # if self.is_generic_response_correct(response.json()):
        #     response.success()
        # else:
        #     response.failure("Incorrect response structure for Java API")