#include <pthread.h>
#include <stdio.h>

// To use gcc with pthreads use lpthread as flag "gcc file.c -lpthread -o outputfilname"

int i = 0;
pthread_t IOThread;


void *func_1(){
	for (int j =0;j<1000000;j++){
		i++;
	}
}

void *func_2(){
	for (int j =0;j<1000000;j++){
		i--;
	}
}

int main(){
	printf("Hello World\n");
	pthread_t thread1;
	pthread_t thread2;
	pthread_create(&thread1,NULL,func_1,NULL);
	pthread_create(&thread2,NULL,func_2,NULL);
	printf("The number i is: %i",i);
	return 0;
}

