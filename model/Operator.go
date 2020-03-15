package model

import (
	"mos/db"
	"fmt"
	"errors"
	h "mos/helper"
	"mos/model/FElement"
	"github.com/valyala/fasthttp"
	"reflect"
	"time"
	"os"
	"encoding/csv"
	"bufio"
	"io"
	"strconv"
	"strings"
)

type Operator struct {
	Id               int64     `db:"id, primarykey"`
	Name             string    `db:"name, size:255"`
	EvolutionDate    string `db:"evolution_date"`
	RegistrationDate string `db:"registration_date"`
	TerminationDate string `db:"termination_date"`
}

func (o Operator) GetAll() []Operator {
	var operators []Operator
	_, err := db.DbMap.Select(&operators, fmt.Sprintf("SELECT * FROM %s ORDER BY name", o.GetTable()));
	h.Error(err, "", h.ERROR_LVL_ERROR);
	return operators;
}

func (_ Operator) Get(id int64) (Operator, error) {
	var operator Operator;
	if (id == 0) {
		return operator, errors.New(fmt.Sprintf("Could not retrieve operator to ID %v", id));
	}

	err := db.DbMap.SelectOne(&operator, fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", operator.GetTable(), operator.GetPrimaryKey()[0]), id);
	h.Error(err, "", h.ERROR_LVL_ERROR);
	if (err != nil) {
		return operator, err;
	}

	if (operator.Id == 0) {
		return operator, errors.New(fmt.Sprintf("Could not retrieve operator to ID %v", id))
	}

	return operator, nil;
}

func (_ Operator) IsLanguageModel() bool {
	return false;
}

func (_ Operator) GetTable() string {
	return "operator";
}

func (_ Operator) GetPrimaryKey() []string {
	return []string{"id"};
}

func GetOperatorForm(data map[string]interface{}, action string) Form {
	var Elements []FormElement;
	var id = FElement.InputHidden{"id", "id", "", false, true, data["id"].(string)}
	Elements = append(Elements, id);
	var name = FElement.InputText{"Name", "name", "", "", "", false, false, data["name"].(string), "", "", "", "", ""}
	Elements = append(Elements, name);
	var evdate = FElement.InputText{"Evolution date", "evolution_date", "evolution_date", "", "fe.: " + time.Now().Format("2006-01-02"), false, false, data["evolution_date"].(string), "", "", "", "", ""}
	Elements = append(Elements, evdate);
	var regdate = FElement.InputText{"Registration date", "registration_date", "registration_date", "", "fe.: " + time.Now().Format("2006-01-02"), false, false, data["registration_date"].(string), "", "", "", "", ""}
	Elements = append(Elements, regdate);
	var termdate = FElement.InputText{"Termination date", "termination_date", "termination_date", "", "fe.: " + time.Now().Format("2006-01-02"), false, false, data["termination_date"].(string), "", "", "", "", ""}
	Elements = append(Elements, termdate);
	var fullColMap = map[string]string{"lg": "12", "md": "12", "sm": "12", "xs": "12"};
	var Fieldsets []Fieldset;
	Fieldsets = append(Fieldsets, Fieldset{"left", Elements, fullColMap});
	button := FElement.InputButton{"Submit", "submit", "submit", "pull-right", false, "", true, false, false, nil}
	Fieldsets = append(Fieldsets, Fieldset{"bottom", []FormElement{button}, fullColMap});
	var form = Form{h.GetUrl(action, nil, true, "admin"), "POST", false, Fieldsets, false, nil, nil}

	return form;
}

func NewOperator(Id int64, Name string, EvolutionDate string, RegistrationDate string) Operator {
	return Operator{
		Id:               Id,
		Name:             Name,
		EvolutionDate:    EvolutionDate,
		RegistrationDate: RegistrationDate,
	};
}

func NewEmptyOperator() Operator {
	return NewOperator(0, "", "", "")
}

func GetOperatorFormValidator(ctx *fasthttp.RequestCtx, Operator Operator) Validator {
	var Validator Validator;
	Validator = Validator.New(ctx);
	Validator.AddField("id", map[string]interface{}{
		"roles": map[string]interface{}{
			"required": false,
		},
	});
	Validator.AddField("evolution_date", map[string]interface{}{
		"roles": map[string]interface{}{
			"required": false,
			"format": map[string]interface{}{
				"type":   "regexp",
				"regexp": "^(\\d{4}-\\d{2}-\\d{2})+$",
			},
		},
	});

	Validator.AddField("registration_date", map[string]interface{}{
		"roles": map[string]interface{}{
			"required": false,
			"format": map[string]interface{}{
				"type":   "regexp",
				"regexp": "^(\\d{4}-\\d{2}-\\d{2})+$",
			},
		},
	});

	Validator.AddField("termination_date", map[string]interface{}{
		"roles": map[string]interface{}{
			"required": false,
			"format": map[string]interface{}{
				"type":   "regexp",
				"regexp": "^(\\d{4}-\\d{2}-\\d{2})+$",
			},
		},
	});
	Validator.AddField("name", map[string]interface{}{
		"roles": map[string]interface{}{
			"required": true,
		},
	});
	return Validator;
}

func (o Operator) GetOptions(defOption map[string]string) []map[string]string {
	var operators = o.GetAll();
	var options []map[string]string;
	if (defOption != nil) {
		_, okl := defOption["label"];
		_, okv := defOption["value"];
		if (okl || okv) {
			options = append(options, defOption);
		}
	}
	for _, operator := range operators {
		options = append(options, map[string]string{"label": operator.Name, "value": strconv.Itoa(int(operator.Id))});
	}
	return options;
}

func (o Operator) BuildStructure() {
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
			"name":   "IDX_OPERATOR_EVOLUTION_DATE",
			"type":   "hash",
			"field":  []string{"evolution_date"},
			"unique": false,
		},
		1: {
			"name":   "IDX_OPERATOR_REGISTRATION_DATE",
			"type":   "hash",
			"field":  []string{"registration_date"},
			"unique": false,
		},
	};
	tablemap, err := dbmap.TableFor(reflect.TypeOf(Operator{}), false);
	h.Error(err, "", h.ERROR_LVL_ERROR);
	for _, index := range indexes {
		h.PrintlnIf(fmt.Sprintf("Create %s index", index["name"].(string)), Conf.Mode.Rebuild_structure);
		tablemap.AddIndex(index["name"].(string), index["type"].(string), index["field"].([]string)).SetUnique(index["unique"].(bool));
	}

	dbmap.CreateIndex();
}

