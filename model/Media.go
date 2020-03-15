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

type Media struct {
	Id          int64  `db:"id, primarykey"`
	Name        string `db:"name, size:255"`
	MediaTypeId int64  `db:"media_type_id"`
	News        bool   `db:"news"`
}

func (m Media) GetAll() []Media {
	var medias []Media
	_, err := db.DbMap.Select(&medias, fmt.Sprintf("select * from %s order by %v", m.GetTable(), m.GetPrimaryKey()[0]));
	h.Error(err, "", h.ERROR_LVL_ERROR);
	return medias;
}

func (m Media) GetOptions(defOption map[string]string) []map[string]string {
	var medias = m.GetAll();
	var options []map[string]string;
	if (defOption != nil) {
		_, okl := defOption["label"];
		_, okv := defOption["value"];
		if (okl || okv) {
			options = append(options, defOption);
		}
	}
	for _, media := range medias {
		options = append(options, map[string]string{"label": media.Name, "value": strconv.Itoa(int(media.Id))});
	}
	return options;
}

func (m Media) GetOwners() []MediaOwner{
	var results []MediaOwner;
	var mediaOwner MediaOwner;
	_,err := db.DbMap.Select(&results, fmt.Sprintf("SELECT * FROM %s WHERE media_id = ?", mediaOwner.GetTable()),m.Id);
	h.Error(err,"",h.ERROR_LVL_WARNING);

	return results;
}

