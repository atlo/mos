package model

import (
	"github.com/go-gorp/gorp"
	"fmt"
	h "mos/helper"
	"mos/db"
	"strconv"
	"errors"
	"database/sql"
	"github.com/valyala/fasthttp"
	"mos/model/FElement"
)

type Interest struct {
	Id   int64  `db:"id, primarykey, autoincrement";json:"id"`
	Name string `db:"name, size:255";json:"name"`
}

// implement the PreInsert and PreUpdate hooks
func (m *Interest) PreInsert(sg gorp.SqlExecutor) error {
	return nil
}

func (m *Interest) PreUpdate(sg gorp.SqlExecutor) error {
	return nil
}

func NewInterest(Id int64, Name string) Interest {
	return Interest{
		Id:   Id,
		Name: Name,
	};
}

func NewEmptyInterest() Interest {
	return NewInterest(0, "")
}

func (_ Interest) Get(id int64) (Interest, error) {
	var interest Interest;
	if (id == 0) {
		return interest, errors.New(fmt.Sprintf("Could not retrieve Interest to ID %v", id));
	}

	err := db.DbMap.SelectOne(&interest, fmt.Sprintf("SELECT * FROM %s WHERE %v = ?", interest.GetTable(), interest.GetPrimaryKey()[0]), id);
	h.Error(err, "", h.ERROR_LVL_ERROR);
	if (err != nil) {
		return interest, err;
	}

	if (interest.Id == 0) {
		return interest, errors.New(fmt.Sprintf("Could not retrieve Media Type to ID %v", id))
	}

	return interest, nil;
}

func (m Interest) GetAll() []Interest {
	var interests []Interest
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY %s", m.GetTable(), m.GetPrimaryKey()[0]);
	_, err := db.DbMap.Select(&interests, query);
	h.Error(err, "", h.ERROR_LVL_ERROR);
	return interests;
}

func (m Interest) SaveAndGetByName(name string) (Interest, error) {
	var interest Interest;
	if (name == "") {
		return interest, errors.New("Could not retrieve Interest to empty name");
	}

	err := db.DbMap.SelectOne(&interest, fmt.Sprintf(`SELECT * FROM %s WHERE %s = ?`, interest.GetTable(), "name"), name);
	if (err == sql.ErrNoRows) {
		interest.Name = name;
		err = db.DbMap.Insert(&interest);
	}

	return interest, err;
}

func (m Interest) GetOptions(defOption map[string]string) []map[string]string {
	var interestes = m.GetAll();
	var options []map[string]string;
	if (defOption != nil) {
		_, okl := defOption["label"];
		_, okv := defOption["value"];
		if (okl || okv) {
			options = append(options, defOption);
		}
	}
	for _, interest := range interestes {
		options = append(options, map[string]string{"label": interest.Name, "value": strconv.Itoa(int(interest.Id))});
	}
	return options;
}

func (m Interest) BuildStructure() {
	dbmap := db.DbMap;
	Conf := h.GetConfig();

	if (!Conf.Mode.Rebuild_structure) {
		return;
	}

	h.PrintlnIf(fmt.Sprintf("Drop %v table", m.GetTable()), Conf.Mode.Rebuild_structure);
	dbmap.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s;", m.GetTable()));
	h.PrintlnIf(fmt.Sprintf("Create %v table", m.GetTable()), Conf.Mode.Rebuild_structure);
	dbmap.CreateTablesIfNotExists();
}

func GetInterestForm(data map[string]interface{}, action string) Form {
	var Elements []FormElement;
	var id = FElement.InputHidden{"id", "id", "", false, true, data["id"].(string)}
	Elements = append(Elements, id);
	var identifier = FElement.InputText{"Name", "name", "", "", "", false, false, data["name"].(string), "", "", "", "", ""}
	Elements = append(Elements, identifier);
	var fullColMap = map[string]string{"lg": "12", "md": "12", "sm": "12", "xs": "12"};
	var Fieldsets []Fieldset;
	Fieldsets = append(Fieldsets, Fieldset{"left", Elements, fullColMap});
	button := FElement.InputButton{"Submit", "submit", "submit", "pull-right", false, "", true, false, false, nil}
	Fieldsets = append(Fieldsets, Fieldset{"bottom", []FormElement{button}, fullColMap});
	var form = Form{h.GetUrl(action, nil, true, "admin"), "POST", false, Fieldsets, false, nil, nil}

	return form;
}

func GetInterestFormValidator(ctx *fasthttp.RequestCtx, interest Interest) Validator {
	var Validator Validator;
	Validator = Validator.New(ctx);
	Validator.AddField("id", map[string]interface{}{
		"roles": map[string]interface{}{
			"required": false,
		},
	});
	Validator.AddField("name", map[string]interface{}{
		"roles": map[string]interface{}{
			"required": true,
		},
	});
	return Validator;
}

func (_ Interest) IsLanguageModel() bool {
	return false;
}

func (_ Interest) GetTable() string {
	return "interest";
}

func (_ Interest) GetPrimaryKey() []string {
	return []string{"id"};
}

func (m Interest) PrepeareData(){}
