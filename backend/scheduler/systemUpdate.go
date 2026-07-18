package scheduler

import (
	"log"
	"os/exec"
)

func SystemUpdate() {
	// 执行系统更新逻辑
	cmd := exec.Command("bash", conf.SystemUpdateScript)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("system update failed: %v, output: %s", err, string(output))
	}
}
