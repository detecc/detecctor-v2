# 🌐 Multiple languages support - i18n

The library used for localizing the messages is [go-i18n](https://github.com/nicksnyder/go-i18n).

## Setting the language

You can change the language in which the bot replies. The language is set per chat/user basis and can be changed by
issuing the `/lang` or `/language` command. The command accepts one parameter - the language tag, which must suffice the
**[BCP 47](https://en.wikipedia.org/wiki/IETF_language_tag)** standard.

Example command(s):

```text 
/lang sv   #switches to swedish 
/lang en   #switches to english 
/language sl   #switches to slovenian 
/language al   #switches to albanian 
```

If the message set does not have the desired language translations in the [translations](../pkg/i18n/translations)
folder, the message will be in English. When first registering a chat/user, the default language set is English.

## 👉 Contributing to translations

Every translation contribution is welcome! Issue a PR with the translation(s). You only need to add a translation yaml
file - everything else is handled. To add new translation messages, be sure to check out how to do so in
the [plugin guide](service/plugin/plugins.md#translating-replies) and
the [go-i18n library](https://github.com/nicksnyder/go-i18n).
