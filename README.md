# gtg
Grind To Gems, a tool to help grind unnecessary items from your steam inventory.. into gems!

```json
{
    "steamID": "",                      // https://steamcommunity.com -> F12 -> console.log(g_steamID)
    "steamVanityLink": "",              // https://steamcommunity.com/my -> Redirect -> https://steamcommunity.com/id/{vanity link}
    "sessionID": "",                    // https://steamcommunity.com -> F12 -> console.log(g_sessionID)
    "accessToken": "",                  // https://steamcommunity.com -> F12 -> console.log(application_config.dataset.loyalty_webapi_token) Thanks! https://steamapi.xpaw.me/
    "keepCount": 1,                     // decides how many copies of an item should be kept (if there are multiple of one)
    "loadEntireInventory": false,       // decides if the entire inventory should be loaded, if false, only 5000 items will be loaded
    "blackList": {
        "keepItemType": [""],           // ignores item(s) from the set type
        "keepAppID": [],                // ignores item(s) from the set appid
        "keepGame": {                   // decides which game(s) should be ignored with the set conditions
            "keepMethod": "Both",       // sets method used for the lookup
            "keepNames": [""],          // sets game name(s) to lookup
            "keepThreshold": 70.0       // sets threshold used for the lookup
        },
        "keepItem": {                   // decides which item(s) should be ignored with the set conditions
            "keepMethod": "Both",       // sets method used for the lookup
            "keepNames": [""],          // sets item name(s) to lookup
            "keepThreshold": 70.0,      // sets threshold used for the lookup
            "includeTypeSearch": true   // sets if long item types should be used for the lookup
        }
    }
}
```

```json
"keepItemType": [],                     // for example, "Emoticon", "Normal", "Profile Background", "{Season} Sale {Year}". Normal = cards. Some emoticons do not have this type, and is recommended to put "Emoticon" in the "keepItem" -> "keepNames"
"keepAppID": []                         // integer array for appIDs. https://store.steampowered.com/app/{This is the appID}/{App name}
"keepNames": [],                        // string array with game names. Will be used by "Contains and/or "Levenshtein"
"keepMethod": "",                       // sets mehod used for the lookup, can be "Contains", "Levenshtein", or "Both" (this is lowered, and not case-sensitive). If "Both" is specified, then "Contains" will run first, then Levenshtein. "Contains" will lower both the input, and the item/game name
"includeTypeSearch": true,              // default is true, will look at the "Type" of the item. Best used with "Contains", as the "Type" also contains the name of the game
"keepThreshold": 70.0,                  // sets threshold used for the Levenshtein lookup, 0-100, keeps anything equal or over the threshold
"keepCount": 1                          // if there are multiple of the same item, "keepCount" determines how many of them should be kept
```

<code>This README is a WIP that contains the very basics, and will be updated to be more concise later.</code>