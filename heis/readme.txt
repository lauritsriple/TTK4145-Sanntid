Module overview:
main
--> driver
--> --> driver.go ------- Driver for elevator
--> --> io.go ----------- Low level io driver
--> --> io.c ------------ Low level io driver in c
--> --> ch.go ----------- Defines for the elevator
--> --> channels.h ------ Defines for the elvator in c
--> --> io.h ------------ Defines for io.c
--> udp
--> --> udp.go ---------- Multicast network module
--> control
--> --> control.go ------ Control logic
--> --> messageparser.go- Message logic
--> localQueue
--> --> localQueue.go --- Queue for orders on running elevator (not global)
--> fsmelev
--> --> fsmelev.go ------ Poll and set driver. Execute orders. (almost a statemachine)

Scripts for easy deploy:
deploy.sh (takes multiple last 3.digits in ip as input, ssh to machines and run localDeploy.sh)
localDeploy.sh (git and build)
runElev.sh (takes multiple last 3-digits in ip as input,spawns multiple terminals with ssh to machines and run localRunElev.sh)
localRunElev.sh(runs elevator)
NOTMADEYET killElevs(takes multiple last 3-digits in ip as input, ssh and cleanup)

TODO:
Messageparser should change name to messagehandler
Seperate messageparser and control to to packages?
Implement stop and obstruction
Cost is suboptimal. Could be better if it also uses the length of queue
No orders is lost, but messages might get lost. Should implement some kind of ack for received messages
FIXED Handle motorstop (Motor cable plugged out)
FIXED Easy deployment scripts to test on elevators
Script to kill program on all running elevators that does not close other students editors editing a main file
Create channels inside init functions and parse them downwards instead of using them locally inside package (FIXED in some packages)
Change name of fsmelev to something smarter, fsm maybe?
