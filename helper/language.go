package helper

import (
	"os"
	"strings"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
)

var Lang *Language;

var DefLang string = "hu";

var LangQueryKey string = "lang";

func InitLanguage() {
	PrintlnIf("Initializing translator", GetConfig().Mode.Debug);
	Lang = &Language{};
	Lang.Init();
	PrintlnIf("Translator initialization done", GetConfig().Mode.Debug);
}

type Language struct {
	storage map[string]map[string]string;
	available []string;
}

func (l *Language) GetStorage() map[string]map[string]string {
	return l.storage;
}

func (l *Language) GetAvailableLanguageCodes() []string {
	return l.available;
}

func (l *Language) Init() {
	l.storage = make(map[string]map[string]string);
	var path string = "./resource/language"
	dir, err := os.Open(path);
	Error(err, "", ERROR_LVL_ERROR);
	files, err := dir.Readdir(0);
	Error(err, "", ERROR_LVL_ERROR)
	for _, f := range files {
		if (!f.IsDir()) {
			var parts []string = strings.Split(f.Name(), ".");
			if (parts[len(parts)-1] == "json") {
				dat, err := ioutil.ReadFile(path + "/" + f.Name());
				Error(err, "", ERROR_LVL_ERROR);
				if (err != nil) {
					continue;
				}
				var toData map[string]string;
				var mapKey string = strings.Replace(f.Name(), ".json", "", -1);
				err = json.Unmarshal(dat, &toData)
				l.storage[mapKey] = toData;
				Error(err, "", ERROR_LVL_ERROR);
			}
		}
	}

	l.setAvailableLanguages();
}

func (l *Language) setAvailableLanguages(){
	var appendDefLang bool = true;
	for c, _ := range l.storage {
		if(Contains(GetConfig().Language.Allowed,c)) {
			l.available = append(l.available, c);
		}
		if (c == DefLang) {
			appendDefLang = false;
		}
	}

	if (appendDefLang) {
		l.available = append(l.available, DefLang);
	}
}

func (l *Language) IsAvailable(lang string) bool {
	_, exists := l.storage[lang];
	return Contains(l.available,lang) && (exists || DefLang == lang);
}

func (l *Language) Trans(txtToTrans string, toLang string) string {
	langMap, ok := l.storage[toLang];
	if (!ok) {
		return txtToTrans;
	}

	translated, ok := langMap[txtToTrans];
	if (!ok) {
		return txtToTrans;
	}

	return translated;
}

func (l *Language) SetLanguage(ctx *fasthttp.RequestCtx, session *Session) {
	var langFromQuery string = string(ctx.FormValue(LangQueryKey));
	if (langFromQuery != "" && l.IsAvailable(langFromQuery)) {
		PrintlnIf(fmt.Sprintf("Changing language to %v", langFromQuery), GetConfig().Mode.Debug);
		session.SetActiveLang(langFromQuery);
	} else if (session.GetActiveLang() == "") {
		PrintlnIf("Setting default language to user", GetConfig().Mode.Debug);
		session.SetActiveLang(DefLang);
	}
}
