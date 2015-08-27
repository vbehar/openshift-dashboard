package main

/*

This file is a hack to fix the cross-arch non-awareness of godep.
See https://github.com/tools/godep/issues/174 for more details.

TL;DR
To run "godep save ./..." on darwin, uncomment the following block.
It will force godep to see the linux-specific dependencies.

*/

/*

import (
	"github.com/docker/libcontainer/cgroups/fs"
)

var (
	_ = &fs.Manager{}
)

*/
