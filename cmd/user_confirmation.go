package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func checkConfirm(detail string, cflag *bool) bool {
	return *cflag || askConfirmf(detail)
}

func checkConfirmF(detail string, cflag *bool, vars ...any) bool {
	return *cflag || askConfirmf(detail, vars...)
}

// TODO: (low) Switch to just using cobra.Command.InOrStdin()

// askConfirmF Does NOT check any flags
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