func (_ Media) Get(id int64) (Media, error) {
	var media Media;
	if (id == 0) {
		return media, errors.New(fmt.Sprintf("Could not retrieve media to ID %v", id));
	}

	err := db.DbMap.SelectOne(&media, fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", media.GetTable(), media.GetPrimaryKey()[0]), id);
	h.Error(err, "", h.ERROR_LVL_ERROR);
	if (err != nil) {
		return media, err;
	}

	if (media.Id == 0) {
		return media, errors.New(fmt.Sprintf("Could not retrieve media to ID %v", id))
	}

	return media, nil;
}

func (_ Media) IsLanguageModel() bool {
	return false;
}

func (_ Media) GetTable() string {
	return "media";
}

func (_ Media) GetPrimaryKey() []string {
	return []string{"id"};
}

func GetMediaForm(data map[string]interface{}, action string) Form {
	var Elements []FormElement;
	var id = FElement.InputHidden{"id", "id", "", false, true, data["id"].(string)}
	Elements = append(Elements, id);
	var identifier = FElement.InputText{"Name", "name", "name", "", "", false, false, data["name"].(string), "", "", "", "", ""}
	Elements = append(Elements, identifier);
	var mediaType MediaType;
	var options = mediaType.GetOptions(nil);
	var mediaTypeInp = FElement.InputSelect{"Media Type", "media_type_id", "media_type_id", "", false, false, []string{data["media_type_id"].(string)}, false, options, ""}
	Elements = append(Elements, mediaTypeInp);
	Checkboxes := FElement.CheckboxGroup{};
	checkbox := FElement.InputCheckbox{
		"News",
		"news",
		"news",
		"",
		false,
		false,
		"1",
		[]string{data["news"].(string)},
		true,
	}
	Checkboxes.Checkbox = append(Checkboxes.Checkbox, checkbox);
	Elements = append(Elements, Checkboxes);

	if(data["id"].(string) != "") {
		var ownersLink= FElement.Static{"", "", "", "", fmt.Sprintf(`<a href="%v">Show Media Owner History</a>`, h.GetUrl("media/owners", []string{data["id"].(string)}, true, "admin"))}
		Elements = append(Elements, ownersLink);
	}

	var fullColMap = map[string]string{"lg": "12", "md": "12", "sm": "12", "xs": "12"};
	var Fieldsets []Fieldset;
	Fieldsets = append(Fieldsets, Fieldset{"left", Elements, fullColMap});
	button := FElement.InputButton{"Submit", "submit", "submit", "pull-right", false, "", true, false, false, nil}
	Fieldsets = append(Fieldsets, Fieldset{"bottom", []FormElement{button}, fullColMap});
	var form = Form{h.GetUrl(action, nil, true, "admin"), "POST", false, Fieldsets, false, nil, nil}

	return form;
}

func NewMedia(Id int64, Name string, MediaTypeId int64, News bool) Media {
	return Media{
		Id:          Id,
		Name:        Name,
		MediaTypeId: MediaTypeId,
		News:        News,
	};
}

func NewEmptyMedia() Media {
	return NewMedia(0, "", 0, false);
}

func GetMediaFormValidator(ctx *fasthttp.RequestCtx, Media Media) Validator {
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

func (m Media) BuildStructure() {
	dbmap := db.DbMap;
	Conf := h.GetConfig();

	if(!Conf.Mode.Rebuild_structure){
		return;
	}

	h.PrintlnIf(fmt.Sprintf("Drop %v table", m.GetTable()), Conf.Mode.Rebuild_structure);
	_,err := dbmap.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s;", m.GetTable()));
	h.Error(err,"",h.ERROR_LVL_ERROR)

	h.PrintlnIf(fmt.Sprintf("Create %v table", m.GetTable()), Conf.Mode.Rebuild_structure);
	dbmap.CreateTablesIfNotExists();
	var indexes map[int]map[string]interface{} = make(map[int]map[string]interface{})

	indexes = map[int]map[string]interface{}{
		0: {
			"name":   "FK_MEDIA_MEDIA_TYPE_ID_MEDIA_TYPE_ID",
			"type":   "fk",
			"field":  []string{"media_type_id"},
			"options": map[string]string{
				"field": "media_type_id",
				"reference_table": "media_type",
				"reference_field": "id",
			},
			"unique": false,
		},
		1: {
			"name":   "IDX_MEDIA_NEWS",
			"type":   "hash",
			"field":  []string{"news"},
			"unique": false,
		},
	};
	tablemap, err := dbmap.TableFor(reflect.TypeOf(Media{}), false);
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
			h.PrintlnIf(indexQ, Conf.Mode.Debug);
			_, err := dbmap.Db.Exec(indexQ);
			h.Error(err, "", h.ERROR_LVL_WARNING)
		} else {
			h.PrintlnIf(fmt.Sprintf("Create %s index", index["name"].(string)), Conf.Mode.Debug);
			tablemap.AddIndex(index["name"].(string), index["type"].(string), index["field"].([]string)).SetUnique(index["unique"].(bool));
		}
	}

	dbmap.CreateIndex();
}

func (m Media) PrepeareData(){
	h.PrintlnIf("Start importing media data...",h.GetConfig().Mode.Debug);
	var csvFiles []string = []string{
		"Data/media_allando_attributumai_1998_2009.csv",
		"Data/Medium_allando_attributumai_hirmedium_2017.csv",
	};
	for _, fileName := range csvFiles {
		h.PrintlnIf(fmt.Sprintf("Parsing file %s\r\n",fileName), h.GetConfig().Mode.Debug)
		csvFile, _ := os.Open(fileName)
		reader := csv.NewReader(bufio.NewReader(csvFile))
		i := 0;
		var mediaType MediaType;
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

			mt, err := mediaType.SaveAndGetByName(strings.Trim(line[2], " "));
			h.Error(err, "", h.ERROR_LVL_ERROR);

			var media Media;
			media,err = media.Get(int64(id))

			if(media.Id == 0) {
				media := Media{
					int64(id),
					line[1],
					mt.Id,
					line[3] == "1",
				}

				err = db.DbMap.Insert(&media);
				h.Error(err, "", h.ERROR_LVL_ERROR);
			}
		}
	}
	h.PrintlnIf("Done importing media data...",h.GetConfig().Mode.Debug);
}