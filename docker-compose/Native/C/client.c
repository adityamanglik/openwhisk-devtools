#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <netdb.h>

#define PORT 9100
#define BUFFER_SIZE 1024

int main() {
    int sockfd, bytes_received;
    char buffer[BUFFER_SIZE];
    struct sockaddr_in serv_addr;
    struct hostent *server;

    // Create socket
    sockfd = socket(AF_INET, SOCK_STREAM, 0);
    if (sockfd < 0) {
        perror("ERROR opening socket");
        exit(EXIT_FAILURE);
    }

    // Get server information
    server = gethostbyname("node0");
    printf("%s\n", server->h_name);
    if (server == NULL) {
        fprintf(stderr, "ERROR, no such host\n");
        exit(EXIT_FAILURE);
    }

    // Clear server address structure
    memset(&serv_addr, 0, sizeof(serv_addr));
    serv_addr.sin_family = AF_INET;
    memcpy(&serv_addr.sin_addr.s_addr, server->h_addr, server->h_length);
    serv_addr.sin_port = htons(PORT);

    // Connect to server
    if (connect(sockfd, (struct sockaddr *)&serv_addr, sizeof(serv_addr)) < 0) {
        perror("ERROR connecting");
        exit(EXIT_FAILURE);
    }

    // Send HTTP request
    char *request = "GET /Python HTTP/1.1\r\nHost: localhost\r\n\r\n";
    if (send(sockfd, request, strlen(request), 0) < 0) {
        perror("ERROR sending request");
        exit(EXIT_FAILURE);
    }

    // Receive response from server
    //printf("Response from server:\n");
    while ((bytes_received = recv(sockfd, buffer, BUFFER_SIZE - 1, 0)) > 0) {
         if (bytes_received < 0) {
            perror("ERROR receiving response");
            exit(EXIT_FAILURE);
        }
        buffer[bytes_received] = '\0';
        printf("%s\n", buffer);
    }
    // Close socket
    close(sockfd);

    return 0;
}