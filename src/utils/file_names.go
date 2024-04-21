package utils

import "fmt"

func GetSubGraphOutputFileName(partitionId int) string {
	return "./output_graphs/SubGraph-" + fmt.Sprint(partitionId) + ".json"
}

func GetSuperStepSubGraphOutputFileName(partitionId int, superStep int) string {
	return "./output_graphs/SuperStep-" + fmt.Sprint(superStep) + "-SubGraph-" + fmt.Sprint(partitionId) + ".json"
}
