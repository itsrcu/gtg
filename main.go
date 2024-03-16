package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/rs/zerolog"
)

type (
	inventoryAPI struct {
		LastAssetID string `json:"last_assetid"`

		Assets []struct {
			AssetID    string `json:"assetid"`
			ClassID    string `json:"classid"`
			InstanceID string `json:"instanceid"`
		} `json:"assets"`

		Descriptions []struct {
			Name       string `json:"name"`
			TypeLong   string `json:"type"`
			ClassID    string `json:"classid"`
			InstanceID string `json:"instanceid"`
			Tags       []struct {
				LocaleTagName string `json:"localized_tag_name"`
			} `json:"tags"`
			AppID int `json:"market_fee_app"`
		} `json:"descriptions"`

		LoadMore int `json:"more_items"`
	}

	gooAPI struct {
		Value   string `json:"goo_value"`
		Message string `json:"message"`
		Success int    `json:"success"`
	}

	grindAPI struct {
		GooTotal string `json:"goo_value_total"`
		Success  int    `json:"success"`
	}

	jsonConfig struct {
		SteamID             *string `json:"steamID"`
		VanityLink          *string `json:"steamVanityLink"`
		SessionID           *string `json:"sessionID"`
		AccessToken         *string `json:"accessToken"`
		KeepCount           *int    `json:"keepCount"`
		LoadEntireInventory *bool   `json:"loadEntireInventory"`
		Blacklist           struct {
			KeepItemType []string `json:"keepItemType"`
			KeepAppID    []int    `json:"keepAppID"`
			KeepGame     struct {
				KeepMethod    string   `json:"keepMethod"`
				KeepNames     []string `json:"keepNames"`
				KeepThreshold float64  `json:"keepThreshold"`
			} `json:"keepGame"`
			KeepItem struct {
				KeepMethod        string   `json:"keepMethod"`
				KeepNames         []string `json:"keepNames"`
				KeepThreshold     float64  `json:"keepThreshold"`
				IncludeTypeSearch bool     `json:"includeTypeSearch"`
			} `json:"keepItem"`
		} `json:"blackList"`
	}
)

