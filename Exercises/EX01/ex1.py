#!/bin/python

from threading import Thread

i=0

def func1 ():
        global i
        for j in range (0,100000):
            i+=1

def func2 ():
        global i
        for j in range (1000000,0):
            i-=1

def main():
    thread1= Thread(target=func1,args=(),)
    thread2= Thread(target=func2,args=(),)
    thread1.start()
    thread2.start()
    thread1.join()
    thread2.join()
    print("finished this stupid program")
    print(i)

main()
