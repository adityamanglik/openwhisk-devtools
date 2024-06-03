from locust import HttpUser, task, events, LoadTestShape, constant_pacing
import random
# from locust_plugins import constant_total_ips
import time
import threading

class ServerLoadTest(HttpUser):
    # wait_time = constant_pacing(1)

    def on_start(self):
        self.API = "http://node0:8801/JS"
        # self.execution_times_file = open("execution_times.txt", "a")

    # def on_stop(self):
    #     self.execution_times_file.close()

    @task
    def send_request(self):
        arraysize = 10000
        requestnumber = random.randint(0, 10000)
        random_seed = random.randint(0, 10000)
        request_url = self.API + "?seed=" + str(random_seed) + "&arraysize=" + str(arraysize) + "&requestnumber=" + str(requestnumber)

        with self.client.get(request_url, catch_response=True) as response:
            if response.status_code != 200:
                # continue
                # data = response.json()
                # execution_time = data.get("executionTime", "NA")
                # self.execution_times_file.write(str(execution_time) + "\n")
            # else:
                response.failure(f"Unexpected status code: {response.status_code}")
        # Each user makes two requests
        # with self.client.get(request_url, catch_response=True) as response:
        #     if response.status_code == 200:
        #         data = response.json()
        #         execution_time = data.get("executionTime", "NA")
        #         self.execution_times_file.write(str(execution_time) + "\n")
        #     else:
        #         response.failure(f"Unexpected status code: {response.status_code}")
        
response_times = []

@events.request_success.add_listener
def on_request_success(request_type, name, response_time, response_length, **kwargs):
    response_times.append(response_time)

@events.request_failure.add_listener
def on_request_failure(request_type, name, response_time, exception, **kwargs):
    response_times.append(response_time)

@events.quitting.add_listener
def on_quitting(environment, **kwargs):
    response_times.sort()
    if response_times:
        p99999 = response_times[int(len(response_times) * 0.99999)]
        print(f"99.999th percentile response time: {p99999} ms")
    else:
        print("No requests were made during the test.")

# class CustomLoadShape(LoadTestShape):
#     time_limit = 10  # Test duration in seconds
#     spawn_rate = 5  # Users to spawn per second
#     initial_users = 10  # Starting number of users
#     user_increment = 10  # Users to add every interval
#     interval = 1  # Interval (in seconds) to increase users

#     def tick(self):
#         run_time = self.get_run_time()

#         if run_time < self.time_limit:
#             # Check if it's time to increase the user count
#             if int(run_time) % self.interval == 0:
#                 # Increase the user count
#                 self.initial_users += self.user_increment
#             return (self.initial_users, self.spawn_rate)
#         return None
