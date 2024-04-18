
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/mman.h>
#include <string.h>

#define SIZE (60l * 1024l * 1024l * 1024l) // 60GB

int main() {
    // Allocate memory
    void *ptr = mmap(NULL, SIZE, PROT_READ | PROT_WRITE, MAP_PRIVATE | MAP_ANONYMOUS, -1, 0);
    if (ptr == MAP_FAILED) {
        perror("mmap");
        exit(EXIT_FAILURE);
    }

    // Advise kernel to use Huge Pages for the allocated memory
    if (madvise(ptr, SIZE, MADV_HUGEPAGE) == -1) {
        perror("madvise");
        exit(EXIT_FAILURE);
    }

    // Set all bytes of the allocated memory to 'A'
    memset(ptr, 'A', SIZE);
    printf("Press any Key\n");
    getchar();

    // Free memory
    if (munmap(ptr, SIZE) == -1) {
        perror("munmap");
        exit(EXIT_FAILURE);
    }

    return 0;
}