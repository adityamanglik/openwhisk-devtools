#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/mman.h>
#include <sys/stat.h>
#include <time.h>
#include <regex.h>

#define PORT 9100
#define REQUEST_NUMBER 0

void initialize_data_sparse(char *mem, size_t size) {
    // Define the pattern to be written sparsely
    char pattern = '\x01';  // Example pattern

    // Write the pattern sparsely into the memory
    for (size_t i = 0; i < size; i += 1024 * 1024) {  // Write 1 byte per megabyte
        mem[i] = pattern;
    }
}

char *allocate_huge_pages(size_t size_gb) {
    // Calculate the size in bytes
    size_t size_bytes = size_gb * (1024 * 1024 * 1024);  // 1 GB = 1024^3 bytes

    // Open a temporary file to back the mmap object
    int fd = open("/dev/zero", O_RDWR);
    if (fd == -1) {
        perror("open");
        exit(EXIT_FAILURE);
    }

    // Create a memory-mapped file using huge pages
    char *mem = mmap(NULL, size_bytes, PROT_READ | PROT_WRITE, MAP_SHARED, fd, 0);
    if (mem == MAP_FAILED) {
        perror("mmap");
        exit(EXIT_FAILURE);
    }
    close(fd);

    return mem;
}

char *get_thp_status() {
    FILE *fp = popen("sudo cat /sys/kernel/mm/transparent_hugepage/enabled", "r");
    if (!fp) {
        perror("popen");
        exit(EXIT_FAILURE);
    }

    char buffer[256];
    if (!fgets(buffer, sizeof(buffer), fp)) {
        perror("fgets");
        exit(EXIT_FAILURE);
    }
    pclose(fp);

    // Define the regular expression pattern
    char *pattern = "\\[(.*?)\\]";

    // Compile the regular expression
    regex_t regex;
    if (regcomp(&regex, pattern, REG_EXTENDED) != 0) {
        fprintf(stderr, "Failed to compile regex\n");
        exit(EXIT_FAILURE);
    }

    // Execute the regular expression
    regmatch_t match[2];
    if (regexec(&regex, buffer, 2, match, 0) != 0) {
        // No match found
        return strdup("unknown");
    }

    // Extract the status from the match
    int start = match[1].rm_so;
    int end = match[1].rm_eo;
    int length = end - start;
    char *status = (char *)malloc(length + 1);
    if (!status) {
        perror("malloc");
        exit(EXIT_FAILURE);
    }
    strncpy(status, buffer + start, length);
    status[length] = '\0';

    // Free the regex resources
    regfree(&regex);

    return status;
}

int get_nr_anon_thp() {
    FILE *fp = popen("sudo egrep nr_anon_transparent_hugepages /proc/vmstat", "r");
    if (!fp) {
        perror("popen");
        exit(EXIT_FAILURE);
    }

    int nr_thp;
    if (fscanf(fp, "%*s %d", &nr_thp) != 1) {
        perror("fscanf");
        exit(EXIT_FAILURE);
    }
    pclose(fp);

    return nr_thp;
}

char *mainLogic() {
    clock_t start_time = clock();
    char *huge_mem = allocate_huge_pages(20);
    initialize_data_sparse(huge_mem, 20l * (1024l * 1024l * 1024l));
    munmap(huge_mem, 20l * (1024l * 1024l * 1024l));
    clock_t end_time = clock();
    double duration_microseconds = ((double)(end_time - start_time)) * 1e6 / CLOCKS_PER_SEC;

    char *thp_status = get_thp_status();
    int nr_thp = get_nr_anon_thp();

    char *response = malloc(256);
    if (!response) {
        perror("malloc");
        exit(EXIT_FAILURE);
    }
    sprintf(response, "{\"state\": \"finished\", \"exec_time\": %.2f, \"request_number\": %d, \"thp_status\": \"%s\", \"nr_thp\": %d}", 
            duration_microseconds, REQUEST_NUMBER + 1, thp_status, nr_thp);

    free(thp_status);
    return response;
}

int main() {
    int server_fd, new_socket, valread;
    struct sockaddr_in address;
    int addrlen = sizeof(address);

    char buffer[1024] = {0};
    char *response;

    if ((server_fd = socket(AF_INET, SOCK_STREAM, 0)) == 0) {
        perror("socket failed");
        exit(EXIT_FAILURE);
    }

    address.sin_family = AF_INET;
    address.sin_addr.s_addr = INADDR_ANY;
    address.sin_port = htons(PORT);

    if (bind(server_fd, (struct sockaddr *)&address, sizeof(address)) < 0) {
        perror("bind failed");
        exit(EXIT_FAILURE);
    }

    if (listen(server_fd, 3) < 0) {
        perror("listen");
        exit(EXIT_FAILURE);
    }
    while (1) {        
        if ((new_socket = accept(server_fd, (struct sockaddr *)&address, (socklen_t*)&addrlen)) < 0) {
            perror("accept");
            exit(EXIT_FAILURE);
        }
        printf("Hello\n");

        response = mainLogic();
        printf("response: %s\n", response);
        send(new_socket, response, strlen(response), 0);
        close(new_socket);
        free(response);
    }
    return 0;
}