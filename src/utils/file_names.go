package utils

import "fmt"

func GetSubGraphOutputFileName(partitionId int) string {
	return "./output_graphs/SubGraph-" + fmt.Sprint(partitionId) + ".json"
}
