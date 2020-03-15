package page

import (
	"mos/helper"
	"time"
	"mos/db"
	"mos/model"
	"fmt"
	"html/template"
	"strings"
)

type Graph struct {
	Title string
	Subtitle []string
	Items []map[string]interface{}
	YearFrom int
	YearTo int
	Years []int
	Lc string
}

func (g *Graph) Init(session *helper.Session){
	g.Lc = session.GetActiveLang();

	g.Title = helper.Lang.Trans("Cégnév",g.Lc);

	g.Subtitle = append(g.Subtitle, helper.Lang.Trans("Alapitasi év",g.Lc));
	g.Subtitle = append(g.Subtitle, helper.Lang.Trans("Székhely",g.Lc));
	g.Subtitle = append(g.Subtitle, helper.Lang.Trans("Anyaország",g.Lc));

	g.Items = append(g.Items, map[string]interface{}{
		"value":"interest",
		"label":helper.Lang.Trans("Érdekeltség",g.Lc),
	});
	g.Items = append(g.Items, map[string]interface{}{
		"value":"finalowner",
		"label":helper.Lang.Trans("Végső tulajdonos",g.Lc),
	});
	g.Items = append(g.Items, map[string]interface{}{
		"value":"owner",
		"label":helper.Lang.Trans("Tulajdonos",g.Lc),
	});
	g.Items = append(g.Items, map[string]interface{}{
		"value":"media",
		"label":helper.Lang.Trans("Média",g.Lc),
	});

	var mediaOwner model.MediaOwner;
	var row = db.DbMap.QueryRow(fmt.Sprintf("SELECT MIN(`year`) as miny, MAX(`year`) as maxy FROM %v",mediaOwner.GetTable()));
	var err = row.Scan(&g.YearFrom,&g.YearTo);
	helper.Error(err,"",helper.ERROR_LVL_ERROR);

	var i int = g.YearFrom;
	for(i <= g.YearTo){
		g.Years = append(g.Years,i);
		i++;
	}
}

func (g Graph) GetContent() template.HTML{
	var temp string = "page/graphoverlay.html";
	hasContent,content := helper.CacheStorage.GetString(temp,[]string{"index","graphcontent",g.Lc});
	if(!hasContent){
		fmt.Println(g);
		content = helper.GetScopeTemplateString(temp,g,"frontend");
		helper.CacheStorage.Set(temp,[]string{"index","graphcontent",g.Lc},time.Hour*12,content);
	}

	return template.HTML(content);
}

func (g Graph) GetSubtitle() string{
	return strings.Join(g.Subtitle," | ");
}
