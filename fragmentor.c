#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/mman.h>
#include <time.h>

#define ITERS 2
#define NUM_PAGES 1024 * 21
#define PAGE_SIZE (2 * 1024 * 1024) // 2 MB

void fill_with_random_alphabets(char *ptr, size_t size) {
    for (size_t i = 0; i < size; ++i) {
        ptr[i] = 'A' + rand() % 26; // Random alphabet
    }
}

int main() {
    // Allocate huge pages
    int i = 0;
    while(1) {
        void *ptr[NUM_PAGES];
        
        printf("ITERATION NUMBER: %d\n", i++);
        for(int j = 0; j < NUM_PAGES; j++)
        {
            ptr[j] = mmap(NULL, PAGE_SIZE, PROT_READ | PROT_WRITE, MAP_PRIVATE | MAP_ANONYMOUS, -1, 0);
            if(ptr[j] == MAP_FAILED) {
                perror("mmap");
                exit(EXIT_FAILURE);
            }

            fill_with_random_alphabets(ptr[j], PAGE_SIZE);
        }
    
        for(int j = 0; j < NUM_PAGES; j++)
        {
            if(ptr[j] != NULL) {
                munmap(ptr[j], PAGE_SIZE);
            }
        }
        printf("All huge pages deallocated.\n");
    }

    
    return 0;
}