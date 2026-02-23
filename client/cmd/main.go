package main

import (
	"fmt"
	"os"
	"tml-sync/client/configs"
	"tml-sync/client/internal/api"
	"tml-sync/client/internal/i18n"
	"tml-sync/client/internal/scanner"
	"tml-sync/client/internal/ui"
	"tml-sync/shared/pkg"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg := configs.Load()

	// We'll use a pointer to the program to send messages from our sync goroutine
	var p *tea.Program

	onConfirm := func(scan bool) tea.Cmd {
		return func() tea.Msg {
			go func() {
				apiClient := api.New(cfg.Host, cfg.Port)

				p.Send(ui.StatusMsg(i18n.T("checking_status")))
				
				// 1. Get Language
				lang, err := apiClient.GetLanguage()
				if err != nil {
					p.Send(ui.LogMsg(i18n.T("error_msg", err.Error())))
					return
				}
				i18n.SetLanguage(lang)

				// 2. Version Check
				verResp, err := apiClient.GetVersion()
				if err != nil {
					p.Send(ui.LogMsg(i18n.T("error_msg", err.Error())))
					return
				}
				
				if verResp.ServerVersion != pkg.Version {
					p.Send(ui.LogMsg(i18n.T("version_mismatch", pkg.Version, verResp.ServerVersion)))
					
					p.Send(ui.LogMsg(i18n.T("triggering_server_update", pkg.Version)))
					err := apiClient.TriggerServerUpdate(pkg.Version)
					if err != nil {
						p.Send(ui.LogMsg(i18n.T("update_trigger_error", err.Error())))
					} else {
						p.Send(ui.LogMsg(i18n.T("server_updating")))
					}
				}

				// 3. Sync Status
				p.Send(ui.StatusMsg(i18n.T("sync_check")))
				syncStatus, err := apiClient.GetSyncStatus()
				if err != nil {
					p.Send(ui.LogMsg(i18n.T("error_msg", err.Error())))
					return
				}

				serverHashes := make(map[string]string)
				for _, m := range syncStatus.Mods {
					serverHashes[m.Name] = m.Hash
				}

				if !scan {
					apiClient.Stop()
					p.Send(ui.StateMsg(ui.StateDone))
					p.Send(ui.StatusMsg(i18n.T("done")))
					return
				}

				// 4. Scan and Sync
				workshopPath := scanner.GetSteamWorkshopPath()
				if workshopPath == "" {
					p.Send(ui.LogMsg(i18n.T("steam_not_found")))
					p.Send(ui.StateMsg(ui.StateDone))
					return
				}

				p.Send(ui.StatusMsg(i18n.T("scanning_spinner")))
				foundMods, err := scanner.ScanMods(workshopPath, verResp.TMLVersion)
				if err != nil {
					p.Send(ui.LogMsg(i18n.T("scan_error", err.Error())))
					return
				}

				p.Send(ui.StateMsg(ui.StateSyncing))

				var modsToUpload []scanner.FoundMod
				var totalBytes int64
				for _, mod := range foundMods {
					if h, exists := serverHashes[mod.Metadata.Name]; !exists || h != mod.Metadata.Hash {
						info, err := os.Stat(mod.Path)
						if err == nil {
							totalBytes += info.Size()
							modsToUpload = append(modsToUpload, mod)
						}
					} else {
						p.Send(ui.LogMsg(i18n.T("mod_exists", mod.Metadata.Name)))
					}
				}

				// Check enabled.json
				enabledPath := scanner.GetEnabledJSONPath()
				var enabledBytes int64
				var uploadEnabled bool
				var enabledHash string
				if enabledPath != "" {
					localHash, err := pkg.CalculateSHA256(enabledPath)
					if err == nil && localHash != syncStatus.EnabledJSONHash {
						enabledHash = localHash
						info, err := os.Stat(enabledPath)
						if err == nil {
							enabledBytes = info.Size()
							totalBytes += enabledBytes
							uploadEnabled = true
						}
					}
				}

				if len(modsToUpload) == 0 && !uploadEnabled {
					p.Send(ui.LogMsg(i18n.T("all_up_to_date")))
				}

				var cumulativeBytes int64
				for _, mod := range modsToUpload {
					p.Send(ui.StatusMsg(i18n.T("mod_uploading", mod.Metadata.Name)))

					err := apiClient.UploadMod(mod.Path, mod.Metadata.Name, mod.Metadata.Version, mod.Metadata.Hash, func(total, sent int64) {
						if totalBytes > 0 {
							p.Send(ui.ProgressMsg(float64(cumulativeBytes+sent) / float64(totalBytes)))
						}
					})

					if err != nil {
						p.Send(ui.LogMsg(i18n.T("mod_error", mod.Metadata.Name, err)))
					} else {
						p.Send(ui.LogMsg(i18n.T("mod_success", mod.Metadata.Name)))
					}
					
					// Update cumulative bytes even on error to keep progress moving
					info, err := os.Stat(mod.Path)
					if err == nil {
						cumulativeBytes += info.Size()
					}
				}

				// Upload enabled.json if needed
				if uploadEnabled {
					p.Send(ui.StatusMsg(i18n.T("uploading_config")))
					err := apiClient.UploadEnabledJSON(enabledPath, enabledHash, func(total, sent int64) {
						if totalBytes > 0 {
							p.Send(ui.ProgressMsg(float64(cumulativeBytes+sent) / float64(totalBytes)))
						}
					})
					if err != nil {
						p.Send(ui.LogMsg(i18n.T("error_uploading_config", err)))
					} else {
						p.Send(ui.LogMsg(i18n.T("config_synced")))
					}
				}

				apiClient.Stop()
				p.Send(ui.StateMsg(ui.StateDone))
				p.Send(ui.StatusMsg(i18n.T("done")))
			}()
			return nil
		}
	}

	m := ui.NewModel(cfg.Host, cfg.Port, onConfirm)
	p = tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
