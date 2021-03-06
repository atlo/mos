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

type BlockList struct {
	List m.List
}

func (bl *BlockList) Init(ctx *fasthttp.RequestCtx, lang string) {
	var block m.Block;
	bl.List.Init(ctx, block, lang)

	bl.List.AddSearchParam(m.SearchParam{"ID","id","number","id", nil});
	bl.List.AddSearchParam(m.SearchParam{"Identifier","identifier","text","identifier",nil});
	bl.List.AddSearchParam(m.SearchParam{"Content","content","text","content",nil});
}

func (bl *BlockList) SetLimitParam(limitParam string) {
	bl.List.SetLimitParam(limitParam)
}

func (bl *BlockList) SetPageParam(pageParam string) {
	bl.List.SetPageParam(pageParam);
}

func (bl *BlockList) Render(elements []m.Block) string {
	var headers []map[string]string;
	var rows []map[string]string;
	var options map[string]string;
	headers = []map[string]string{
		{"col": "id","title":"ID"},
		{"col":"identifier","title":"Identifier"},
		{"col":"content","title":"Content"},
		{"col":"actions","title":"Actions"},
	}

	for _, b := range elements {
		var actions []string;
		actions = append(actions, h.Replace(
			`<a href="%link%">%title%</a>`,
			[]string{"%link%","%title%"},
			[]string{h.GetUrl("block/edit",  []string{strconv.Itoa(int(b.Id))},true,"admin"), "[Edit]"},
		))

		actions = append(actions, h.Replace(
			`<a href="%link%" onclick="return window.confirm('Are you sure you want to delete the item?')">%title%</a>`,
			[]string{"%link%","%title%"},
			[]string{h.GetUrl("block/delete",  []string{strconv.Itoa(int(b.Id))},true,"admin"), "[Delete]"},
		))

		rows = append(rows, map[string]string{
			"id":         strconv.Itoa(int(b.Id)),
			"identifier":      b.Identifier,
			"content": b.Content,
			"actions":strings.Join(actions,"&nbsp;&nbsp;"),
		});
	}

	options = map[string]string{
		"class": "table-striped table-bordered table-hover",
		"id":    "block-list-table",
	};
	return bl.List.Render(headers, rows, options);
}

func (bl *BlockList) GetAll() []m.Block {
	var results []m.Block
	var where string = bl.List.GetSqlParams();
	if(where != ""){
		where = fmt.Sprintf(" WHERE %v",where);
	}
	sql := fmt.Sprintf("SELECT * FROM %v%v ORDER BY %v %v", bl.List.Table, where, bl.List.GetOrder(), bl.List.GetOrderDir());
	h.PrintlnIf(sql,h.GetConfig().Mode.Debug);
	_, err := db.DbMap.Select(&results, sql);
	h.Error(err, "", h.ERROR_LVL_ERROR)
	return results;
}

func (bl *BlockList) GetToPage() []m.Block {
	var results []m.Block
	var where string = bl.List.GetSqlParams();
	if(where != ""){
		where = fmt.Sprintf(" WHERE %v",where);
	}
	sql := fmt.Sprintf("SELECT * FROM %v%v ORDER BY %v %v LIMIT %v", bl.List.Table, where, bl.List.GetOrder(), bl.List.GetOrderDir(), bl.List.GetLimitString());
	h.PrintlnIf(sql,h.GetConfig().Mode.Debug);
	_, err := db.DbMap.Select(&results, sql);
	h.Error(err, "", h.ERROR_LVL_ERROR)
	return results;
}
