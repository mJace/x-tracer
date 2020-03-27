# x-tracer

Server Streaming the bcc tool(tcpconnect) output as follows in Real Time:

<pre>
2020/03/27 17:12:02 Start x-tracer

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
2020/03/27 17:12:06 Hostname :  dad
Start Agent Pod
Start Agent Service

 0.000    15281  coredns      4  127.0.0.1        127.0.0.1        8080
PID: 15281
ProbeName: tcpconnect

 0.132    13503  kubelet      4  169.254.25.10    169.254.25.10    9254
PID: 13503
ProbeName: tcpconnect

 0.229    13503  kubelet      4  192.168.123.38   192.168.123.38   8081
PID: 13503
ProbeName: tcpconnect

 0.369    13503  kubelet      4  192.168.123.38   192.168.123.38   8081
PID: 13503
ProbeName: tcpconnect

 0.604    14450  node-cache   4  169.254.25.10    169.254.25.10    9254
PID: 14450
ProbeName: tcpconnect

 1.000    15281  coredns      4  127.0.0.1        127.0.0.1        8080
</pre>

