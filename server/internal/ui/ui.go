package ui

import (
	"fmt"
	"strings"
	"tml-sync/server/internal/i18n"
	"tml-sync/server/internal/network"

	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(false)

	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EEEEEE"))

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1)

	commandBoxStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Background(lipgloss.Color("#111111")).
			Padding(0, 1)
)

func PrintBanner(addr string) {
	var s strings.Builder

	// Header
	s.WriteString(headerStyle.Render(i18n.T("server_title")) + "\n\n")

	// Description
	s.WriteString(descStyle.Render(i18n.T("description")) + "\n\n")

	var content strings.Builder

	// Status line
	content.WriteString(fmt.Sprintf("● %s: %s\n", labelStyle.Render(i18n.T("status")), statusStyle.Render(i18n.T("online"))))
	
	// Address line
	content.WriteString(fmt.Sprintf("➜ %s: %s\n", labelStyle.Render(i18n.T("address")), lipgloss.NewStyle().Underline(true).Render("http://0.0.0.0"+addr)))

	ip, err := network.GetPublicIP()
	if err == nil {
		content.WriteString(fmt.Sprintf("➜ %s: %s\n", labelStyle.Render(i18n.T("public_ip")), lipgloss.NewStyle().Underline(true).Render(ip)))
		content.WriteString("\n" + infoStyle.Render(i18n.T("sync_commands")) + "\n\n")
		
		// Windows Command
		content.WriteString(labelStyle.Render(i18n.T("windows_ps")) + "\n")
		windowsCmd := fmt.Sprintf("powershell -c \"`$env:TML_HOST='%s';`$env:TML_PORT='%s';irm https://raw.githubusercontent.com/Ashu11-A/tModLoader-sync/main/sync.ps1|iex\"", ip, addr)
		content.WriteString(commandBoxStyle.Render(windowsCmd) + "\n\n")

		// Linux Command
		content.WriteString(labelStyle.Render(i18n.T("linux_bash")) + "\n")
		linuxCmd := fmt.Sprintf("curl -fsSL https://raw.githubusercontent.com/Ashu11-A/tModLoader-sync/main/sync.sh | h=%s p=%s bash", ip, addr)
		content.WriteString(commandBoxStyle.Render(linuxCmd) + "\n")
	}

	s.WriteString(boxStyle.Render(content.String()))

	fmt.Println(s.String())
}
