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

type MediaOwnerList struct {
	List m.List
}

func (mol *MediaOwnerList) Init(ctx *fasthttp.RequestCtx, lang string) {
	var mediaOwner m.MediaOwner;
	mol.List.Init(ctx, mediaOwner, lang)

	var media m.Media;
	mediaList := media.GetOptions(map[string]string{"value":"","label":"Anything"});

	var owner m.Owner;
	ownerList := owner.GetOptions(map[string]string{"value":"","label":"Anything"});

	mol.List.AddSearchParam(m.SearchParam{"Media", "media_id", "select", "media_id", map[string]interface{}{"options": mediaList}});
	mol.List.AddSearchParam(m.SearchParam{"Owner", "owner_id", "select", "owner_id", map[string]interface{}{"options": ownerList}});
	mol.List.AddSearchParam(m.SearchParam{"Year","year","number","year",nil});
}

func (mol *MediaOwnerList) SetLimitParam(limitParam string) {
	mol.List.SetLimitParam(limitParam)
}

func (mol *MediaOwnerList) SetPageParam(pageParam string) {
	mol.List.SetPageParam(pageParam);
}

func (mol *MediaOwnerList) Render(elements []m.MediaOwner) string {
	var headers []map[string]string;
	var rows []map[string]string;
	var options map[string]string;
	headers = []map[string]string{
		{"col": "year","title":"Year"},
		{"col":"media_id","title":"Media"},
		{"col":"owner_id","title":"Owner"},
		{"col":"actions","title":"Actions"},
	}

	for _, mo:= range elements {
		media := mo.GetMedia();
		owner := mo.GetOwner();
		var actions []string;

		actions = append(actions, h.Replace(
			`<a href="%link%" onclick="return window.confirm('Are you sure you want to delete the item?')">%title%</a>`,
			[]string{"%link%","%title%"},
			[]string{h.GetUrl("mediaowner/delete",  []string{strconv.Itoa(int(mo.Id))},true,"admin"), "[Delete]"},
		))

		rows = append(rows, map[string]string{
			"year":      strconv.Itoa(mo.Year),
			"media_id":      media.Name,
			"owner_id":      owner.Name,
			"actions":strings.Join(actions,"&nbsp;&nbsp;"),
		});
	}

	options = map[string]string{
		"class": "table-striped table-bordered table-hover",
		"id":    "mediaowner-list-table",
	};
	return mol.List.Render(headers, rows, options);
}

func (mol *MediaOwnerList) GetAll() []m.MediaOwner {
	var results []m.MediaOwner
	var where string = mol.List.GetSqlParams();
	if(where != ""){
		where = fmt.Sprintf(" WHERE %v",where);
	}
	sql := fmt.Sprintf("SELECT * FROM %v%v ORDER BY %v %v", mol.List.Table, where, mol.List.GetOrder(), mol.List.GetOrderDir());
	h.PrintlnIf(sql,h.GetConfig().Mode.Debug);
	_, err := db.DbMap.Select(&results, sql);
	h.Error(err, "", h.ERROR_LVL_ERROR)
	return results;
}

func (mol *MediaOwnerList) GetToPage() []m.MediaOwner {
	var results []m.MediaOwner
	var where string = mol.List.GetSqlParams();
	if(where != ""){
		where = fmt.Sprintf(" WHERE %v",where);
	}
	sql := fmt.Sprintf("SELECT * FROM %v%v ORDER BY %v %v LIMIT %v", mol.List.Table, where, mol.List.GetOrder(), mol.List.GetOrderDir(), mol.List.GetLimitString());
	h.PrintlnIf(sql,h.GetConfig().Mode.Debug);
	_, err := db.DbMap.Select(&results, sql);
	h.Error(err, "", h.ERROR_LVL_ERROR)
	return results;
}