func (o Operator) PrepeareData(){
	h.PrintlnIf("Start importing operator data...",h.GetConfig().Mode.Debug);
	var csvFiles []string = []string{
		"Data/Uzemelteto_allando_attributumai_alapitasiev.csv",
	};
	for _, fileName := range csvFiles {
		h.PrintlnIf(fmt.Sprintf("Parsing file %s\r\n",fileName), h.GetConfig().Mode.Debug)

		csvFile, _ := os.Open(fileName)
		reader := csv.NewReader(bufio.NewReader(csvFile))
		var header []string;
		i := 0;
		for {
			i++;
			line, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				h.Error(err, "", h.ERROR_LVL_ERROR)
			}

			if(i == 1){
				header = line;
			}

			if (i == 1 || strings.Trim(line[1], " ") == "") {
				continue;
			}

			id, err := strconv.Atoi(line[0]);
			if (err != nil) {
				h.Error(err, "", h.ERROR_LVL_WARNING);
				continue;
			}

			if (line[2] == "") {
				line[2] = "01/01/1970";
			}

			if (line[3] == "") {
				line[3] = "01/01/1970";
			}

			evDate, err := time.Parse("01/02/2006", line[2]);
			if (err != nil) {
				h.Error(err, "", h.ERROR_LVL_WARNING);
			}

			regDate, err := time.Parse("01/02/2006", line[3]);
			if (err != nil) {
				h.Error(err, "", h.ERROR_LVL_WARNING);
			}

			var operator Operator;

			operator = Operator{
				int64(id),
				line[1],
				evDate.Format("2006-01-02"),
				regDate.Format("2006-01-02"),
				"",
			}

			err = db.DbMap.Insert(&operator);
			h.Error(err, "", h.ERROR_LVL_ERROR);

			var k int = 4;
			for(k<len(line)){
				var oyd OperatorYearData;
				var headerY int;
				headerY,err = strconv.Atoi(header[k]);
				oyd,err = oyd.Get(operator.Id,headerY);
				var oydNew bool = oyd.OperatorId == 0;
				oyd.OperatorId = operator.Id;
				oyd.Year = headerY;
				oyd.GovernmentMedia = strings.Trim(line[k]," ") == "I"
				if(oydNew){
					err = db.DbMap.Insert(&oyd);
				} else {
					_, err = db.DbMap.Update(&oyd);
				}
				h.Error(err, "", h.ERROR_LVL_ERROR);
				k++
			}
		}
	}
	csvFiles = []string{
		"Data/uzemelteto_allando_attributumai_alapitas_eve_1998_2009.csv",
		"Data/Uzemelteto_allando_attributumai_alapitasiev_2017is.csv",
	};
	for _, fileName := range csvFiles {
		h.PrintlnIf(fmt.Sprintf("Parsing file %s\r\n",fileName), h.GetConfig().Mode.Debug)
		csvFile, _ := os.Open(fileName)
		reader := csv.NewReader(bufio.NewReader(csvFile))
		i := 0;
		for {
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

			if (line[2] == "") {
				line[2] = "01/01/1970";
			}

			/*if (line[3] == "") {
				line[3] = "01/01/1970";
			}*/

			evDate, err := time.Parse("01/02/2006", line[2]);
			if (err != nil) {
				h.Error(err, "", h.ERROR_LVL_WARNING);
			}

			regDate,err := time.Parse("01/02/2006", "01/01/1970");
			/*regDate, err := time.Parse("01/02/2006", line[3]);
			if (err != nil) {
				h.Error(err, "", h.ERROR_LVL_WARNING);
			}*/

			var operator Operator;
			operator,err = operator.Get(int64(id));

			if(operator.Id == 0) {
				operator = Operator{
					int64(id),
					line[1],
					evDate.Format("2006-01-02"),
					regDate.Format("2006-01-02"),
					"",
				}

				err = db.DbMap.Insert(&operator);
				h.Error(err, "", h.ERROR_LVL_ERROR);
			}
		}
	}
	h.PrintlnIf("Done importing operator data...",h.GetConfig().Mode.Debug);
}