func setupLogger() zerolog.Logger {
	logFile, logErr := os.OpenFile("log.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o600)
	if logErr != nil {
		log.Panicf("Failed to create log.txt file: %s", logErr.Error())
	}

	cWriter := zerolog.ConsoleWriter{Out: os.Stderr}
	multiWriter := zerolog.MultiLevelWriter(logFile, cWriter)

	return zerolog.New(multiWriter).With().Timestamp().Logger()
}

func loadConfig() jsonConfig {
	var cfg jsonConfig

	configFile, configErr := os.OpenFile("config.json", os.O_RDONLY, 0o600)
	if configErr != nil {
		log.Panicf("Failed to open config file: %s", configErr.Error())
	}

	if configDecodeErr := json.NewDecoder(configFile).Decode(&cfg); configDecodeErr != nil {
		log.Panicf("Failed to decode config file: %s", configDecodeErr.Error())
	}

	if cfg.SteamID == nil || *cfg.SteamID == "" ||
		cfg.VanityLink == nil || *cfg.VanityLink == "" ||
		cfg.SessionID == nil || *cfg.SessionID == "" ||
		cfg.AccessToken == nil || *cfg.AccessToken == "" ||
		cfg.KeepCount == nil || cfg.LoadEntireInventory == nil {
		log.Panic("missing or empty config fields")
	}

	lowerGameMethod := strings.ToLower(cfg.Blacklist.KeepGame.KeepMethod)
	lowerItemMethod := strings.ToLower(cfg.Blacklist.KeepItem.KeepMethod)

	if lowerGameMethod != "both" && lowerGameMethod != "levenshtein" && lowerGameMethod != "contains" ||
		lowerItemMethod != "both" && lowerItemMethod != "levenshtein" && lowerItemMethod != "contains" {
		log.Println("one of your keep methods seems invalid, is it correct? if so, press enter.")
		fmt.Scanln()
	}

	if len(cfg.Blacklist.KeepAppID) == 0 ||
		len(cfg.Blacklist.KeepItemType) == 0 ||
		len(cfg.Blacklist.KeepGame.KeepNames) == 0 ||
		len(cfg.Blacklist.KeepItem.KeepNames) == 0 {
		log.Println("there are some empty fields in your blacklist, is it correct? if so, press enter.")
		fmt.Scanln()
	}

	return cfg
}

func matchContains[T comparable](item T, array []T, toString func(T) string) bool {
	for i := 0; i < len(array); i++ {
		if strings.Contains(toString(array[i]), toString(item)) {
			return true
		}
	}

	return false
}

func matchLevenshtein[T comparable](item T, array []T, toString func(T) string, metric *metrics.Levenshtein, treshold float64) bool {
	for i := 0; i < len(array); i++ {
		if similarity := strutil.Similarity(toString(array[i]), toString(item), metric) * 100; similarity >= treshold {
			return true
		}
	}

	return false
}

func main() {
	log.Println(`
         888
         888
         888
 .d88b.  888888 .d88b.
d88P"88b 888   d88P"88b
888  888 888   888  888
Y88b 888 Y88b. Y88b 888
 "Y88888  "Y888 "Y88888
     888            888
Y8b d88P       Y8b d88P
 "Y88P"         "Y88P"`)

	config := loadConfig()
	gtgLogger := setupLogger()

	timeoutClient := &http.Client{
		Timeout: 1 * time.Minute,
	}

	levenshteinMetric := &metrics.Levenshtein{
		CaseSensitive: false,
		InsertCost:    1,
		DeleteCost:    1,
		ReplaceCost:   1,
	}

	inventoryRequest, inventoryRequestErr := http.NewRequestWithContext(context.Background(), http.MethodGet, fmt.Sprintf("https://steamcommunity.com/inventory/%s/753/6?l=english&count=5000", *config.SteamID), http.NoBody)
	if inventoryRequestErr != nil {
		gtgLogger.Fatal().Err(inventoryRequestErr).Msg("failed to create inventory request")
	}

	inventoryResponse, inventoryErr := timeoutClient.Do(inventoryRequest)
	if inventoryErr != nil {
		gtgLogger.Fatal().Err(inventoryErr).Msg("failed to get inventory")
	}
	defer inventoryResponse.Body.Close()

	var inventory inventoryAPI

	if inventoryDecodeErr := json.NewDecoder(inventoryResponse.Body).Decode(&inventory); inventoryDecodeErr != nil {
		gtgLogger.Fatal().Err(inventoryDecodeErr).Msg("failed to decode inventory data")
	}

	for inventory.LoadMore == 1 && *config.LoadEntireInventory {
		gtgLogger.Info().Str("lastassetid", inventory.LastAssetID).Msg("sending request to load more inventory items")

		ldMoreRequest, ldMoreRequestErr := http.NewRequestWithContext(context.Background(), http.MethodGet, fmt.Sprintf("https://steamcommunity.com/inventory/%s/753/6?l=english&count=5000&start_assetid=%s", *config.SteamID, inventory.LastAssetID), http.NoBody)
		if ldMoreRequestErr != nil {
			gtgLogger.Fatal().Err(ldMoreRequestErr).Msg("failed to get inventory")
		}

		ldMoreResponse, ldMoreResponseErr := timeoutClient.Do(ldMoreRequest)
		if ldMoreResponseErr != nil {
			gtgLogger.Fatal().Err(ldMoreResponseErr).Msg("failed to get inventory")
		}
		defer ldMoreResponse.Body.Close()

		var ldInventory inventoryAPI

		if ldInventoryDecodeErr := json.NewDecoder(ldMoreResponse.Body).Decode(&ldInventory); ldInventoryDecodeErr != nil {
			gtgLogger.Fatal().Err(ldInventoryDecodeErr).Msg("failed to decode inventory data")
		}

		inventory.Assets = append(inventory.Assets, ldInventory.Assets...)
		inventory.Descriptions = append(inventory.Descriptions, ldInventory.Descriptions...)
		inventory.LoadMore = ldInventory.LoadMore
		inventory.LastAssetID = ldInventory.LastAssetID
	}

	gtgLogger.Info().Int("itemcount", len(inventory.Descriptions)).Msg("loaded inventory")

	for i := 0; i < len(inventory.Descriptions); i++ {
		item := inventory.Descriptions[i]
		itemTags := item.Tags

		// safety check
		if !(len(itemTags) > 2) {
			gtgLogger.Error().Int("appid", item.AppID).Int("tagcount", len(itemTags)).Msg("missing minimum tag count of 2, skipping")
			continue
		}

		gameName := itemTags[1].LocaleTagName
		itemType := itemTags[2].LocaleTagName

		// appID check
		if matchContains(item.AppID, config.Blacklist.KeepAppID, strconv.Itoa) {
			gtgLogger.Info().Int("appid", item.AppID).Msg("appid is on the blacklist, skipping")
			continue
		}

		// itemType check
		if matchContains(itemType, config.Blacklist.KeepItemType, strings.ToLower) {
			gtgLogger.Info().Str("game", gameName).Str("itemname", item.Name).Str("itemtype", item.TypeLong).Msg("item type is on the blacklist, skipping")
			continue
		}

		keepFunc := make(map[string]func(name string, array []string, treshold float64) bool)
		keepFunc["contains"] = func(name string, array []string, _ float64) bool {
			return matchContains(name, array, strings.ToLower)
		}

		keepFunc["levenshtein"] = func(name string, array []string, treshold float64) bool {
			return matchLevenshtein(name, array, strings.ToLower, levenshteinMetric, treshold)
		}

		keepFunc["both"] = func(name string, array []string, treshold float64) bool {
			if matchContains(name, array, strings.ToLower) {
				return true
			}

			return matchLevenshtein(name, array, strings.ToLower, levenshteinMetric, treshold)
		}

		// gameName check
		if keepFunc[strings.ToLower(config.Blacklist.KeepGame.KeepMethod)](gameName, config.Blacklist.KeepGame.KeepNames, config.Blacklist.KeepGame.KeepThreshold) {
			gtgLogger.Info().Str("game", gameName).Str("itemname", item.Name).Str("itemtype", item.TypeLong).Msg("gamename is on the blacklist, skipping")
			continue
		}

		// itemName check
		if keepFunc[strings.ToLower(config.Blacklist.KeepItem.KeepMethod)](item.Name, config.Blacklist.KeepItem.KeepNames, config.Blacklist.KeepItem.KeepThreshold) {
			gtgLogger.Info().Str("game", gameName).Str("itemname", item.Name).Str("itemtype", item.TypeLong).Msg("itemname is on the blacklist, skipping")
			continue
		}

		// itemTypeLong check
		if config.Blacklist.KeepItem.IncludeTypeSearch {
			if keepFunc[strings.ToLower(config.Blacklist.KeepItem.KeepMethod)](item.TypeLong, config.Blacklist.KeepItem.KeepNames, config.Blacklist.KeepItem.KeepThreshold) {
				gtgLogger.Info().Str("game", gameName).Str("itemname", item.Name).Str("itemtype", item.TypeLong).Msg("item type (long) is on the blacklist, skipping")
				continue
			}
		}

		var assetIDs []string

		for k := 0; k < len(inventory.Assets); k++ {
			if item.ClassID == inventory.Assets[k].ClassID &&
				item.InstanceID == inventory.Assets[k].InstanceID {
				assetIDs = append(assetIDs, inventory.Assets[k].AssetID)
			}
		}

		for v := 0; v < len(assetIDs)-*config.KeepCount; v++ {
			gooValueRequest, gooValueRequestErr := http.NewRequestWithContext(context.Background(), http.MethodGet, fmt.Sprintf("https://steamcommunity.com/id/%s/ajaxgetgoovalue/?sessionid=%s&appid=%d&assetid=%s&contextid=6", *config.VanityLink, *config.SessionID, item.AppID, assetIDs[v]), http.NoBody)
			if gooValueRequestErr != nil {
				gtgLogger.Error().Err(gooValueRequestErr).Str("game", gameName).Str("itemname", item.Name).Str("itemtype", item.TypeLong).Msg("failed to create gem value request")
				continue
			}

			gooValueRequest.Header.Set("Cookie", "sessionid="+*config.SessionID+";steamLoginSecure="+*config.SteamID+"%7C%7C"+*config.AccessToken)

			gooValueResponse, gooValueErr := timeoutClient.Do(gooValueRequest)
			if gooValueErr != nil {
				gtgLogger.Error().Err(gooValueErr).Str("game", gameName).Str("itemname", item.Name).Str("itemtype", item.TypeLong).Msg("failed to get gem value for item")
				continue
			}
			defer gooValueResponse.Body.Close()

			var goo gooAPI

			if gooDecodeErr := json.NewDecoder(gooValueResponse.Body).Decode(&goo); gooDecodeErr != nil {
				gtgLogger.Error().Err(gooDecodeErr).Str("game", gameName).Str("itemname", item.Name).Str("itemtype", item.TypeLong).Msg("failed to decode gem value data")
				continue
			}

			if goo.Success != 0 {
				gtgLogger.Error().
					Str("status", gooValueResponse.Status).
					Str("url", gooValueRequest.URL.String()).
					Str("item", item.Name).
					Str("message", goo.Message).
					Int("success", goo.Success).
					Msg("non-zero success code")

				continue
			}

			body := bytes.NewBufferString(fmt.Sprintf(`sessionid=%s&appid=%d&assetid=%s&contextid=6&goo_value_expected=%s`, *config.SessionID, item.AppID, assetIDs[v], goo.Value))

			grindRequest, grindRequestErr := http.NewRequestWithContext(context.Background(), http.MethodPost, fmt.Sprintf("https://steamcommunity.com/id/%s/ajaxgrindintogoo", *config.VanityLink), body)
			if grindRequestErr != nil {
				gtgLogger.Error().Err(grindRequestErr).Str("game", gameName).Str("itemname", item.Name).Str("itemtype", item.TypeLong).Msg("failed to create item grind request")
				continue
			}

			grindRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
			grindRequest.Header.Set("Cookie", "sessionid="+*config.SessionID+";steamLoginSecure="+*config.SteamID+"%7C%7C"+*config.AccessToken)

			grindResponse, grindResponseErr := timeoutClient.Do(grindRequest)
			if grindResponseErr != nil {
				gtgLogger.Error().Err(grindResponseErr).Str("game", gameName).Str("itemname", item.Name).Str("itemtype", item.TypeLong).Msg("failed to grind item into gems")
				continue
			}
			defer grindResponse.Body.Close()

			var grind grindAPI

			if grindDecodeErr := json.NewDecoder(grindResponse.Body).Decode(&grind); grindDecodeErr != nil {
				gtgLogger.Error().Err(grindDecodeErr).Str("game", gameName).Str("itemname", item.Name).Str("itemtype", item.TypeLong).Msg("failed to decode grind data")
				continue
			}

			if grind.Success != 1 {
				gtgLogger.Error().
					Str("status", grindResponse.Status).
					Str("url", grindRequest.URL.String()).
					Str("itemname", item.Name).
					Int("success", grind.Success).
					Msg("non-one success code")

				continue
			}

			gtgLogger.Info().
				Str("itemname", item.Name).
				Str("itemtype", item.TypeLong).
				Str("gem value", goo.Value).
				Str("total gems", grind.GooTotal).
				Msg("item grinded into gems")
		}
	}
	gtgLogger.Info().Msg("finished")
}
