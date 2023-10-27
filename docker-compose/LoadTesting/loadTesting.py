from locust import HttpUser, task, between
import paramiko
import time
import os

class ServerLoadTest(HttpUser):
    wait_time = between(1, 2)

    def ssh_execute(self, command):
        """
        Execute an SSH command on the OW_SERVER_NODE.
        Returns the command's stdout.
        """
        client = paramiko.SSHClient()
        client.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        client.connect('OW_SERVER_NODE', username='am_CU', password='your_password')  # Use the appropriate authentication method
        stdin, stdout, stderr = client.exec_command(command)
        output = stdout.read().decode('utf-8')
        client.close()
        return output

    @task
    def warm_up_and_test(self):
        self.kill_java_server()
        Xmx_values = ["64m", "128m", "256m", "512m", "1g", "2g", "4g"]
        MaxGCPauseMillis_values = [50, 100, 150, 200, 250, 300]

        for current_Xmx in Xmx_values:
            for current_MaxGCPauseMillis in MaxGCPauseMillis_values:
                self.start_server(current_Xmx, current_MaxGCPauseMillis)
                self.warm_up()
                # Now, run the actual test using `load_gen.py` and store results
                self.run_test()
                self.kill_java_server()

    def kill_java_server(self):
        pid_command = "jps | awk '/JsonServer/ {print $1}'"
        pid = self.ssh_execute(pid_command).strip()
        if pid:
            kill_command = f"kill {pid}"
            self.ssh_execute(kill_command)
            print(f"Killed JsonServer with PID {pid}.")
        else:
            print("JsonServer is not running.")

    def start_server(self, xmx_value, max_gc_pause_millis):
        GC_FLAGS = f"-Xmx{xmx_value} -XX:MaxGCPauseMillis={max_gc_pause_millis} -XX:+PrintGCDetails -XX:+PrintGCDateStamps -Xloggc:/users/am_CU/openwhisk-devtools/docker-compose/PureJava/gc_log_{xmx_value}_{max_gc_pause_millis}"
        start_command = f"cd /users/am_CU/openwhisk-devtools/docker-compose/PureJava/; taskset -c 1 java -cp .:gson-2.10.1.jar {GC_FLAGS} JsonServer > /users/am_CU/openwhisk-devtools/docker-compose/PureJava/server_log 2>&1 &"
        self.ssh_execute(start_command)

    def warm_up(self):
        while True:
            response = self.client.get("/jsonresponse")
            if response.status_code == 200 and response.text.startswith("{"):
                break
            time.sleep(1)

    def run_test(self):
        # Assuming load_gen.py is locally executable and available
        output = os.popen("python load_gen.py").read()
        median_throughput = [line.split()[-1] for line in output.splitlines() if "Median Throughput:" in line][0]
        print(f"Calculated median throughput: {median_throughput}")
