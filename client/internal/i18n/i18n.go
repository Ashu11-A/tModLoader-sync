package i18n

import (
	"fmt"
	"os"
	"strings"
)

var (
	currentLanguage = "en"
	messages        = map[string]map[string]string{
		"en": {
			"welcome":          "TMODLOADER-SYNC CLIENT",
			"connecting":       "Connecting to server at %s:%d...",
			"lang_detected":    "Server language: %s",
			"version_check":    "Checking versions...",
			"version_match":    "Versions match (Client: %s, Server: %s)",
			"version_mismatch": "Version mismatch! Client: %s, Server: %s.",
			"tml_version":       "Server tModLoader version: %s",
			"sync_check":       "Checking synchronized mods...",
			"scan_prompt":      "Allow scanning system for Steam Workshop mods?",
			"scanning":         "Scanning for mods in: %s",
			"mod_found":        "Found mod: %s (Version: %s)",
			"mod_exists":       "Mod %s is already up to date.",
			"mod_uploading":    "Uploading: %s",
			"mod_success":      "Synced: %s",
			"mod_error":        "Error syncing %s: %v",
			"all_up_to_date":   "All mods and config are up to date.",
			"uploading_config": "Uploading enabled.json...",
			"error_uploading_config": "Error uploading enabled.json: %v",
			"config_synced":    "enabled.json synced successfully.",
			"done":             "Synchronization complete!",
			"scanning_spinner": "Searching for files...",
			"checking_status":  "Checking server status...",
			"confirm_yes":      "YES",
			"confirm_no":       "NO",
			"press_to_start":   "Press Enter to start synchronization",
			"quit":             "Press Q to quit",
			"error_msg":        "Error: %v",
			"steam_not_found":  "Steam Workshop not found.",
			"scan_error":       "Scan error: %v",
			"triggering_server_update": "Triggering server update to v%s...",
			"update_trigger_error": "Failed to trigger server update: %s",
			"server_updating": "Server is updating and will restart soon.",
		},
		"pt": {
			"welcome":          "TMODLOADER-SYNC CLIENT",
			"connecting":       "Conectando ao servidor em %s:%d...",
			"lang_detected":    "Idioma do servidor: %s",
			"version_check":    "Verificando versões...",
			"version_match":    "As versões coincidem (Cliente: %s, Servidor: %s)",
			"version_mismatch": "Versões diferentes! Cliente: %s, Servidor: %s.",
			"tml_version":       "Versão do tModLoader no servidor: %s",
			"sync_check":       "Verificando mods sincronizados...",
			"scan_prompt":      "Permitir escanear o sistema por mods da Steam?",
			"scanning":         "Procurando mods em: %s",
			"mod_found":        "Mod encontrado: %s (Versão: %s)",
			"mod_exists":       "O mod %s já está atualizado.",
			"mod_uploading":    "Enviando: %s",
			"mod_success":      "Sincronizado: %s",
			"mod_error":        "Erro ao sincronizar %s: %v",
			"all_up_to_date":   "Todos os mods e configurações estão atualizados.",
			"uploading_config": "Enviando enabled.json...",
			"error_uploading_config": "Erro ao enviar enabled.json: %v",
			"config_synced":    "enabled.json sincronizado com sucesso.",
			"done":             "Sincronização concluída!",
			"scanning_spinner": "Procurando arquivos...",
			"checking_status":  "Verificando status do servidor...",
			"confirm_yes":      "SIM",
			"confirm_no":       "NÃO",
			"press_to_start":   "Pressione Enter para iniciar a sincronização",
			"quit":             "Pressione Q para sair",
			"error_msg":        "Erro: %v",
			"steam_not_found":  "Oficina Steam não encontrada.",
			"scan_error":       "Erro ao escanear: %v",
			"triggering_server_update": "Iniciando atualização do servidor para v%s...",
			"update_trigger_error": "Falha ao iniciar atualização do servidor: %s",
			"server_updating": "O servidor está se atualizando e reiniciará em breve.",
		},
	}
)

func init() {
	lang := os.Getenv("LANG")
	if strings.HasPrefix(lang, "pt") {
		currentLanguage = "pt"
	}
}

func SetLanguage(lang string) {
	if lang == "pt" || lang == "br" || lang == "pt-br" {
		currentLanguage = "pt"
	} else {
		currentLanguage = "en"
	}
}

func GetLanguage() string {
	return currentLanguage
}

func T(key string, args ...interface{}) string {
	msg, ok := messages[currentLanguage][key]
	if !ok {
		// Fallback to English if not found in current language
		msg, ok = messages["en"][key]
		if !ok {
			return key
		}
	}
	return fmt.Sprintf(msg, args...)
}
