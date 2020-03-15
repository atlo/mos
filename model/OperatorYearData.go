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
	"strconv"
	"strings"
	"regexp"
	"time"
)

type OperatorYearData struct {
	OperatorId        int64   `db:"operator_id"`
	Year              int     `db:"year, size:4"`
	Address           string  `db:"address, size: 500"`
	IncomeNet         int64 `db:"income_net"`
	IncomeTax         int64 `db:"income_tax"`
	IncomeOperational int64 `db:"income_operational"`
	GovernmentMedia bool `db:"government, default: 0"`
}

func (o OperatorYearData) GetAll() []OperatorYearData {
	var operatoryds []OperatorYearData
	_, err := db.DbMap.Select(&operatoryds, fmt.Sprintf("select * from %s order by operator_id, year", o.GetTable()));
	h.Error(err, "", h.ERROR_LVL_ERROR);
	return operatoryds;
}

func (oyd OperatorYearData) GetOperator() Operator{
	var operator Operator;
	operator,err := operator.Get(oyd.OperatorId);
	h.Error(err,"",h.ERROR_LVL_WARNING);
	return operator;
}

func (_ OperatorYearData) Get(operatorId int64, year int) (OperatorYearData, error) {
	var operatoryd OperatorYearData;
	if (operatorId == 0 || year == 0) {
		return operatoryd, errors.New(fmt.Sprintf("Could not retrieve operatoryear to Operator Id %v and year %v", operatorId, year));
	}

	query := fmt.Sprintf("SELECT * FROM %v WHERE %s = ? and %s = ?", operatoryd.GetTable(), "operator_id", "year");
	err := db.DbMap.SelectOne(&operatoryd, query, operatorId, year);
	h.Error(err, "", h.ERROR_LVL_ERROR);
	if (err != nil) {
		return operatoryd, err;
	}

	if (operatoryd.OperatorId == 0) {
		return operatoryd, errors.New(fmt.Sprintf("Could not retrieve operator data to operator id %v and year %v", operatorId, year))
	}

	return operatoryd, nil;
}

func (_ OperatorYearData) IsLanguageModel() bool {
	return false;
}

func (_ OperatorYearData) GetTable() string {
	return "operator_yeardata";
}

func (_ OperatorYearData) GetPrimaryKey() []string {
	return []string{"operator_id", "year"};
}

func GetOperatorYearDataForm(newModel bool, data map[string]interface{}, action string, actionParams ...string) Form {
	var operator Operator;
	var options = operator.GetOptions(nil);

	var halfColmap map[string]string = map[string]string{"lg":"6","md":"6","sm":"12","xs":"12"};
	var fullColMap = map[string]string{"lg": "12", "md": "12", "sm": "12", "xs": "12"};

	var FieldsetLeft Fieldset = Fieldset{"left",nil,halfColmap}
	var FieldsetRight Fieldset = Fieldset{"right",nil,halfColmap}
	var FieldsetBottom Fieldset = Fieldset{"bottom",nil, fullColMap};

	if(newModel){
		var operatorInp = FElement.InputSelect{"Operator", "operator_id", "operator_id", "", false, false, []string{data["operator_id"].(string)}, false, options, ""}
		FieldsetLeft.AddElement(operatorInp);
	} else {
		operatorId,err := strconv.Atoi(data["operator_id"].(string));
		h.Error(err,"",h.ERROR_LVL_WARNING);

		var operator Operator;
		operator,err = operator.Get(int64(operatorId));

		var operatorInp = FElement.InputHidden{"operator_id","operator_id","",false,true,data["operator_id"].(string)}
		FieldsetLeft.AddElement(operatorInp);

		var operatorNameInp = FElement.InputText{"Operator","","","","",false,true,operator.Name,"","","","",""}
		FieldsetLeft.AddElement(operatorNameInp);
	}

	var year = FElement.InputText{"year", "year", "", "", "", false, !newModel, data["year"].(string), "", "", "", "", ""}
	FieldsetLeft.AddElement(year);

	var address = FElement.InputText{"address", "address", "", "", "", false, false, data["address"].(string), "", "", "", "", ""}
	FieldsetLeft.AddElement(address);

	var incomeNetInp = FElement.InputText{"Income Without tax", "income_net", "", "", "", false, false, data["income_net"].(string), "", "", "", "", ""}
	FieldsetRight.AddElement(incomeNetInp);

	var incomeTaxInp = FElement.InputText{"Income With tax", "income_tax", "", "", "", false, false, data["income_tax"].(string), "", "", "", "", ""}
	FieldsetRight.AddElement(incomeTaxInp);

	var incomeOpeInp = FElement.InputText{"Income Operational", "income_operational", "", "", "", false, false, data["income_operational"].(string), "", "", "", "", ""}
	FieldsetRight.AddElement(incomeOpeInp);


	button := FElement.InputButton{"Submit", "submit", "submit", "pull-right", false, "", true, false, false, nil}
	FieldsetBottom.AddElement(button);

	var form = Form{h.GetUrl(action, actionParams, true, "admin"), "POST", false, []Fieldset{FieldsetLeft,FieldsetRight,FieldsetBottom}, false, nil, nil}

	return form;
}

