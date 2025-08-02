package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func askConfirmf(detail string, vars ...any) bool {
	reader := bufio.NewReader(os.Stdin)
	n := 4
	fdetail := "[request to confirm]: " + detail + "\n:::"
	var noindex, yesindex int = -1, -1
	for {
		fmt.Printf(fdetail, vars...)
		response, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		response = strings.ToLower(response)
		for _, no := range []string{"n", "no", "false"} {
			if hasno := strings.Contains(response, no); hasno {
				noindex = strings.Index(response, no)
				break
			}
		}
		for _, yes := range []string{"y", "yes", "true"} {
			if hasYes := strings.Contains(response, yes); hasYes {
				yesindex = strings.Index(response, yes)
				break
			}
		}
		if yesindex >= 0 && noindex >= 0 {
			if yesindex < noindex {
				return true
			} else {
				return false
			}

		} else if yesindex >= 0 {
			return true
		} else if noindex >= 0 {
			return false
		}
		n--
		if n < 0 {
			return false
		}

	}
}
