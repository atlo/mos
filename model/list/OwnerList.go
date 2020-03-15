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

type OwnerList struct {
	List m.List
}

func (ol *OwnerList) Init(ctx *fasthttp.RequestCtx, lang string) {
	var owner m.Owner;
	ol.List.Init(ctx, owner, lang)

	ol.List.AddSearchParam(m.SearchParam{"ID","id","number","id", nil});
	ol.List.AddSearchParam(m.SearchParam{"Name","name","text","name",nil});
	ol.List.AddSearchParam(m.SearchParam{"Hungarian","hungarian","bool","hungarian",nil});
}

func (ol *OwnerList) SetLimitParam(limitParam string) {
	ol.List.SetLimitParam(limitParam)
}

func (ol *OwnerList) SetPageParam(pageParam string) {
	ol.List.SetPageParam(pageParam);
}

func (ol *OwnerList) Render(elements []m.Owner) string {
	var headers []map[string]string;
	var rows []map[string]string;
	var options map[string]string;
	headers = []map[string]string{
		{"col": "id","title":"ID","order":"true"},
		{"col":"name","title":"Name"},
		{"col":"hungarian","title":"Hungarian"},
		{"col":"actions","title":"Actions"},
	}

	for _, o := range elements {
		var actions []string;
		actions = append(actions, h.Replace(
			`<a href="%link%">%title%</a>`,
			[]string{"%link%","%title%"},
			[]string{h.GetUrl("owner/edit",  []string{strconv.Itoa(int(o.Id))},true,"admin"), "[Edit]"},
		))

		/*actions = append(actions, h.Replace(
			`<a href="%link%" onclick="return window.confirm('Are you sure you want to delete the item?')">%title%</a>`,
			[]string{"%link%","%title%"},
			[]string{h.GetUrl("owner/delete",  []string{strconv.Itoa(int(o.Id))},true,"admin"), "[Delete]"},
		))*/

		rows = append(rows, map[string]string{
			"id":         strconv.Itoa(int(o.Id)),
			"name":      o.Name,
			"hungarian": strconv.FormatBool(o.Hungarian),
			"actions":strings.Join(actions,"&nbsp;&nbsp;"),
		});
	}

	options = map[string]string{
		"class": "table-striped table-bordered table-hover",
		"id":    "owner-list-table",
	};
	return ol.List.Render(headers, rows, options);
}

func (ol *OwnerList) GetAll() []m.Owner {
	var results []m.Owner
	var where string = ol.List.GetSqlParams();
	if(where != ""){
		where = fmt.Sprintf(" WHERE %v",where);
	}
	sql := fmt.Sprintf("SELECT * FROM %v%v ORDER BY %v %v", ol.List.Table, where, ol.List.GetOrder(), ol.List.GetOrderDir());
	h.PrintlnIf(sql,h.GetConfig().Mode.Debug);
	_, err := db.DbMap.Select(&results, sql);
	h.Error(err, "", h.ERROR_LVL_ERROR)
	return results;
}

func (ol *OwnerList) GetToPage() []m.Owner {
	var results []m.Owner
	var where string = ol.List.GetSqlParams();
	if(where != ""){
		where = fmt.Sprintf(" WHERE %v",where);
	}
	sql := fmt.Sprintf("SELECT * FROM %v%v ORDER BY %v %v LIMIT %v", ol.List.Table, where, ol.List.GetOrder(), ol.List.GetOrderDir(), ol.List.GetLimitString());
	h.PrintlnIf(sql,h.GetConfig().Mode.Debug);
	_, err := db.DbMap.Select(&results, sql);
	h.Error(err, "", h.ERROR_LVL_ERROR)
	return results;
}
