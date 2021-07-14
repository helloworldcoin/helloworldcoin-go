package DataStructureUtil

func IsExistDuplicateElement(datas *[]string) bool {
	visited := make(map[string]bool, 0)
	for i := 0; i < len(*datas); i++ {
		if visited[(*datas)[i]] == true {
			return true
		} else {
			visited[(*datas)[i]] = true
		}
	}
	return false
}
