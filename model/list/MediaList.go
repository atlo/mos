package list

import (
	"github.com/valyala/fasthttp"
	"mos/db"
	"fmt"
	h "mos/helper"
	m "mos/model"
	"strconv"
	"strings"
)

type MediaList struct {
	List m.List
}

func (ml *MediaList) Init(ctx *fasthttp.RequestCtx, lang string) {
	var media m.Media;
	ml.List.Init(ctx, media, lang)

	var mediaType m.MediaType;
	mediaTypeOptions := mediaType.GetOptions(map[string]string{"value":"","label":"Anything"});

	ml.List.AddSearchParam(m.SearchParam{"ID","id","number","id", nil});
	ml.List.AddSearchParam(m.SearchParam{"Name","name","text","name",nil});
	ml.List.AddSearchParam(m.SearchParam{"Media Type", "media_type_id", "select", "media_type_id", map[string]interface{}{"options": mediaTypeOptions}});
	ml.List.AddSearchParam(m.SearchParam{"News","news","bool","news",nil});
}

func (ml *MediaList) SetLimitParam(limitParam string) {
	ml.List.SetLimitParam(limitParam)
}

func (ml *MediaList) SetPageParam(pageParam string) {
	ml.List.SetPageParam(pageParam);
}

func (ml *MediaList) Render(elements []m.Media) string {
	var headers []map[string]string;
	var rows []map[string]string;
	var options map[string]string;
	headers = []map[string]string{
		{"col": "id","title":"ID","order":"true"},
		{"col":"name","title":"Name"},
		{"col":"media_type_id","title":"Media Type"},
		{"col":"news","title":"News"},
		{"col":"actions","title":"Actions"},
	}

	var mediaType m.MediaType;
	for _, med := range elements {
		mediaType,err := mediaType.Get(med.MediaTypeId);
		h.Error(err,"",h.ERROR_LVL_WARNING);

		var actions []string;
		actions = append(actions, h.Replace(
			`<a href="%link%">%title%</a>`,
			[]string{"%link%","%title%"},
			[]string{h.GetUrl("media/edit",  []string{strconv.Itoa(int(med.Id))},true,"admin"), "[Edit]"},
		))

		/*actions = append(actions, h.Replace(
			`<a href="%link%" onclick="return window.confirm('Are you sure you want to delete the item?')">%title%</a>`,
			[]string{"%link%","%title%"},
			[]string{h.GetUrl("owner/delete",  []string{strconv.Itoa(int(o.Id))},true,"admin"), "[Delete]"},
		))*/

		rows = append(rows, map[string]string{
			"id":         strconv.Itoa(int(med.Id)),
			"name":      med.Name,
			"media_type_id":      mediaType.Name,
			"news": strconv.FormatBool(med.News),
			"actions":strings.Join(actions,"&nbsp;&nbsp;"),
		});
	}

	options = map[string]string{
		"class": "table-striped table-bordered table-hover",
		"id":    "media-list-table",
	};
	return ml.List.Render(headers, rows, options);
}

func (ml *MediaList) GetAll() []m.Media {
	var results []m.Media
	var where string = ml.List.GetSqlParams();
	if(where != ""){
		where = fmt.Sprintf(" WHERE %v",where);
	}
	sql := fmt.Sprintf("SELECT * FROM %v%v ORDER BY %v %v", ml.List.Table, where, ml.List.GetOrder(), ml.List.GetOrderDir());
	h.PrintlnIf(sql,h.GetConfig().Mode.Debug);
	_, err := db.DbMap.Select(&results, sql);
	h.Error(err, "", h.ERROR_LVL_ERROR)
	return results;
}

func (ml *MediaList) GetToPage() []m.Media {
	var results []m.Media
	var where string = ml.List.GetSqlParams();
	if(where != ""){
		where = fmt.Sprintf(" WHERE %v",where);
	}
	sql := fmt.Sprintf("SELECT * FROM %v%v ORDER BY %v %v LIMIT %v", ml.List.Table, where, ml.List.GetOrder(), ml.List.GetOrderDir(), ml.List.GetLimitString());
	h.PrintlnIf(sql,h.GetConfig().Mode.Debug);
	_, err := db.DbMap.Select(&results, sql);
	h.Error(err, "", h.ERROR_LVL_ERROR)
	return results;
}
