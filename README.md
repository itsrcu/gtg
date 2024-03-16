# gtg
Grind To Gems, a tool to help grind unnecessary items from your steam inventory.. into gems!

```json
{
    "steamID": "",
    "steamVanityLink": "",
    "sessionID": "",
    "accessToken": "",
    "keepCount": 1,
    "loadEntireInventory": false,
    "blackList": {
        "keepItemType": [""],
        "keepAppID": [753],
        "keepGame": {
            "keepMethod": "Both",
            "keepNames": [""],
            "keepThreshold": 70.0
        },
        "keepItem": {
            "keepMethod": "Both",
            "keepNames": [""],
            "keepThreshold": 70.0,
            "includeTypeSearch": true
        }
    }
}
```

### steamID
> The term SteamID refers to a unique identifier for a specific Steam account. It can be formatted in a few different ways, but the most commonly used is the account's 64 bit ID which is a 17 digit number. The instructions below will help you find your SteamID.
> Excerpt from https://help.steampowered.com/en/faqs/view/2816-BE67-5B69-0FEC

**Method A)**
* Open your Steam Community Profile (https://steamcommunity.com/my)
* Click Edit Profile
* If you've never set a custom Steam Community URL for your account, your 64 bit ID will will be shown in the URL under the CUSTOM URL box in the format 76561198#########
* If you have set a custom URL for your account, you can delete the text in the CUSTOM URL box to see your account's 64 bit ID in the URL listed below.

**Method B)**
* Open your Steam Community Profile (https://steamcommunity.com/my)
* Open the developer console (F12)
* Type the following: `console.log(g_steamID)`

**Method C)**
* Open SteamDB (https://steamdb.info)
* Login with your Steam account
* Open SteamDB's calculator (https://steamdb.info/calculator)
* Your SteamID will be visible above the lookup section: `You are currently signed in as 76561198#########. Click this or your avatar in the navbar to view your profile.`

### steamVanityLink
A custom URL you can set that will redirect to your Steam profile

Navigating to https://steamcommunity.com/my will redirect you to the value you set it as, for example, https://steamcommunity.com/id/###, where ### is your vanity/custom link

### sessionID
A 24 character long hex string, used in the `steamLoginSecure` cookie. This renews every 24 hours

**Method A)**
* Open your Steam Community Profile (https://steamcommunity.com/my)
* Open the developer console (F12)
* Type the following: `console.log(g_sessionID)`

**Method B)**
* Open your Steam Community Profile (https://steamcommunity.com/my)
* Open the developer console (F12)
* Type the following: `console.log(document.cookie.split("; ").find((s)=>s.startsWith("sessionid="))?.split("=")[1])`
> https://developer.mozilla.org/en-US/docs/Web/API/Document/cookie

### accessToken
A very long character string, used in the `steamLoginSecure` cookie. This renews every 24 hours

**Method A)**
* Open your Steam Community Profile (https://steamcommunity.com/my)
* Open the developer console (F12)
* Type the following: `console.log(application_config.dataset.loyalty_webapi_token)`
> https://steamapi.xpaw.me

### keepCount
If you have multiple copies of the same item, keepCount will determine how many of them should be kept. 1 will keep one copy of each item, 2 will keep two, and so on..

### loadEntireInventory
Initially, only 5000 items will be loaded into the internal `inventory` variable. Setting this to true will check if there are more items to load, and won't continue until the entire inventory is loaded into the variable

### keepItemType
Keeps specified item types from being shredded. **It's advised to also put this in the keepItem's keepNames**, as for example, sale items aren't marked as such, but  instead marked as, for example, "2023 Winter Sale"

### keepAppID
Keeps items from the specified appID from being shredded

**Method A)**
* Open the game's store page
* https://store.steampowered.com/app/###/$$$, where ### is the appID, and $$$ is the game's name

**Method B)**
* Open SteamDB (https://steamdb.info)
* Search for the game's name, and the `App ID` column should be right below the game's name

### keepMethod
Determines which method should be used for game/item name lookups. Valid options are `Both`, `Contains` or `Levenshtein` (these are lowered internally, and are not case-sensitive in the config file). If `Both` is specified, then `Contains` will run first, then `Levenshtein`

### keepNames
Keeps specified games/items from being shredded. This is not case-sensitive, as `Contains` will lower both the config's input, and the game/item's names, and `Levenshtein`'s config is set to not be case-sensitive.

### keepThreshold
Used in the Levenshtein lookup, can be `0.0`-`100.0`, keeps anything equal or over the threshold

### includeTypeSearch
**This is different from keepItemType**. This uses the inventory's item description type for the lookups, for example, `Detroit: Become Human Profile Background` or `Resident Evil 4 (2005) Emoticon`