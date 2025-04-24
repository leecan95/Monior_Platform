package config

const MemTotalReal = "1.3.6.1.4.1.2021.4.5.0"
const MemAvailReal = "1.3.6.1.4.1.2021.4.6.0"
const MemCacheReal = "1.3.6.1.4.1.2021.4.14.0"
const MemBufferReal = "1.3.6.1.4.1.2021.4.13.0"
const MemSwapTotal = "1.3.6.1.4.1.2021.4.3.0"
const MemSwapAvail = "1.3.6.1.4.1.2021.4.3.0"
const IOReceive = "1.3.6.1.4.1.2021.11.6.0"
const IOSent = "1.3.6.1.4.1.2021.11.5.0"

const CpuUsageUserPercent = "1.3.6.1.4.1.2021.11.9.0"
const SsCpuRawUser = "1.3.6.1.4.1.2021.11.50.0"
const SsCpuRawNice = "1.3.6.1.4.1.2021.11.51.0"
const SsCpuRawSystem = "1.3.6.1.4.1.2021.11.52.0"
const SsCpuRawIdle = "1.3.6.1.4.1.2021.11.53.0"
const SsCpuRawWait = "1.3.6.1.4.1.2021.11.54.0"
const SsCpuRawKernel = "1.3.6.1.4.1.2021.11.55.0"
const SsCpuRawInterrupt = "1.3.6.1.4.1.2021.11.56.0"
const SsCpuRawSoftIRQ = "1.3.6.1.4.1.2021.11.61.0"
const SsCpuUser = "1.3.6.1.4.1.2021.11.9.0"
const SsCpuSystem = "1.3.6.1.4.1.2021.11.10.0"

// The total amount of memory free or available for use on
// this host.  This value typically covers both real memory
// and swap space or virtual memory
const memTotalFree = "1.3.6.1.4.1.2021.11.0"
const memBuffer = "1.3.6.1.4.1.2021.4.13.0"
const memCached = "1.3.6.1.4.1.2021.4.14.0"
const memSysAvail = "1.3.6.1.4.1.2021.4.27.0"
const ssCpuSystem = "1.3.6.1.4.1.2021.11.10.0"
const DskTotal = "1.3.6.1.4.1.2021.9.1.6.3"
const DskAvail = "1.3.6.1.4.1.2021.9.1.7.3"

const DskTotalNew = "1.3.6.1.4.1.2021.9.1.6.1"

const DskAvailNew = "1.3.6.1.4.1.2021.9.1.7.1"

const Target = "localhost"

const Server246 = "172.21.5.246"
const Server245 = "172.21.5.245"
const Server244 = "172.21.5.244"
const Server250 = "172.21.5.250"
const Server251 = "172.21.5.251"
const Server252 = "172.21.5.252"

const Server94 = "10.207.191.94"
const Server95 = "10.207.191.95"
const Server96 = "10.207.191.96"
const Server97 = "10.207.191.97"
const Server98 = "10.207.191.98"
const Server99 = "10.207.191.99"
const Server100 = "10.207.189.100"
const Server101 = "10.207.189.101"
const Server102 = "10.207.189.102"
const Server104 = "10.207.189.104"
const Server105 = "10.207.189.105"
const Server106 = "10.207.189.106"
const Server108 = "10.207.189.108"
const Server109 = "10.207.189.109"
const Server110 = "10.207.189.110"
const Server112 = "10.207.189.112"
const Server113 = "10.207.189.113"
const Server114 = "10.207.189.114"
const Server116 = "10.207.189.116"
const Server117 = "10.207.189.117"
const Server118 = "10.207.189.118"
const Server119 = "10.207.189.119"
const Server120 = "10.207.189.120"
const Server121 = "10.207.189.121"
const Server123 = "10.207.189.123"
const Server124 = "10.207.189.124"
const Server125 = "10.207.189.125"
const Server127 = "10.207.189.127"
const Server128 = "10.207.189.128"
const Server129 = "10.207.189.129"
const Server130 = "10.207.189.130"
const Server132 = "10.207.189.132"
const Server133 = "10.207.189.133"
const Sv1 = "10.210.103.1"
const Sv2 = "10.210.103.2"
const Sv3 = "10.210.103.3"
const Sv4 = "10.210.103.4"
const Sv5 = "10.210.103.5"
const Sv6 = "10.210.103.6"
const Sv7 = "10.210.103.7"
const Sv88 = "10.207.191.88"
const Sv89 = "10.207.191.89"
const Sv90 = "10.207.191.90"
const Sv91 = "10.207.191.91"
const Sv92 = "10.207.191.92"
const Sv123 = "10.207.191.123"
const Sv124 = "10.207.191.124"
const Sv125 = "10.207.191.125"
const NumServers = 50

var Servers = [NumServers]string{
	Server94,
	Server95,
	Server96,
	Server97,
	Server98,
	Server99,
	Server100,
	Server101,
	Server102,
	Server104,
	Server105,
	Server106,
	Server108,
	Server109,
	Server110,
	Server112,
	Server113,
	Server114,
	Server116,
	Server117,
	Server118,
	Server119,
	Server120,
	Server121,
	Server123,
	Server124,
	Server125,
	Server127,
	Server128,
	Server129,
	Server130,
	Server132,
	Server133,
	Sv1,
	Sv2,
	Sv3,
	Sv4,
	Sv5,
	Sv6,
	Sv7,
	Sv88,
	Sv89,
	Sv90,
	Sv91,
	Sv92,
	Sv123,
	Sv124,
	Sv125,
}

const SNMP_Port = 161
