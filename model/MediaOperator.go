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
	"regexp"
)

type MediaOperator struct {
	Id    int64 `db:"id, primarykey, autoincrement"`
	MediaId    int64 `db:"media_id"`
	OperatorId int64 `db:"operator_id"`
	Year       int   `db:"year,size:4"`
}

func (m MediaOperator) GetAll() []MediaOperator {
	var mediaoperators []MediaOperator
	_, err := db.DbMap.Select(&mediaoperators, fmt.Sprintf("select * from %s", m.GetTable()));
	h.Error(err, "", h.ERROR_LVL_ERROR);
	return mediaoperators;
}

func (_ MediaOperator) Get(id int64) (MediaOperator, error) {
	var mediaOperator MediaOperator;
	if ( id < 1) {
		return mediaOperator, errors.New(fmt.Sprintf("Could not retrieve media operator to id %v ", id));
	}

	err := db.DbMap.SelectOne(&mediaOperator, fmt.Sprintf("SELECT * FROM %v WHERE %s = ?", mediaOperator.GetTable(), "id"), id);
	h.Error(err, "", h.ERROR_LVL_ERROR);
	if (err != nil) {
		return mediaOperator, err;
	}

	if (mediaOperator.MediaId == 0) {
		return mediaOperator, errors.New(fmt.Sprintf("Could not retrieve media operator to id %v", id))
	}

	return mediaOperator, nil;
}

func (_ MediaOperator) IsLanguageModel() bool {
	return false;
}

func (_ MediaOperator) GetTable() string {
	return "media_operator";
}

func (_ MediaOperator) GetPrimaryKey() []string {
	return []string{"id"};
}

func (mo MediaOperator) GetMedia() Media{
	var media Media;
	media,err := media.Get(mo.MediaId);
	h.Error(err,"",h.ERROR_LVL_WARNING);
	return media;
}

func (mo MediaOperator) GetOperator() Operator{
	var operator Operator;
	operator,err := operator.Get(mo.OperatorId);
	h.Error(err,"",h.ERROR_LVL_WARNING);
	return operator;
}

func GetMediaOperatorForm(newModelForm bool, data map[string]interface{}, action string, actionParams ...string) Form {
	var lineColMap = map[string]string{"lg": "4", "md": "4", "sm": "6", "xs": "12"};
	var fullColMap = map[string]string{"lg": "12", "md": "12", "sm": "12", "xs": "12"};

	var LeftFieldSet = Fieldset{"left", nil, lineColMap};
	var MiddleFieldSet = Fieldset{"middle", nil, lineColMap};
	var RightFieldSet = Fieldset{"right", nil, lineColMap};
	var FullFieldSet = Fieldset{"bottom", nil, fullColMap};

	var media Media;
	var operator Operator;
	var mediaIdInput FormElement;

	idInput := FElement.InputHidden{"id", "id", "", false, true, data["id"].(string)}
	LeftFieldSet.AddElement(idInput);

	if (!newModelForm) {
		intMediaId, err := strconv.Atoi(data["media_id"].(string));
		h.Error(err, "", h.ERROR_LVL_ERROR);
		media, err := media.Get(int64(intMediaId));
		h.Error(err, "", h.ERROR_LVL_ERROR);
		mediaIdInput = FElement.InputHidden{"media_id", "media_id", "", false, true, data["media_id"].(string)}
		var mediaIdTextInput = FElement.InputText{"Media", "", "", "", "", false, true, media.Name, "", "", "", "", ""}
		LeftFieldSet.AddElement(mediaIdTextInput);
	} else {
		mediaIdInput = FElement.InputSelect{"Media", "media_id", "media_id", "", false, false, []string{data["media_id"].(string)}, false, media.GetOptions(nil), ""}
	}
	LeftFieldSet.AddElement(mediaIdInput);
	var yearInput = FElement.InputText{"Year", "year", "year", "", "", false, !newModelForm, data["year"].(string), "", "", "", "", ""}
	MiddleFieldSet.AddElement(yearInput);
	var operators = operator.GetOptions(nil);
	var ownerIdInput = FElement.InputSelect{"Operator", "operator_id", "operator_id", "", false, false, []string{data["operator_id"].(string)}, false, operators, ""}
	RightFieldSet.AddElement(ownerIdInput);
	var buttonSubmit = FElement.InputButton{"Save", "save", "save", "", false, "", true, false, true, nil}
	FullFieldSet.AddElement(buttonSubmit);
	var form = Form{h.GetUrl(action, actionParams, true, "admin"), "POST", false, []Fieldset{LeftFieldSet, MiddleFieldSet, RightFieldSet, FullFieldSet}, false, nil, nil}

	return form;
}

