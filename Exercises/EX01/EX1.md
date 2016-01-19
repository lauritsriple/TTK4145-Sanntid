What is concurrency?
In computer science concurrency means that multiple computations are executing at the same time/simultaneously. The different computationts might communicate or interact with each other.

What is parallelism?
Operations that happen independent and simultaneously 

Whats the difference between concurrency and parallelism?
Concurrency is used to run multiple operations that needs to use the same resources. For example you are able to open multiple programs on your computer at once. In parallelism on the other hand, the operations happens physically simultaneously in hardware. Concurrency might also use parallell processing.

Why have machines become increasingly multicore in the past decade?
They have become increasingly multicore in the past decade to make it possible to execute more instructions at the same time. This often yields more performance and concurrent execution.

What kinds of problems motivates the need for concurrent execution? (Or phrased differently: What problems do concurrency help in solving?)
Everything you want to do at "the same time". For example two leds at the same time at different intervals, or playing a game at the same time as browsing the internet. Or making the mouse move while using the computer.

Does creating concurrent programs make the programmer's life easier? Harder? Maybe both? (Come back to this after you have worked on part 4 of this exercise)
Much easier. The logic becoms esier, no need for complex state machines, and we need fewer timers.

What are the differences between processes, threads, green threads, and coroutines?
Process: OS-managed. Truly concurrent if hardware support.
Thread: OS-managed. In the same adress space as its parent and all its other threads.
Green-Thread: Threads scheduled by a runtime libary or virtual machine. Managed in user space and not in kernel space.
Corutines: Routines that can be stopped and started again later on.

Which one of these do pthread_create() (C/POSIX), threading.Thread() (Python), go (Go) create?
They create threads.

How does pythons Global Interpreter Lock (GIL) influence the way a python Thread behaves?
It is a lock that prevents multiple threads from exeuting Python bytecodes at once. The lock is necessary because Pythons memory management is not threadsafe.

With this in mind: What is the workaround for the GIL (Hint: it's another module)?
cPython

What does func GOMAXPROCS(n int) int change?
It limits the number of operating system threads that can execute user-level Go code.






