package view

import (
	h "mos/helper"
	"mos/model"
	"fmt"
	"html/template"
	"time"
	"strings"
	"mos/db"
)

type Header struct {
	Title    string
	SubTitle string
	SubtitleIdentifier string
	Switches []map[string]interface{};
	/**
	structture is fe.:
	{
		"id" : "hu-nohu"
		"label" : "Magyar vagy kulfoldi tulajdon"
	}
	*/
	Medias []map[string]string
	/**
	structure is fe.:
	{
		"label" : "Television",
		"class" : "color-news"
	}
	*/
	MediaTitle string;
	Options      []map[string]string;
	/**
	structure is fe.:
	{
		"option" : "Television",
		"label" : "color-news"
	}
	*/
	OptionTitle string;
	IncomeOptions      []map[string]string;
	/**
	structure is fe.:
	{
		"option" : "Television",
		"label" : "color-news"
	}
	*/
	Introduction string
	MoreIdentifier     string
	MoreLabel    string
	Years []int;
	StartYear int;
	EndYear int;
	FooterBlocks []model.Block
}

func (head *Header) Init(session *h.Session) {
	head.Title = h.Lang.Trans("Media Ownership Project", session.GetActiveLang());
	head.SubTitle = h.Lang.Trans("rólunk", session.GetActiveLang())
	head.SubtitleIdentifier = "about";
	head.MediaTitle = h.Lang.Trans("Média", session.GetActiveLang())
	head.OptionTitle = h.Lang.Trans("Hierarchia", session.GetActiveLang())

	head.setIntroduction(session);
	head.setSwitches(session);
	head.setMediatypes(session);
	head.setIncomeOptions(session);
	head.setOptions(session);
	head.setFooterBocks(session);
	head.setYears();
}

func (head *Header) setYears(){
	var mediaOwner model.MediaOwner;
	var row = db.DbMap.QueryRow(fmt.Sprintf("SELECT MIN(`year`) as miny, MAX(year) as maxy FROM %s",mediaOwner.GetTable()));
	var err = row.Scan(&head.StartYear, &head.EndYear);
	h.Error(err,"",h.ERROR_LVL_ERROR);

	var i int = head.StartYear;
	for(i <= time.Now().Year()){
		head.Years = append(head.Years,i);
		i++;
	}
}

func (head *Header) setSwitches(session *h.Session) {
	var span = "<span class=\"switch--red\">%v </span>%v<span class=\"switch--blue\"> %v </span> %v</label>";
	var hun = h.Lang.Trans("Magyar",session.GetActiveLang());
	var or = h.Lang.Trans("vagy",session.GetActiveLang());
	var foreign = h.Lang.Trans("külföldi",session.GetActiveLang());
	var grov = h.Lang.Trans("Kormánypárti",session.GetActiveLang());
	var opposit = h.Lang.Trans("Ellenzéki",session.GetActiveLang());
	var owner = h.Lang.Trans("tulajdon",session.GetActiveLang());
	var news = h.Lang.Trans("Hírmédia",session.GetActiveLang());
	var notnews = h.Lang.Trans("nem hírmédia",session.GetActiveLang());

	var hunswhtml = fmt.Sprintf(span,hun,or,foreign,owner);
	var leftrightswhtml = fmt.Sprintf(span,grov,or,opposit,"");
	var newsnotnewshtml = fmt.Sprintf(span,news,or,notnews,"");

	head.Switches = append(head.Switches, map[string]interface{}{
		"id":      "hu-nonhu",
		"label":   template.HTML(hunswhtml),
		"checked": false,
	});
	head.Switches = append(head.Switches, map[string]interface{}{
		"id":      "left-right",
		"label":   template.HTML(leftrightswhtml),
		"checked": false,
	});
	head.Switches = append(head.Switches, map[string]interface{}{
		"id":      "media-nonmedia",
		"label":   template.HTML(newsnotnewshtml),
		"checked": false,
	});
}

func (head *Header) setMediatypes(session *h.Session) {
	var mt model.MediaType;
	var mediaTypes []model.MediaType;
	mediaTypes = mt.GetAll();
	for _, mte := range mediaTypes {
		head.Medias = append(head.Medias, map[string]string{
			"label": h.Lang.Trans(mte.Name, session.GetActiveLang()),
			"class": mte.GetColorClass(),
		});
	}
}

func (head *Header) setIntroduction(session *h.Session) {
	var intro model.Block;
	intro, err := intro.GetByIdentifier("introduction", session.GetActiveLang());
	if (err != nil) {
		return;
	}
	if (intro.Id > 0) {
		head.Introduction = intro.Content;
		head.MoreLabel = h.Lang.Trans("Többet a módszertanról", session.GetActiveLang());
		head.MoreLabel = "methodology";
	}
}

func (head *Header) setOptions(session *h.Session) {
	head.Options = append(head.Options, map[string]string{
		"option": "owner",
		"label":  h.Lang.Trans("tulajdonos", session.GetActiveLang()),
	});

	head.Options = append(head.Options, map[string]string{
		"option": "income",
		"label":  h.Lang.Trans("árbevétel", session.GetActiveLang()),
	});

	head.Options = append(head.Options, map[string]string{
		"option": "interest",
		"label":  h.Lang.Trans("érdekeltség", session.GetActiveLang()),
	});
}

func (head *Header) setIncomeOptions(session *h.Session) {
	head.IncomeOptions = append(head.IncomeOptions, map[string]string{
		"option": "net",
		"label":  h.Lang.Trans("Nettó árbevétel", session.GetActiveLang()),
	});

	head.IncomeOptions = append(head.IncomeOptions, map[string]string{
		"option": "tax",
		"label":  h.Lang.Trans("Adózott eredmény", session.GetActiveLang()),
	});

	head.IncomeOptions = append(head.IncomeOptions, map[string]string{
		"option": "operational",
		"label":  h.Lang.Trans("Üzemi eredmény", session.GetActiveLang()),
	});
}

func (head *Header) setFooterBocks(session *h.Session){
	var blockIdentifiers []string = []string{"impressum","privacy"};
	var results []model.Block;
	var block model.Block;

	var query = `SELECT * FROM %v`;
	var where []string;
	var orWhere []string;

	where = append(where,fmt.Sprintf("`lc` = \"%v\"",session.GetActiveLang()));

	for _,bi := range blockIdentifiers{
		orWhere = append(orWhere,fmt.Sprintf(`identifier = "%v"`,bi));
	}

	where = append(where,"(" + strings.Join(orWhere, " OR ") + ")");

	query += " WHERE " + strings.Join(where, " AND ");

	db.DbMap.Select(&results,fmt.Sprintf(query,block.GetTable()));
	head.FooterBlocks = results;
}