func NewMediaOperator(MediaId int64, OperatorId int64, Year int) MediaOperator {
	return MediaOperator{
		MediaId:    MediaId,
		OperatorId: OperatorId,
		Year:       Year,
	};
}

func NewEmptyMediaOperator() MediaOperator {
	return NewMediaOperator(0, 0, 0);
}

func GetMediaOperatorFormValidator(ctx *fasthttp.RequestCtx, MediaOperator MediaOperator) Validator {
	var Validator Validator;
	Validator = Validator.New(ctx);
	Validator.AddField("media_id", map[string]interface{}{
		"roles": map[string]interface{}{
			"required": true,
			"format": map[string]interface{}{
				"type":   "regexp",
				"regexp": "^\\d+$",
			},
		},
	});
	Validator.AddField("operator_id", map[string]interface{}{
		"roles": map[string]interface{}{
			"required": true,
			"format": map[string]interface{}{
				"type":   "regexp",
				"regexp": "^\\d+$",
			},
		},
	});
	Validator.AddField("year", map[string]interface{}{
		"roles": map[string]interface{}{
			"required": true,
			"format": map[string]interface{}{
				"type":   "regexp",
				"regexp": "^\\d{4}$",
			},
		},
	});
	return Validator;
}

func (m MediaOperator) BuildStructure() {
	dbmap := db.DbMap;
	Conf := h.GetConfig();

	if(!Conf.Mode.Rebuild_structure){
		return;
	}

	h.PrintlnIf(fmt.Sprintf("Drop %v table", m.GetTable()), Conf.Mode.Rebuild_structure);
	dbmap.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s;", m.GetTable()));

	h.PrintlnIf(fmt.Sprintf("Create %v table", m.GetTable()), Conf.Mode.Rebuild_structure);
	dbmap.CreateTablesIfNotExists();
	var indexes map[int]map[string]interface{} = make(map[int]map[string]interface{})

	indexes = map[int]map[string]interface{}{
		0: {
			"name":  "FK_MEDIA_OPERATOR_MEDIA_ID",
			"type":  "fk",
			"field": []string{"media_id"},
			"options": map[string]string{
				"field":           "media_id",
				"reference_table": "media",
				"reference_field": "id",
			},
			"unique": false,
		},
		1: {
			"name":  "FK_MEDIA_OPERATOR_OPERATOR_ID",
			"type":  "fk",
			"field": []string{"operator_id"},
			"options": map[string]string{
				"field":           "operator_id",
				"reference_table": "operator",
				"reference_field": "id",
			},
			"unique": false,
		},
		2: {
			"name":  "IDX_MEDIA_MEDIA_ID_YEAR",
			"type":  "hash",
			"field": []string{"media_id","year"},
			"unique": true,
		},
	};
	tablemap, err := dbmap.TableFor(reflect.TypeOf(MediaOperator{}), false);
	h.Error(err, "", h.ERROR_LVL_ERROR);
	for _, index := range indexes {
		if (index["type"] == "fk") {
			options := index["options"].(map[string]string);
			var indexQ string = fmt.Sprintf(
				"ALTER TABLE %v ADD CONSTRAINT %v FOREIGN KEY (%v) REFERENCES %v (%v);",
				m.GetTable(),
				index["name"].(string),
				options["field"],
				options["reference_table"],
				options["reference_field"],
			);
			h.PrintlnIf(fmt.Sprintf("Create foreign key %v ", index["name"].(string)), Conf.Mode.Rebuild_structure);
			h.PrintlnIf(indexQ, Conf.Mode.Rebuild_structure);
			_, err := dbmap.Db.Exec(indexQ);
			h.Error(err, "", h.ERROR_LVL_WARNING)
		} else {
			h.PrintlnIf(fmt.Sprintf("Create %s index", index["name"].(string)), Conf.Mode.Rebuild_structure);
			tablemap.AddIndex(index["name"].(string), index["type"].(string), index["field"].([]string)).SetUnique(index["unique"].(bool));
		}
	}

	dbmap.CreateIndex();
}

