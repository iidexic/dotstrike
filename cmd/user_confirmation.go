package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func checkConfirm(detail string, cflag *bool) bool {
	if *cflag {
		return true
	}
	yes, e := askConfirmf(detail)
	if e != nil {
		fmt.Fprintf(os.Stderr, "confirm read failed: %s\n", e)
		return false
	}
	return yes
}

func checkConfirmF(detail string, cflag *bool, vars ...any) bool {
	if *cflag {
		return true
	}
	yes, e := askConfirmf(detail, vars...)
	if e != nil {
		fmt.Fprintf(os.Stderr, "confirm read failed: %s\n", e)
		return false
	}
	return yes
}

// promptYN asks a yes/no question; on stdin error prints to stderr and returns false.
func promptYN(detail string, vars ...any) bool {
	yes, e := askConfirmf(detail, vars...)
	if e != nil {
		fmt.Fprintf(os.Stderr, "confirm read failed: %s\n", e)
		return false
	}
	return yes
}

// TODO: (low) Switch to just using cobra.Command.InOrStdin()

// TODO:(med) Add a full prompt for user input

// askConfirmf prompts until yes/no or 4 blanks. Returns (false, err) on stdin failure.
func askConfirmf(detail string, vars ...any) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	n := 4
	fdetail := "[request to confirm]: " + detail + "\n:::"
	var noindex, yesindex int = -1, -1
	for {
		fmt.Printf(fdetail, vars...)
		response, err := reader.ReadString('\n')
		if err != nil {
			return false, err
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
			return yesindex < noindex, nil
		} else if yesindex >= 0 {
			return true, nil
		} else if noindex >= 0 {
			return false, nil
		}
		n--
		if n < 0 {
			return false, nil
		}

	}
}
