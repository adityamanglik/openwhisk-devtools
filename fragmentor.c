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
    void *ptr[ITERS][NUM_PAGES];
    for(int i = 0; i < ITERS; i++ )
    {
        printf("ITERATION NUMBER: %d\n", i);
        for(int j = 0; j < NUM_PAGES; j++)
        {
            ptr[i][j] = mmap(NULL, PAGE_SIZE, PROT_READ | PROT_WRITE, MAP_PRIVATE | MAP_ANONYMOUS | MAP_HUGETLB, -1, 0);
            if(ptr[i][j] == MAP_FAILED) {
                perror("mmap");
                exit(EXIT_FAILURE);
            }

            fill_with_random_alphabets(ptr[i][j], 1);
        }

        printf("Allocated %d huge pages.\n", NUM_PAGES);
        // getchar();
        srand(time(NULL));

        for(int j = 0; j < NUM_PAGES/2; j++)
        {
            int index = rand() % NUM_PAGES;
            if(ptr[i][j] != NULL) {
                munmap(ptr[i][j], PAGE_SIZE);
                ptr[i][j] = NULL;
            } else {
                j--;
            }
        }
        printf("Deallocated %d huge pages randomly.\n", NUM_PAGES / 2);
        // getchar();

    }
    getchar();
    // Free remaining pages
    for(int i = 0; i < ITERS; i++)
    {
        for(int j = 0; j < NUM_PAGES; j++)
        {
            if(ptr[i][j] != NULL) {
                munmap(ptr[i][j], PAGE_SIZE);
            }
        }
    }

    printf("All huge pages deallocated.\n");
    return 0;
}