package model

import (
	"mos/db"
	"fmt"
	"errors"
	h "mos/helper"
	"mos/model/FElement"
	"github.com/valyala/fasthttp"
	"reflect"
	"os"
	"encoding/csv"
	"bufio"
	"io"
	"strings"
	"strconv"
)

type Owner struct {
	Id        int64  `db:"id, primarykey"`
	Name      string `db:"name, size:255"`
	Hungarian bool   `db:"hungarian"`
}

func (o Owner) GetOptions(defOption map[string]string) []map[string]string {
	var owners = o.GetAll();
	var options []map[string]string;
	if (defOption != nil) {
		_, okl := defOption["label"];
		_, okv := defOption["value"];
		if (okl || okv) {
			options = append(options, defOption);
		}
	}
	for _, owner := range owners {
		options = append(options, map[string]string{"label": owner.Name, "value": strconv.Itoa(int(owner.Id))});
	}
	return options;
}

func (o Owner) GetAll() []Owner {
	var owners []Owner
	_, err := db.DbMap.Select(&owners, fmt.Sprintf("select * from %s order by %v", o.GetTable(), "name ASC"));
	h.Error(err, "", h.ERROR_LVL_ERROR);
	return owners;
}

func (_ Owner) Get(id int64) (Owner, error) {
	var owner Owner;
	if (id == 0) {
		return owner, errors.New(fmt.Sprintf("Could not retrieve owner to ID %v", id));
	}

	err := db.DbMap.SelectOne(&owner, fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", owner.GetTable(), owner.GetPrimaryKey()[0]), id);
	h.Error(err, "", h.ERROR_LVL_ERROR);
	if (err != nil) {
		return owner, err;
	}

	if (owner.Id == 0) {
		return owner, errors.New(fmt.Sprintf("Could not retrieve owner to ID %v", id))
	}

	return owner, nil;
}

func (_ Owner) GetByName(name string) (Owner, error) {
	var owner Owner;
	if (name == "") {
		return owner, errors.New(fmt.Sprintf("Could not retrieve owner to name %s", name));
	}

	err := db.DbMap.SelectOne(&owner, fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", owner.GetTable(), "name"), name);
	h.Error(err, "", h.ERROR_LVL_ERROR);
	if (err != nil) {
		return owner, err;
	}

	if (owner.Id == 0) {
		return owner, errors.New(fmt.Sprintf("Could not retrieve owner to name %s", name))
	}

	return owner, nil;
}

func (_ Owner) IsLanguageModel() bool {
	return false;
}

func (_ Owner) GetTable() string {
	return "owner";
}

func (_ Owner) GetPrimaryKey() []string {
	return []string{"id"};
}

func GetOwnerForm(data map[string]interface{}, action string) Form {
	var Elements []FormElement;
	var id = FElement.InputHidden{"id", "id", "", false, true, data["id"].(string)}
	Elements = append(Elements, id);
	var identifier = FElement.InputText{"Name", "name", "name", "", "", false, false, data["name"].(string), "", "", "", "", ""}
	Elements = append(Elements, identifier);

	Checkboxes := FElement.CheckboxGroup{};
	checkbox := FElement.InputCheckbox{
		"Magyar",
		"hungarian",
		"hungarian",
		"",
		false,
		false,
		"1",
		[]string{data["hungarian"].(string)},
		true,
	}
	Checkboxes.Checkbox = append(Checkboxes.Checkbox, checkbox);
	Elements = append(Elements, Checkboxes);

	var fullColMap = map[string]string{"lg": "12", "md": "12", "sm": "12", "xs": "12"};
	var Fieldsets []Fieldset;
	Fieldsets = append(Fieldsets, Fieldset{"left", Elements, fullColMap});
	button := FElement.InputButton{"Submit", "submit", "submit", "pull-right", false, "", true, false, false, nil}
	Fieldsets = append(Fieldsets, Fieldset{"bottom", []FormElement{button}, fullColMap});
	var form = Form{h.GetUrl(action, nil, true, "admin"), "POST", false, Fieldsets, false, nil, nil}

	return form;
}

func NewOwner(Id int64, Name string, Hungarian bool) Owner {
	return Owner{
		Id:               Id,
		Name:             Name,
		Hungarian:Hungarian,
	};
}

func NewEmptyOwner() Owner {
	return NewOwner(0, "", false)
}

func GetOwnerFormValidator(ctx *fasthttp.RequestCtx, Owner Owner) Validator {
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

func (o Owner) BuildStructure() {
	dbmap := db.DbMap;
	Conf := h.GetConfig();

	if(!Conf.Mode.Rebuild_structure){
		return;
	}

	h.PrintlnIf(fmt.Sprintf("Drop %v table", o.GetTable()), Conf.Mode.Rebuild_structure);
	dbmap.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s;", o.GetTable()));

	h.PrintlnIf(fmt.Sprintf("Create %v table", o.GetTable()), Conf.Mode.Rebuild_structure);
	dbmap.CreateTablesIfNotExists();
	var indexes map[int]map[string]interface{} = make(map[int]map[string]interface{})

	indexes = map[int]map[string]interface{}{
		0: {
			"name":   "IDX_OWNER_HUNGARIAN",
			"type":   "hash",
			"field":  []string{"hungarian"},
			"unique": false,
		},
	};
	tablemap, err := dbmap.TableFor(reflect.TypeOf(Owner{}), false);
	h.Error(err, "", h.ERROR_LVL_ERROR);
	for _, index := range indexes {
		h.PrintlnIf(fmt.Sprintf("Create %s index", index["name"].(string)), Conf.Mode.Rebuild_structure);
		tablemap.AddIndex(index["name"].(string), index["type"].(string), index["field"].([]string)).SetUnique(index["unique"].(bool));
	}

	dbmap.CreateIndex();
}

func (o Owner) PrepeareData(){
	h.PrintlnIf("Start importing owner data...",h.GetConfig().Mode.Debug);
	var csvFiles []string = []string{
		"Data/Tulajdonosok_allando_attributumai.csv",
		"Data/Vegtulajdonos_allando_attributumai_2017is.csv",
	};
	for _, fileName := range csvFiles {
		h.PrintlnIf(fmt.Sprintf("Parsing file %s\r\n",fileName), h.GetConfig().Mode.Debug)
		csvFile, _ := os.Open(fileName)
		reader := csv.NewReader(bufio.NewReader(csvFile))
		i := 0;
		for {
			var owner Owner;
			i++;
			line, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				h.Error(err, "", h.ERROR_LVL_ERROR)
			}

			if (i == 1 || strings.Trim(line[1], " ") == "") {
				continue;
			}

			id, err := strconv.Atoi(line[0]);
			if (err != nil) {
				h.Error(err, "", h.ERROR_LVL_WARNING);
				continue;
			}

			owner, err = owner.Get(int64(id));
			if (owner.Id > 0) {
				owner.Name = line[1];
				owner.Hungarian = strings.Trim(line[2], " ") == "1";
				_, err = db.DbMap.Update(&owner);
			} else {
				owner = Owner{
					int64(id),
					line[1],
					strings.Trim(line[2], " ") == "1",
				}
				err = db.DbMap.Insert(&owner);
			}
			h.Error(err, "", h.ERROR_LVL_ERROR);
		}
	}
	h.PrintlnIf("Done importing owner data...",h.GetConfig().Mode.Debug);
}
