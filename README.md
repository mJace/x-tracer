# x-tracer

Server Streaming the filtered PID information from testpod's as follows in Real Time:

<pre>
Choose pod : 2
---------------------------------------------
The pod you chose is testpod
Container ID is ...
19fb910a711f5eabf2cad6569a01db3702752e2f8155059654130711b8bf2c8f
2020/04/01 12:45:25 Hostname :  dad
Start Agent Pod
Start Agent Service

 {Probe:tcpaccept |T: 348.477 | PID:5644 | PNAME:iperf3 | IP:6 | RADDR:::ffff:127.0.0.1 | RPORT:45952 | LADDR:::ffff:127.0.0.1 | LPORT:6001
{Probe:tcpaccept |T: 348.521 | PID:5644 | PNAME:iperf3 | IP:6 | RADDR:::ffff:127.0.0.1 | RPORT:45954 | LADDR:::ffff:127.0.0.1 | LPORT:6001
{Probe:tcpaccept |T: 351.667 | PID:5644 | PNAME:iperf3 | IP:6 | RADDR:::ffff:127.0.0.1 | RPORT:45978 | LADDR:::ffff:127.0.0.1 | LPORT:6001
{Probe:tcptracer |T: 331.377 | PID:32703 | PNAME:iperf3 |IP->4 | SADDR:127.0.0.1 | DADDR:127.0.0.1 | SPORT:45848 | DPORT:6001
{Probe:tcptracer |T: 334.427 | PID:5644 | PNAME:iperf3 |IP->6 | SADDR:[::] | DADDR:[0:ffff:7f00:1::] | SPORT:0 | DPORT:65535
{Probe:tcptracer |T: 334.512 | PID:32703 | PNAME:iperf3 |IP->4 | SADDR:127.0.0.1 | DADDR:127.0.0.1 | SPORT:45848 | DPORT:6001
{Probe:tcptracer |T: 334.513 | PID:5644 | PNAME:iperf3 |IP->6 | SADDR:[::] | DADDR:[0:ffff:7f00:1::] | SPORT:0 | DPORT:65535
{Probe:tcptracer |T: 334.516 | PID:32703 | PNAME:iperf3 |IP->4 | SADDR:127.0.0.1 | DADDR:127.0.0.1 | SPORT:45846 | DPORT:6001
{Probe:tcptracer |T: 334.520 | PID:32717 | PNAME:iperf3 |IP->4 | SADDR:127.0.0.1 | DADDR:127.0.0.1 | SPORT:45864 | DPORT:6001
</pre>
