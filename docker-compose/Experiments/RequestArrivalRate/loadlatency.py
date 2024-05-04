from locust import HttpUser, task, events, LoadTestShape, constant_pacing
import random
# from locust_plugins import constant_total_ips
import time
import threading

class ServerLoadTest(HttpUser):
    wait_time = constant_pacing(1)

    def on_start(self):
        self.API = "http://node0:8601/jsonresponse"
        self.execution_times_file = open("execution_times.txt", "a")

    def on_stop(self):
        self.execution_times_file.close()

    @task
    def send_request(self):
        arraysize = 10000
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
        # Each user makes two requests
        # with self.client.get(request_url, catch_response=True) as response:
        #     if response.status_code == 200:
        #         data = response.json()
        #         execution_time = data.get("executionTime", "NA")
        #         self.execution_times_file.write(str(execution_time) + "\n")
        #     else:
        #         response.failure(f"Unexpected status code: {response.status_code}")

# class CustomLoadShape(LoadTestShape):
#     time_limit = 240  # Test duration in seconds
#     spawn_rate = 0  # Users to spawn per second
#     initial_users = 1  # Starting number of users
#     user_increment = 0  # Users to add every interval
#     interval = 100  # Interval (in seconds) to increase users

#     def tick(self):
#         run_time = self.get_run_time()

#         if run_time < self.time_limit:
#             # Check if it's time to increase the user count
#             if int(run_time) % self.interval == 0:
#                 # Increase the user count
#                 self.initial_users += self.user_increment
#             return (self.initial_users, self.spawn_rate)
#         return None

# # Start the custom load shape
# events.test_start.add_listener(lambda _: threading.Thread(target=CustomLoadShape().tick).start())