func (m MediaOperator) PrepeareData() {
	h.PrintlnIf("Start importing media operator data...",h.GetConfig().Mode.Debug);
	var mediaSkip,operatorSkip int;
	var csvFiles []string = []string{
		"Data/Halozati_tablak_kesz_1998.csv",
		"Data/Halozati_tablak_kesz_1999.csv",
		"Data/Halozati_tablak_kesz_2000.csv",
		"Data/Halozati_tablak_kesz_2001.csv",
		"Data/Halozati_tablak_kesz_2002.csv",
		"Data/Halozati_tablak_kesz_2003.csv",
		"Data/Halozati_tablak_kesz_2004.csv",
		"Data/Halozati_tablak_kesz_2005.csv",
		"Data/Halozati_tablak_kesz_2006.csv",
		"Data/Halozati_tablak_kesz_2007.csv",
		"Data/Halozati_tablak_kesz_2008.csv",
		"Data/Halozati_tablak_kesz_2009.csv",
		"Data/Halozati_tablak_evente_vegtulaj_ID_3_2010.csv",
		"Data/Halozati_tablak_evente_vegtulaj_ID_3_2011.csv",
		"Data/Halozati_tablak_evente_vegtulaj_ID_3_2012.csv",
		"Data/Halozati_tablak_evente_vegtulaj_ID_3_2013.csv",
		"Data/Halozati_tablak_evente_vegtulaj_ID_3_2014.csv",
		"Data/Halozati_tablak_evente_vegtulaj_ID_3_2015.csv",
		"Data/Halozati_tablak_evente_vegtulaj_ID_3_2016.csv",
	};
	for _,fileName := range csvFiles{
		h.PrintlnIf(fmt.Sprintf("Parsing file %s\r\n",fileName), h.GetConfig().Mode.Debug)
		var year int;
		csvFile, _ := os.Open(fileName)
		reader := csv.NewReader(bufio.NewReader(csvFile))
		i := 0;
		for {
			i++;
			//for escaping if end of file or error reading file
			line, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				h.Error(err, "", h.ERROR_LVL_ERROR)
			}

			//getting/setting year
			if(year == 0){
				yearExp := regexp.MustCompile("^\\d{4}$");
				if(i == 1 && yearExp.MatchString(line[0])){
					year, err = strconv.Atoi(line[0]);
					continue;
				}
				yearExp = regexp.MustCompile("^\\d{4}.*");
				if(i == 1 && yearExp.MatchString(line[3])){
					year, err = strconv.Atoi(line[3][0:4]);
					continue;
				}
			}

			if(strings.Trim(line[0], " ") == ""){
				continue;
			}

			mediaId, err := strconv.Atoi(line[0]);
			if (err != nil) {
				//h.PrintlnIf(fmt.Sprintf(line[0] + " is not integer"), h.GetConfig().Mode.Debug);
				//h.PrintlnIf(fmt.Sprintf("year is %v\r\n", year), h.GetConfig().Mode.Debug);
				continue;
			}

			operatorId, err := strconv.Atoi(line[2]);
			if (err != nil) {
				h.Error(err, "", h.ERROR_LVL_WARNING);
				continue;
			}

			var mediaOp MediaOperator;
			var media Media;
			var operator Operator;

			media,err = media.Get(int64(mediaId))
			if(err != nil){
				mediaSkip++;
				continue;
			}
			operator,err = operator.Get(int64(operatorId))
			if(err != nil){
				operatorSkip++;
				continue;
			}

			mediaOp.MediaId = int64(mediaId);
			mediaOp.OperatorId = int64(operatorId);
			mediaOp.Year = year;

			err = db.DbMap.Insert(&mediaOp);
			h.Error(err, "", h.ERROR_LVL_ERROR);

			/*var l = 3;
			for(l<len(headers)){
				y,err := strconv.Atoi(headers[l]);
				h.Error(err,"",h.ERROR_LVL_WARNING);
				operatoryd := OperatorYearData{
					int64(id),
					y,
					line[l],
				}
				err = db.DbMap.Insert(&operatoryd);
				h.Error(err, "", h.ERROR_LVL_ERROR);
				l++;
			}*/
		}
	}

	h.PrintlnIf(fmt.Sprintf("Missing media and skip row: %v\r\n",mediaSkip), h.GetConfig().Mode.Debug)
	h.PrintlnIf(fmt.Sprintf("Missing operator and skip row: %v\r\n",operatorSkip), h.GetConfig().Mode.Debug)
	h.PrintlnIf("Done importing media operator data...",h.GetConfig().Mode.Debug);
}
