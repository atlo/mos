package page

import (
	"mos/model"
	"mos/db"
	"fmt"
	"mos/helper"
	"strings"
	"html/template"
)

type Index struct {
	Blocks []model.Block
	Graph Graph
}

func (i *Index) Init(session *helper.Session){
	var blockIdentifiers []string = []string{"about","impressum","privacy","methodology"};
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

	db.DbMap.Select(&i.Blocks,fmt.Sprintf(query,block.GetTable()));

	i.Graph.Init(session);
}

func (i Index) GraphContent() template.HTML{
	return i.Graph.GetContent();
}
