package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// checkConfirm short-circuits to true if the assume-yes flag pointed to by cflag
// is set; otherwise it prompts. Stdin read failure is logged to stderr and
// treated as a deny (returns false) so an unusable terminal cannot silently
// approve a destructive op.
func checkConfirm(detail string, cflag *bool, vars ...any) bool {
	if *cflag {
		return true
	}
	return promptYN(detail, vars...)
}

// promptYN asks the user a yes/no question. Stdin read failure is logged to
// stderr and treated as a deny (returns false).
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

// askConfirmf prompts the user with detail (formatted via vars) until they
// answer yes/no or 4 blank responses pass (returning false). Returns
// (false, err) if reading stdin fails; callers should treat that as a deny.
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
