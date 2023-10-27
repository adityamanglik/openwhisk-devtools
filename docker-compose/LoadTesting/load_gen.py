from locust import HttpUser, task, between
import numpy as np
import random

class ServerLoadTest(HttpUser):
    # Specify a wait time between tasks
    # wait_time = between(0.1, 1)
    
    # Server URLs
    NATIVE_JAVA_API = "/jsonresponse?seed="
    
    @task
    def send_native_java_request(self):
        # Generate a random number and append to the NATIVE_JAVA_API string
        random_seed = random.randint(0, 1e6)
        request_url = self.NATIVE_JAVA_API + str(random_seed)

        with self.client.get(request_url, catch_response=True) as response:
            if response.status_code == 404:
                response.failure("Got wrong response")
            # TODO: Use time to isolate cold starts
            # elif response.elapsed.total_seconds() > 1.0:
            #     response.failure("Request took too long")

        # response = self.client.get()
        # print(response.json())
        # if self.is_generic_response_correct(response.json()):
        #     response.success()
        # else:
        #     response.failure("Incorrect response structure for Java API")