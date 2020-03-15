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

type MediaOperatorList struct {
	List m.List
}

func (mol *MediaOperatorList) Init(ctx *fasthttp.RequestCtx, lang string) {
	var mediaOperator m.MediaOperator;
	mol.List.Init(ctx, mediaOperator, lang)

	var media m.Media;
	mediaList := media.GetOptions(map[string]string{"value":"","label":"Anything"});

	var operator m.Operator;
	opaeratorList := operator.GetOptions(map[string]string{"value":"","label":"Anything"});

	mol.List.AddSearchParam(m.SearchParam{"Media", "media_id", "select", "media_id", map[string]interface{}{"options": mediaList}});
	mol.List.AddSearchParam(m.SearchParam{"Operator", "operator_id", "select", "operator_id", map[string]interface{}{"options": opaeratorList}});
	mol.List.AddSearchParam(m.SearchParam{"Year","year","number","year",nil});
}

func (mol *MediaOperatorList) SetLimitParam(limitParam string) {
	mol.List.SetLimitParam(limitParam)
}

func (mol *MediaOperatorList) SetPageParam(pageParam string) {
	mol.List.SetPageParam(pageParam);
}

func (mol *MediaOperatorList) Render(elements []m.MediaOperator) string {
	var headers []map[string]string;
	var rows []map[string]string;
	var options map[string]string;
	headers = []map[string]string{
		{"col": "year","title":"Year"},
		{"col":"media_id","title":"Media"},
		{"col":"operator_id","title":"Operator"},
		{"col":"actions","title":"Actions"},
	}

	for _, mo:= range elements {
		media := mo.GetMedia();
		operator := mo.GetOperator();
		var actions []string;

		actions = append(actions, h.Replace(
			`<a href="%link%" onclick="return window.confirm('Are you sure you want to delete the item?')">%title%</a>`,
			[]string{"%link%","%title%"},
			[]string{h.GetUrl("mediaoperator/delete",  []string{strconv.Itoa(int(mo.Id))},true,"admin"), "[Delete]"},
		))

		rows = append(rows, map[string]string{
			"year":      strconv.Itoa(mo.Year),
			"media_id":      media.Name,
			"operator_id":      operator.Name,
			"actions":strings.Join(actions,"&nbsp;&nbsp;"),
		});
	}

	options = map[string]string{
		"class": "table-striped table-bordered table-hover",
		"id":    "mediaoperator-list-table",
	};
	return mol.List.Render(headers, rows, options);
}

func (mol *MediaOperatorList) GetAll() []m.MediaOperator {
	var results []m.MediaOperator
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

func (mol *MediaOperatorList) GetToPage() []m.MediaOperator {
	var results []m.MediaOperator
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
