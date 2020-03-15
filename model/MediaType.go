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
	"strings"
)

type MediaType struct {
	Id   int64  `db:"id, primarykey, autoincrement"`
	Name string `db:"name, size:255"`
}

// implement the PreInsert and PreUpdate hooks
func (m *MediaType) PreInsert(sg gorp.SqlExecutor) error {
	return nil
}

func (m *MediaType) PreUpdate(sg gorp.SqlExecutor) error {
	return nil
}

func NewMediaType(Id int64, Name string) MediaType {
	return MediaType{
		Id:   Id,
		Name: Name,
	};
}

func (m MediaType) GetColorClass() string {
	var typeClean string = strings.Trim(strings.ToLower(h.RemoveAccents(m.Name)), " ");
	switch (typeClean) {
	case "sajto":
		return "color--news";
		break;
	case "televizio":
		return "color--tv";
		break;
	case "internet":
		return "color--online";
		break;
	case "radio":
		return "color--radio";
		break;
	}

	return "";
}

func NewEmptyMediaType() MediaType {
	return NewMediaType(0, "")
}

func (_ MediaType) Get(id int64) (MediaType, error) {
	var mediaType MediaType;
	if (id == 0) {
		return mediaType, errors.New(fmt.Sprintf("Could not retrieve Media Type to ID %v", id));
	}

	err := db.DbMap.SelectOne(&mediaType, fmt.Sprintf("SELECT * FROM %s WHERE %v = ?", mediaType.GetTable(), mediaType.GetPrimaryKey()[0]), id);
	h.Error(err, "", h.ERROR_LVL_ERROR);
	if (err != nil) {
		return mediaType, err;
	}

	if (mediaType.Id == 0) {
		return mediaType, errors.New(fmt.Sprintf("Could not retrieve Media Type to ID %v", id))
	}

	return mediaType, nil;
}

func (m MediaType) GetAll() []MediaType {
	var mediaTypees []MediaType
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY %s", m.GetTable(), m.GetPrimaryKey()[0]);
	h.PrintlnIf(query, h.GetConfig().Mode.Debug);
	_, err := db.DbMap.Select(&mediaTypees, query);
	h.Error(err, "", h.ERROR_LVL_ERROR);
	return mediaTypees;
}

func (m MediaType) SaveAndGetByName(name string) (MediaType, error) {
	var mediaType MediaType;
	if (name == "") {
		return mediaType, errors.New("Could not retrieve Media Type to empty name");
	}

	err := db.DbMap.SelectOne(&mediaType, fmt.Sprintf(`SELECT * FROM %s WHERE %s = ?`, mediaType.GetTable(), "name"), name);
	if (err == sql.ErrNoRows) {
		mediaType.Name = name;
		err = db.DbMap.Insert(&mediaType);
	}

	return mediaType, err;
}

func (m MediaType) GetOptions(defOption map[string]string) []map[string]string {
	var mediaTypees = m.GetAll();
	var options []map[string]string;
	if (defOption != nil) {
		_, okl := defOption["label"];
		_, okv := defOption["value"];
		if (okl || okv) {
			options = append(options, defOption);
		}
	}
	for _, mediaType := range mediaTypees {
		options = append(options, map[string]string{"label": mediaType.Name, "value": strconv.Itoa(int(mediaType.Id))});
	}
	return options;
}

func (m MediaType) BuildStructure() {
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

func GetMediaTypeForm(data map[string]interface{}, action string) Form {
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

func GetMediaTypeFormValidator(ctx *fasthttp.RequestCtx, mediaType MediaType) Validator {
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

func (_ MediaType) IsLanguageModel() bool {
	return false;
}

func (_ MediaType) GetTable() string {
	return "media_type";
}

func (_ MediaType) GetPrimaryKey() []string {
	return []string{"id"};
}

func (m MediaType) PrepeareData(){}
