# Localization and Native Menu Lifecycle

Single-language non-English MyGo apps (for example, Simplified Chinese apps written without an `i18n` catalog) commonly hit mixed-language menus and system dialogs unless all three layers are configured.

## 1. Three-layer checklist

1. **Link MyGo's built-in locale package** (otherwise MyGo role menus, dialogs, and controls stay in English):
   ```go
   import _ "github.com/egoist/mygo/locales/zh-Hans" // or zh-Hant, ja
   ```
2. **Declare the app's written language before `App.Run()`**:
   ```go
   mygo.App.SetName("MyApp")
   mygo.App.SetLanguages("zh-Hans")
   ```
   Without an `i18n` catalog or `SetLanguages`, `mygo.App.Locale()` defaults to `"en"`, causing `plugins/updater`'s `updater.MenuItem()` and update window to render in English even on a Chinese OS.
3. **Declare macOS bundle localizations in `mygo.json`**:
   ```json
   "macos": {
     "infoPlist": {
       "CFBundleDevelopmentRegion": "zh-Hans",
       "CFBundleLocalizations": ["zh-Hans"]
     }
   }
   ```
   AppKit reads these keys to localize system-injected menu items (Services, Dictation, Emoji & Symbols) and native Open/Save/Alert dialog buttons.

## 2. Install `App.SetMenu(...)` inside `WhenReady`

Do **not** call `mygo.App.SetMenu(mygo.NewMenu(...))` before `App.Run()` when using standard menu roles such as `mygo.RoleEditMenu`:

- Before `App.Run()` calls `startLocale()`, MyGo's internal `shown` locale is still `nil`, so role submenus expand in English.
- During startup, `relabelMenus()` rewrites items that have a non-empty `Role`, but role-less submenus (such as the `Edit -> Speech` submenu generated inside `RoleEditMenu`) keep their initial English label (`"Speech"`).
- Calling `mygo.App.SetMenu(...)` at the top of `mygo.App.WhenReady(...)` runs right after `startLocale()`, so every role item and generated submenu receives the resolved locale (`"语音"`, `"检查更新…"`, etc.) on first build.
