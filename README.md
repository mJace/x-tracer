# x-tracer

Server Streaming the bcc tool(tcpconnect) output as follows in Real Time:

<pre>
2020/03/27 16:14:43 Start x-tracer

---------------------------------------------
0 : default
1 : kube-node-lease
2 : kube-public
3 : kube-system
Choose namespace : 0
---------------------------------------------
The namespace you chose is default

---------------------------------------------
0 : learnpod
1 : nginx
Choose pod : 1
---------------------------------------------
The pod you chose is nginx
Container ID is ...
f03e72f265ee911ca5e7b6ae4a70e4b149402fcd964273713d616d574f4e3133
2020/03/27 16:14:48 Hostname :  dad
Start Agent Pod
Start Agent Service

 0.000    14450  node-cache   4  169.254.25.10    169.254.25.10    9254
PID: 123
ProbeName:  probe

 0.395    15281  coredns      4  127.0.0.1        127.0.0.1        8080
PID: 123
ProbeName:  probe

 1.000    14450  node-cache   4  169.254.25.10    169.254.25.10    9254
PID: 123
ProbeName:  probe

 1.395    15281  coredns      4  127.0.0.1        127.0.0.1        8080
PID: 123
ProbeName:  probe

 1.705    13503  kubelet      4  169.254.25.10    169.254.25.10    9254
PID: 123
ProbeName:  probe

 2.000    14450  node-cache   4  169.254.25.10    169.254.25.10    9254
PID: 123
ProbeName:  probe

 2.395    15281  coredns      4  127.0.0.1        127.0.0.1        8080
PID: 123
ProbeName:  probe

 2.776    12010  calico-node  4  127.0.0.1        127.0.0.1        9099
PID: 123
ProbeName:  probe</pre>

