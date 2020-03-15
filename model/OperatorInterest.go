package model

import (
	"mos/db"
	"fmt"
	h "mos/helper"
	"mos/model/FElement"
	"github.com/valyala/fasthttp"
	"reflect"
	"strconv"
	"errors"
)

type OperatorInterest struct {
	Id int64 `db:"id, primarykey, autoincrement"`
	OperatorId        int64   `db:"operator_id"`
	Year              int     `db:"year, size:4"`
	InterestId            int64 `db:"interest_id"`
}

func (o OperatorInterest) GetAll() []OperatorInterest {
	var operatorinterests []OperatorInterest
	_, err := db.DbMap.Select(&operatorinterests, fmt.Sprintf("SELECT * FROM %s ORDER BY operator_id, year", o.GetTable()));
	h.Error(err, "", h.ERROR_LVL_ERROR);
	return operatorinterests;
}

func (oi OperatorInterest) GetOperator() Operator{
	var operator Operator;
	operator,err := operator.Get(oi.OperatorId);
	h.Error(err,"",h.ERROR_LVL_WARNING);
	return operator;
}

func (oi OperatorInterest) GetInterest() Interest{
	var interest Interest;
	interest,err := interest.Get(oi.InterestId);
	h.Error(err,"",h.ERROR_LVL_WARNING);
	return interest;
}

func (_ OperatorInterest) IsLanguageModel() bool {
	return false;
}

func (_ OperatorInterest) GetTable() string {
	return "operator_interest";
}

func (_ OperatorInterest) GetPrimaryKey() []string {
	return []string{"id"};
}

func GetOperatorInterestForm(newModel bool, data map[string]interface{}, action string, actionParams ...string) Form {
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

func NewOperatorInterest(OperatorId int64, Year int, InterestId int64) OperatorInterest {
	return OperatorInterest{
		OperatorId: OperatorId,
		Year:       Year,
		InterestId:    InterestId,
	};
}

func GetOperatorInterestFormValidator(ctx *fasthttp.RequestCtx, OperatorInterest OperatorInterest) Validator {
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
	Validator.AddField("interest_id", map[string]interface{}{
		"roles": map[string]interface{}{
			"required": true,
		},
	});
	return Validator;
}

func (_ OperatorInterest) Get(id int64) (OperatorInterest, error) {
	var oi OperatorInterest;
	if (id == 0) {
		return oi, errors.New(fmt.Sprintf("Could not retrieve operator interest to id %v", id));
	}

	query := fmt.Sprintf("SELECT * FROM %v WHERE %s = ?", oi.GetTable(), "id");
	err := db.DbMap.SelectOne(&oi, query, id);
	h.Error(err, "", h.ERROR_LVL_ERROR);
	if (err != nil) {
		return oi, err;
	}

	if (oi.OperatorId == 0) {
		return oi, errors.New(fmt.Sprintf("Could not retrieve operator interest to id %v", id));
	}

	return oi, nil;
}

func (_ OperatorInterest) GetByData(operatorId int64, year int, interestId int64) (OperatorInterest, error) {
	var oi OperatorInterest;
	if (operatorId == 0 || year == 0 || interestId == 0) {
		return oi, errors.New(fmt.Sprintf("Could not retrieve operator interest to Operator Id %v, year %v and interest id %v", operatorId, year, interestId));
	}

	query := fmt.Sprintf("SELECT * FROM %v WHERE %s = ? and %s = ? and %s = ?", oi.GetTable(), "operator_id", "year", "interest_id");
	err := db.DbMap.SelectOne(&oi, query, operatorId, year, interestId);
	h.Error(err, "", h.ERROR_LVL_ERROR);
	if (err != nil) {
		return oi, err;
	}

	if (oi.OperatorId == 0) {
		return oi, errors.New(fmt.Sprintf("Could not retrieve operator interest to Operator Id %v, year %v and interest id %v", operatorId, year, interestId));
	}

	return oi, nil;
}

func (o OperatorInterest) BuildStructure() {
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
			"name":   "IDX_OPERATOR_INTEREST_OPERATOR_ID_YEAR_INTEREST_ID",
			"type":   "hash",
			"field":  []string{"operator_id", "year", "interest_id"},
			"unique": true,
		},
	};
	tablemap, err := dbmap.TableFor(reflect.TypeOf(OperatorInterest{}), false);
	h.Error(err, "", h.ERROR_LVL_ERROR);
	for _, index := range indexes {
		h.PrintlnIf(fmt.Sprintf("Create %s index", index["name"].(string)), Conf.Mode.Rebuild_structure);
		tablemap.AddIndex(index["name"].(string), index["type"].(string), index["field"].([]string)).SetUnique(index["unique"].(bool));
	}

	dbmap.CreateIndex();
}

func (o OperatorInterest) PrepeareData() {}
