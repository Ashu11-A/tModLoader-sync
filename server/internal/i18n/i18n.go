package i18n

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	LangFile = "logs/language.conf"
)

var (
	messages = map[string]map[string]string{
		"en": {
			"server_title":     "tModLoader-sync Server",
			"description":      "Automatically synchronizes your Steam mods with the server, facilitating the upload and maintenance. This tool was created to be used with the tModLoader egg found here: https://github.com/Ashu11-A/Ashu_eggs",
			"status":           "Status",
			"online":           "Online",
			"address":          "Address",
			"public_ip":        "Public IP",
			"sync_commands":    "Sync Commands for Clients:",
			"windows_ps":       "Windows (PowerShell):",
			"linux_bash":       "Linux (Bash):",
		},
		"pt": {
			"server_title":     "Servidor tModLoader-sync",
			"description":      "Sincroniza seus mods steam automaticamente com o servidor, facilitando o upload e a manutenção, essa ferramenta foi criada para ser usada com o egg Tmodloader encontrado aqui: https://github.com/Ashu11-A/Ashu_eggs",
			"status":           "Status",
			"online":           "Online",
			"address":          "Endereço",
			"public_ip":        "IP Público",
			"sync_commands":    "Comandos de Sincronização para Clientes:",
			"windows_ps":       "Windows (PowerShell):",
			"linux_bash":       "Linux (Bash):",
		},
	}
)

func Setup() {
	if _, err := os.Stat(LangFile); os.IsNotExist(err) {
		_ = os.MkdirAll("logs", 0755)
		fmt.Println("\n  \033[1;35mtModLoader-sync Server Setup\033[0m")
		fmt.Println("  ----------------------------------------")
		fmt.Println("  Select your language / Selecione seu idioma:")
		fmt.Println("  1. English (en)")
		fmt.Println("  2. Portuguese (pt)")
		fmt.Print("\n  Choice (1-2): ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		lang := "en"
		if input == "2" || strings.ToLower(input) == "pt" {
			lang = "pt"
		}

		_ = os.WriteFile(LangFile, []byte(lang), 0644)
		fmt.Printf("  ✅ Language set to: %s\n", lang)
		fmt.Println("  ----------------------------------------\n")
	}
}

func GetLanguage() string {
	data, err := os.ReadFile(LangFile)
	if err != nil {
		return "en"
	}
	return string(data)
}

func T(key string, args ...interface{}) string {
	lang := GetLanguage()
	msg, ok := messages[lang][key]
	if !ok {
		// Fallback to English if not found in current language
		msg, ok = messages["en"][key]
		if !ok {
			return key
		}
	}
	return fmt.Sprintf(msg, args...)
}
