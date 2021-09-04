package DataStructureUtil

/*
 @author king 409060350@qq.com
*/

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