func NewOperatorYearData(OperatorId int64, Year int, Address string) OperatorYearData {
	return OperatorYearData{
		OperatorId: OperatorId,
		Year:       Year,
		Address:    Address,
	};
}

func GetOperatorYearDataFormValidator(ctx *fasthttp.RequestCtx, OperatorYearData OperatorYearData) Validator {
	var Validator Validator;
	Validator = Validator.New(ctx);
	Validator.AddField("operator_id", map[string]interface{}{
		"roles": map[string]interface{}{
			"required": false,
		},
	});
	Validator.AddField("year", map[string]interface{}{
		"roles": map[string]interface{}{
			"required": true,
		},
	});
	Validator.AddField("income_net", map[string]interface{}{
		"roles": map[string]interface{}{
			"required": false,
			"format": map[string]interface{}{
				"type": "regexp",
				"regexp":"^-?\\d*$",
			},
		},
	});
	Validator.AddField("income_tax", map[string]interface{}{
		"roles": map[string]interface{}{
			"required": false,
			"format": map[string]interface{}{
				"type": "regexp",
				"regexp":"^-?\\d*$",
			},
		},
	});
	Validator.AddField("income_operational", map[string]interface{}{
		"roles": map[string]interface{}{
			"required": false,
			"format": map[string]interface{}{
				"type": "regexp",
				"regexp":"^-?\\d*$",
			},
		},
	});
	return Validator;
}

func (o OperatorYearData) BuildStructure() {
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
			"name":   "IDX_OPERATOR_YEARDATA_OPERATOR_ID_YEAR",
			"type":   "hash",
			"field":  []string{"operator_id", "year"},
			"unique": true,
		},
	};
	tablemap, err := dbmap.TableFor(reflect.TypeOf(OperatorYearData{}), false);
	h.Error(err, "", h.ERROR_LVL_ERROR);
	for _, index := range indexes {
		h.PrintlnIf(fmt.Sprintf("Create %s index", index["name"].(string)), Conf.Mode.Rebuild_structure);
		tablemap.AddIndex(index["name"].(string), index["type"].(string), index["field"].([]string)).SetUnique(index["unique"].(bool));
	}

	dbmap.CreateIndex();
}

