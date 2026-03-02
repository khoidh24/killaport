package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	defer waitEnter()
	showBanner()

	port, err := readPort()
	if err != nil {
		fmt.Println("Invalid port:", err)
		waitEnter()
		os.Exit(1)
	}

	pids, err := findPIDsByPort(port)
	if err != nil {
		fmt.Println("Find PID error:", err)
		waitEnter()
		os.Exit(1)
	}
	if len(pids) == 0 {
		fmt.Printf("No process is listening on port %d\n", port)
		waitEnter()
		return
	}

	fmt.Printf("Found PID(s) on port %d: %v\n", port, pids)

	for _, pid := range uniqueInts(pids) {
		if err := killPID(pid); err != nil {
			fmt.Printf("Kill PID %d failed: %v\n", pid, err)
		} else {
			fmt.Printf("Killed PID %d\n", pid)
		}
	}

	waitEnter()
}

func readPort() (int, error) {
	fmt.Print("Enter port to kill: ")
	reader := bufio.NewReader(os.Stdin)
	s, _ := reader.ReadString('\n')
	s = strings.TrimSpace(s)

	p, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if p < 1 || p > 65535 {
		return 0, errors.New("port must be 1..65535")
	}
	return p, nil
}

func findPIDsByPort(port int) ([]int, error) {
	switch runtime.GOOS {
	case "windows":
		// netstat -ano | findstr :<port>
		cmd := exec.Command("cmd", "/C", fmt.Sprintf("netstat -ano | findstr :%d", port))
		out, err := cmd.CombinedOutput()
		if err != nil {
			// netstat returns non-zero if findstr finds nothing
			// so treat as "no result" if output empty
			if len(strings.TrimSpace(string(out))) == 0 {
				return []int{}, nil
			}
		}
		return parseWindowsNetstat(string(out), port), nil

	case "linux":
		// Prefer lsof if available; fallback to ss
		if commandExists("lsof") {
			cmd := exec.Command("bash", "-lc", fmt.Sprintf("lsof -i :%d -sTCP:LISTEN -t 2>/dev/null", port))
			out, _ := cmd.CombinedOutput()
			return parseLinesOfInts(string(out)), nil
		}
		cmd := exec.Command("bash", "-lc", fmt.Sprintf("ss -ltnp 'sport = :%d' 2>/dev/null", port))
		out, _ := cmd.CombinedOutput()
		return parseLinuxSS(string(out)), nil

	default:
		return nil, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func killPID(pid int) error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/C", fmt.Sprintf("taskkill /PID %d /F", pid))
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%v: %s", err, strings.TrimSpace(string(out)))
		}
		return nil
	case "linux":
		cmd := exec.Command("bash", "-lc", fmt.Sprintf("kill -9 %d", pid))
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%v: %s", err, strings.TrimSpace(string(out)))
		}
		return nil
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func parseLinesOfInts(s string) []int {
	var res []int
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if n, err := strconv.Atoi(line); err == nil {
			res = append(res, n)
		}
	}
	return res
}

// Windows netstat sample line:
// TCP    0.0.0.0:3000   0.0.0.0:0   LISTENING   12345
func parseWindowsNetstat(out string, port int) []int {
	var res []int
	lines := strings.Split(out, "\n")
	rePID := regexp.MustCompile(`\s+(\d+)\s*$`)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Make sure line contains :port and LISTENING
		if !strings.Contains(line, fmt.Sprintf(":%d", port)) {
			continue
		}
		if !strings.Contains(strings.ToUpper(line), "LISTENING") {
			continue
		}
		m := rePID.FindStringSubmatch(line)
		if len(m) == 2 {
			if pid, err := strconv.Atoi(m[1]); err == nil {
				res = append(res, pid)
			}
		}
	}
	return res
}

// ss output contains users:(("proc",pid=123,fd=3))
func parseLinuxSS(out string) []int {
	var res []int
	re := regexp.MustCompile(`pid=(\d+)`)
	matches := re.FindAllStringSubmatch(out, -1)
	for _, m := range matches {
		if len(m) == 2 {
			if pid, err := strconv.Atoi(m[1]); err == nil {
				res = append(res, pid)
			}
		}
	}
	return res
}

func uniqueInts(nums []int) []int {
	seen := map[int]bool{}
	var res []int
	for _, n := range nums {
		if !seen[n] {
			seen[n] = true
			res = append(res, n)
		}
	}
	return res
}

func waitEnter() {
	fmt.Print("Press Enter to exit...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func showBanner() {
	banner := `
╔══════════════════════════════════════════════════════════════════════════════════════════════════╗
║                                                                                                  ║
║     ___  __    ___  ___       ___       ________  ________  ________  ________  _________        ║
║    |\  \|\  \ |\  \|\  \     |\  \     |\   __  \|\   __  \|\   __  \|\   __  \|\___   ___\      ║
║    \ \  \/  /|\ \  \ \  \    \ \  \    \ \  \|\  \ \  \|\  \ \  \|\  \ \  \|\  \|___ \  \_|      ║
║     \ \   ___  \ \  \ \  \    \ \  \    \ \   __  \ \   ____\ \  \\\  \ \   _  _\   \ \  \       ║
║      \ \  \\ \  \ \  \ \  \____\ \  \____\ \  \ \  \ \  \___|\ \  \\\  \ \  \\  \|   \ \  \      ║
║       \ \__\\ \__\ \__\ \_______\ \_______\ \__\ \__\ \__\    \ \_______\ \__\\ _\    \ \__\     ║
║        \|__| \|__|\|__|\|_______|\|_______|\|__|\|__|\|__|     \|_______|\|__|\|__|    \|__|     ║
║                                                                                                  ║
║                                                                                                  ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════╝
`
	fmt.Println(banner)
}
