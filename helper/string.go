package helper

import (
	"math/rand"
	"time"
	"strings"
	"regexp"
)

const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789";

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomString(length int) string {
	return StringWithCharset(length, charset)
}

func Replace(replaceIn string,replaceKeys []string, replaceVals []string) string{
	for i,replaceKey := range replaceKeys{
		replaceVal := replaceVals[i];
		replaceIn = strings.Replace(replaceIn,replaceKey,replaceVal,-1);
	}

	return replaceIn;
}

func HtmlAttribute(key string, value string) string{
	var attrTemp string = `%key%="%val%"`
	attrVal := strings.Replace(attrTemp,"%key%", key, -1);
	attrVal = strings.Replace(attrVal,"%val%", value, -1);

	return attrVal;
}

func RemoveNewLines(subject string, removeLot bool) string {
	var replace []map[string]string = []map[string]string{
		{"exp":"(\r\n)","to":"<br />"},
		{"exp":"(\r)","to":"<br />"},
		{"exp":"(\n)","to":"<br />"},
	}

	if(removeLot){
		replace = append(replace,map[string]string{"exp":"(\\s*<br />){3,}","to":"<br /><br />"});
	}

	for _,rm := range replace{
		re := regexp.MustCompile(rm["exp"]);
		subject = re.ReplaceAllString(subject,rm["to"]);
	}

	return subject;
}

func TrimPath(path string) string{
	path = strings.Trim(path,"/");
	path = strings.Trim(path,"./");

	return path;
}

func Contains(slice []string, entry string) bool {
	for _,se := range slice{
		if se == entry{
			return true;
		}
	}
	return false;
}

func RemoveAccents(stringToClean string) string{
	replacer := strings.NewReplacer("á","a","é","e","í","i","ó","o","ö","o","ő","o","ú","u","ü","u","ű","u"," ","-")
	cleanString := replacer.Replace(stringToClean)
	return cleanString;
}
