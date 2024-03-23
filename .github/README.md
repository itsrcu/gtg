<h1 align="center">gtg 💎</h1>
<p align="center"><code>g</code>rind <code>t</code>o <code>g</code>ems, a tool to help automatically grind specified items from your Steam inventory... into gems!</p>

> [!CAUTION]
> **While gtg should be safe to use under normal circumstances (if it is configured properly), there is absolutely NO warranty for it. (see the <a href="https://github.com/itsrcu/gtg/blob/main/LICENSE">license</a>)**
> 
> **Use it at your own discretion.**

<h1 align="center">Installation and usage 💿</h1>

**Method A)**
1. Download the latest Go version from <a href="https://go.dev/dl">here</a> (https://go.dev/dl)
2. Follow the <a href="https://go.dev/doc/install">Installation instructions</a> (https://go.dev/doc/install)
3. Clone, or download the repository
    - To clone, type: `git clone https://github.com/itsrcu/gtg.git`
    - To download, click <a href="https://github.com/itsrcu/gtg/archive/refs/heads/main.zip">here</a> (https://github.com/itsrcu/gtg/archive/refs/heads/main.zip)
4. After cloning/unzipping, enter the `gtg(-main)` directory, and rename `config.example.json` to `config.json`
5. Configure gtg according to your preferences. See the <a href="https://github.com/itsrcu/gtg#configuration-">Configuration ⚙️</a> section
6. Run `go run main.go`
    - If configured properly, items based on the config are starting to be turned into gems

**Method B)**
1. Navigate to the <a href="https://github.com/itsrcu/gtg/releases/tag/latest">latest release</a> (https://github.com/itsrcu/gtg/releases/tag/latest)
2. Download your OS' zip file.
    - Every push will build gtg for several OS' and architectures
        - `legacy` builds use an older version of Go. **These exist for compatability, and should ONLY be used if your OS does not support Go 1.21+**. Those being:
            - `Windows 7 or newer`
            - `Kernel 2.6.32 or newer`
            - `macOS 10.13 High Sierra or newer`
3. Extract the zip, and run `gtg-<architecture>-<os>`
    - If you aren't sure about your architecture, run `uname -m`
        - `x86_64` = should run `amd64`
        - `i686` or `i386` = should run `i386`
        - `armv7l` = should run `arm`
        - `aarch64` = should run `arm64`

<h1 align="center">Configuration ⚙️</h1>

```json
{
    "steamID": "",
    "steamVanityLink": "",
    "sessionID": "",
    "accessToken": "",
    "keepCount": 1,
    "testDrive": true,
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
* If you've never set a custom Steam Community URL for your account, your 64 bit ID will be shown in the URL under the CUSTOM URL box in the format 76561198#########
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
A custom URL you can set that will redirect to your Steam profile.

Navigating to https://steamcommunity.com/my will redirect you to the value you set it as, for example, https://steamcommunity.com/id/###, where ### is your vanity/custom link.

### sessionID
A 24 character long hex string, used in the `steamLoginSecure` cookie. This renews every 24 hours.

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
A very long character string, used in the `steamLoginSecure` cookie. This renews every 24 hours.

**Method A)**
* Open your Steam Community Profile (https://steamcommunity.com/my)
* Open the developer console (F12)
* Type the following: `console.log(application_config.dataset.loyalty_webapi_token)`
> https://steamapi.xpaw.me

### keepCount
If you have multiple copies of the same item, keepCount will determine how many of them should be kept. 1 will keep one copy of each item, 2 will keep two, and so on..

### testDrive
Default set to true, in this mode, `sessionID` and `accessToken` aren't required, and items will not be shredded. Items will still be logged whether they got detected by one of the blacklist filters or not.

### loadEntireInventory
Initially, only 5000 items will be loaded into the internal `inventory` variable. Setting this to true will check if there are more items to load, and won't continue until the entire inventory is loaded into the variable.

### keepItemType
> [!IMPORTANT]
> **It's advised to also put this in the keepItem's keepNames**, as for example, sale items aren't marked as such, but  instead marked as, for example, "2023 Winter Sale".

Keeps specified item types from being shredded. Common types are `Emoticon`, `Normal` (for cards), `Profile Background`, `{Season} Sale {Year}`.

### keepAppID
Keeps items from the specified appID from being shredded.

**Method A)**
* Open the game's store page
* https://store.steampowered.com/app/###/$$$, where ### is the appID, and $$$ is the game's name

**Method B)**
* Open SteamDB (https://steamdb.info)
* Search for the game's name, and the `App ID` column should be right below the game's name

### keepMethod
Determines which method should be used for game/item name lookups. Valid options are `Both`, `Contains` or `Levenshtein` (these are lowered internally, and are not case-sensitive in the config file). If `Both` is specified, then `Contains` will run first, then `Levenshtein`.

### keepNames
Keeps specified games/items from being shredded. This is not case-sensitive, as `Contains` will lower both the config's input, and the game/items names, and `Levenshtein`'s config is set to not be case-sensitive.

### keepThreshold
Used in the Levenshtein lookup, can be `0.0` to `100.0`, keeps anything equal or over the threshold.

### includeTypeSearch
> [!NOTE]
> **This is different from keepItemType**.

This uses the inventory's item description type for the lookups, for example, `Detroit: Become Human Profile Background` or `Resident Evil 4 (2005) Emoticon`.