func (o OperatorYearData) PrepeareData() {
	h.PrintlnIf("Start importing operator year address data...",h.GetConfig().Mode.Debug);
	var csvFiles []string = []string{
		"Data/uzemelteto_cim_idosoros_1998_2009.csv",
		"Data/Uzemelteto_cim_idosoros_2017is.csv",
	};
	var headers []string;
	_ = headers;
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

			if (i == 1) {
				headers = line;
				continue;
			} else if (strings.Trim(line[1], " ") == "") {
				continue;
			}

			id, err := strconv.Atoi(line[0]);
			if (err != nil) {
				h.Error(err, "", h.ERROR_LVL_WARNING);
				continue;
			}

			var l = 3;
			for (l < len(headers) && headers[l]!="") {
				y, err := strconv.Atoi(headers[l]);
				h.Error(err, "", h.ERROR_LVL_WARNING);
				var operatoryd OperatorYearData;
				operatoryd,_ = operatoryd.Get(int64(id),y)

				if(operatoryd.OperatorId == 0) {
					operatoryd = OperatorYearData{
						int64(id),
						y,
						line[l],
						0,
						0,
						0,
						false,
					}
					err = db.DbMap.Insert(&operatoryd);
				} else {
					operatoryd.Address = line[l];
					_,err = db.DbMap.Update(&operatoryd);
				}
				h.Error(err, "", h.ERROR_LVL_ERROR);
				l++;
			}
		}
	}
	h.PrintlnIf("Done importing operator year address data...",h.GetConfig().Mode.Debug);

	h.PrintlnIf("Start importing operator year income data...",h.GetConfig().Mode.Debug);
	csvFiles = []string{
		"Data/arbevetel_idosoros_1998_2009.csv",
		"Data/Arbevetel_idosoros_2017is.csv",
	};
	headers = []string{};
	yearexp, err := regexp.Compile("[^\\d]*(\\d{4,4}).*");
	h.Error(err, "", h.ERROR_LVL_WARNING);
	var opyds map[string]OperatorYearData = make(map[string]OperatorYearData);
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

			if (i == 1) {
				headers = line;
				continue;
			} else if (strings.Trim(line[1], " ") == "") {
				continue;
			}

			id, err := strconv.Atoi(line[0]);
			if (err != nil) {
				h.Error(err, "", h.ERROR_LVL_WARNING);
				continue;
			}

			var op Operator;
			op, err = op.Get(int64(id));

			var megjegyzes string = strings.Trim(line[len(line)-1]," ");
			if(megjegyzes != ""){
				//A cég megszűnt: 2015.07.01.
				termexp := regexp.MustCompile("megszűnt[^\\d]*(\\d{4}\\.\\d{2}\\.\\d{2}\\.)");
				subterm := termexp.FindStringSubmatch(line[len(line)-1]);
				if(len(subterm) == 2){
					termTime,err := time.Parse("2006.01.02.",subterm[1]);
					h.Error(err,"",h.ERROR_LVL_NOTICE);
					op.TerminationDate = termTime.Format(MYSQL_DATE_FORMAT);
					db.DbMap.Update(&op);
				}
			}

			if (err != nil) {
				h.PrintlnIf(fmt.Sprintf("Operator is missing with id %v", id), h.GetConfig().Mode.Debug);
				continue;
			}
			var exist bool;
			var m int = 2;
			var opyd OperatorYearData;
			var ok bool;

			for (m < len(line)) {
				exist = true;
				substrs := yearexp.FindStringSubmatch(headers[m]);
				if (len(substrs) > 1) {
					y, err := strconv.Atoi(substrs[1]);
					h.Error(err, "", h.ERROR_LVL_ERROR);
					opydykey := fmt.Sprintf("%v%v", y, id);
					opyd, ok = opyds[opydykey];
					if (!ok) {
						opyd, err = opyd.Get(int64(id), y);
						h.Error(err, "", h.ERROR_LVL_ERROR);
					}

					if (opyd.OperatorId == 0) {
						exist = false;
						opyd.OperatorId = int64(id);
						opyd.Year = y;
					}

					rowColVal, err := strconv.Atoi(line[m]);
					if (m%3 == 2) {
						opyd.IncomeNet = int64(rowColVal);
					} else if (m%3 == 0) {
						opyd.IncomeTax = int64(rowColVal);
					} else {
						opyd.IncomeOperational = int64(rowColVal);
					}

					opyds[opydykey] = opyd;
					if(exist){
						db.DbMap.Update(&opyd);
					} else {
						db.DbMap.Insert(&opyd);
					}
				}
				m++;
			}
		}
	}

	csvFiles = []string{
		"Data/Uzemelteto_allando_attributumai_alapitasiev_erdekeltseg.csv",
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
				if(i == 1){
					headers = line;
				}
				continue;
			}

			id, err := strconv.Atoi(line[0]);
			if (err != nil) {
				h.Error(err, "", h.ERROR_LVL_WARNING);
				continue;
			}

			var operator Operator;
			operator,err = operator.Get(int64(id));

			if(operator.Id == 0){
				fmt.Println("OPERATOR MISSING");
				continue
			}

			j:=3;
			for(j+1 < len(line)){
				var year int;
				var interest Interest;

				j++;

				year,err = strconv.Atoi(headers[j])
				if(err != nil){
					continue;
				}

				var interests []string = strings.Split(line[j], ",");

				k:=0;
				for(k<len(interests)){
					var interestTrim string = strings.Trim(interests[k], " ");
					if(interestTrim == ""){
						k++;
						continue;
					}
					interest,err = interest.SaveAndGetByName(interestTrim);
					if(interest.Id == 0){
						h.PrintlnIf(fmt.Sprintf("Missing interests: %s",interest.Name),h.GetConfig().Mode.Debug);
						k++;
						continue;
					}

					var oi OperatorInterest;
					oi,err = oi.GetByData(operator.Id,year,interest.Id)
					oi = OperatorInterest{
						0,
						operator.Id,
						year,
						interest.Id,
					}
					err = db.DbMap.Insert(&oi);

					h.Error(err,"",h.ERROR_LVL_ERROR);
					k++;
				}
			}

		}
	}

	h.PrintlnIf("Done importing operator year income data...",h.GetConfig().Mode.Debug);
}
