package utils

var FILECOIN_GENESIS_UNIX_EPOCH int64 = 1598306400

func UnixToFilEpoch(unixEpoch int64) int64 {
	return unixEpoch - FILECOIN_GENESIS_UNIX_EPOCH/30
}
