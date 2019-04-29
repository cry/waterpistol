package messages

import "math/rand"

/**
Creating actual structs for protobuf in go is really ugly,
 these are just functions that abstract the creation process away
*/

func Implant_wrap(module string, message isCheckCmdRequest_Message) *CheckCmdRequest {
	bytes := make([]byte, rand.Intn(100)+1) // Make a buffer of len 1-100
	rand.Read(bytes)                        // Fill with random data

	return &CheckCmdRequest{Module: module, Message: message, RandomPadding: bytes}
}

func Implant_heartbeat(heartbeat int64) *CheckCmdRequest {
	return Implant_wrap("", &CheckCmdRequest_Heartbeat{Heartbeat: heartbeat})
}

func Implant_error(module string, err int32) *CheckCmdRequest {
	return Implant_wrap(module, &CheckCmdRequest_Error{Error: err})
}

func Implant_data(module string, data []byte) *CheckCmdRequest {
	return Implant_wrap(module, &CheckCmdRequest_Data{Data: data})
}

func Implant_ipscan_in_progress(module string, ip string, port uint32) *CheckCmdRequest {
	ipscan := &IPScanReply{Status: IPScanReply_IN_PROGRESS, Ip: ip, Port: port}
	return Implant_wrap(module, &CheckCmdRequest_Ipscan{Ipscan: ipscan})
}

func Implant_ipscan_complete(module string) *CheckCmdRequest {
	ipscan := &IPScanReply{Status: IPScanReply_COMPLETE}
	return Implant_wrap(module, &CheckCmdRequest_Ipscan{Ipscan: ipscan})
}

func Implant_portscan_in_progress(module string, port uint32) *CheckCmdRequest {
	portscan := &PortScanReply{Status: PortScanReply_IN_PROGRESS, Found: port}
	return Implant_wrap(module, &CheckCmdRequest_Portscan{Portscan: portscan})
}

func Implant_portscan_complete(module string) *CheckCmdRequest {
	portscan := &PortScanReply{Status: PortScanReply_COMPLETE}
	return Implant_wrap(module, &CheckCmdRequest_Portscan{Portscan: portscan})
}
