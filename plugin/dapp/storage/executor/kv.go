package executor

/*
 *       kv   ，key
 *  key = keyPrefix + userKey
 *          ，  ’-‘
 */

var (
	//KeyPrefixStateDB state db key
	KeyPrefixStateDB = "mavl-storage-"
	//KeyPrefixLocalDB local db key
	KeyPrefixLocalDB = "LODB-storage-"
)

// Key Storage to save key
func Key(txHash string) (key []byte) {
	key = append(key, []byte(KeyPrefixStateDB)...)
	key = append(key, []byte(txHash)...)
	return key
}

func getLocalDBKey(txHash string) (key []byte) {
	key = append(key, []byte(KeyPrefixLocalDB)...)
	key = append(key, []byte(txHash)...)
	return key
}
