package reader

import (
	"bufio"
	"os"
	"simulation/common"
	"strings"
)

type JobDes struct {
	Name	   string
	Sub        string
	Running    string
	Allocation string
	Status     string
}

func (j *JobDes) Str() string {
	result := j.Sub + " " + j.Running + " " + j.Allocation
	return result
}

func ReadFile(path string) []JobDes {
	file, err := os.Open(path)
	common.Check(err)
	defer file.Close()
	result := make([]JobDes, 0)
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		index := strings.Index(line,";")
		if index > -1 {
			continue
		}
		tokens := strings.Fields(line)
		status := tokens[10]
		if status != "1" {
			continue
		}
		job := JobDes{
			Name:		strings.TrimSpace(tokens[0]),
			Sub:        strings.TrimSpace(tokens[1]),
			Running:    strings.TrimSpace(tokens[3]),
			Allocation: strings.TrimSpace(tokens[7]),
			Status:     status,
		}
		result = append(result, job)
	}
	return result
}
