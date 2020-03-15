package model

import (
	"mos/db"
	"fmt"
	"errors"
	h "mos/helper"
	"mos/model/FElement"
	"github.com/valyala/fasthttp"
	"reflect"
	"strconv"
	"os"
	"encoding/csv"
	"bufio"
	"io"
	"strings"
	"database/sql"
	regexp2 "regexp"
)

type MediaOwner struct {
	Id      int64 `db:"id, primarykey, autoincrement"`
	MediaId int64 `db:"media_id"`
	OwnerId int64 `db:"owner_id"`
	Year    int   `db:"year,size:4"`
}

func (m MediaOwner) GetAll() []MediaOwner {
	var mediaowners []MediaOwner
	_, err := db.DbMap.Select(&mediaowners, fmt.Sprintf("SELECT * FROM %s", m.GetTable()));
	h.Error(err, "", h.ERROR_LVL_ERROR);
	return mediaowners;
}

func (m MediaOwner) Get(id int64) (MediaOwner, error) {
	var mediaowner MediaOwner

	if (id < 1) {
		return mediaowner, errors.New(fmt.Sprintf("Could not retrieve media owners to id %v", id));
	}

	err := db.DbMap.SelectOne(&mediaowner, fmt.Sprintf("SELECT * FROM %s WHERE id = ?", m.GetTable()), id);

	return mediaowner, err;
}

func GetMediaOwnerForm(newModelForm bool, data map[string]interface{}, action string, actionParams ...string) Form {
	var lineColMap = map[string]string{"lg": "4", "md": "4", "sm": "6", "xs": "12"};
	var fullColMap = map[string]string{"lg": "12", "md": "12", "sm": "12", "xs": "12"};
	var LeftFieldSet = Fieldset{"left", nil, lineColMap};
	var MiddleFieldSet = Fieldset{"middle", nil, lineColMap};
	var RightFieldSet = Fieldset{"right", nil, lineColMap};
	var FullFieldSet = Fieldset{"bottom", nil, fullColMap};
	var media Media;
	var owner Owner;
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
	var owners = owner.GetOptions(nil);
	var ownerIdInput = FElement.InputSelect{"Owner", "owner_id", "owner_id", "", false, false, []string{data["owner_id"].(string)}, false, owners, ""}
	RightFieldSet.AddElement(ownerIdInput);
	var buttonSubmit = FElement.InputButton{"Save", "save", "save", "", false, "", true, false, true, nil}
	FullFieldSet.AddElement(buttonSubmit);
	var form = Form{h.GetUrl(action, actionParams, true, "admin"), "POST", false, []Fieldset{LeftFieldSet, MiddleFieldSet, RightFieldSet, FullFieldSet}, false, nil, nil}

	return form;
}

func (mo MediaOwner) GetMedia() Media {
	var media Media;
	media, err := media.Get(mo.MediaId);
	h.Error(err, "", h.ERROR_LVL_WARNING);
	return media;
}

func (mo MediaOwner) GetOwner() Owner {
	var owner Owner;
	owner, err := owner.Get(mo.OwnerId);
	h.Error(err, "", h.ERROR_LVL_WARNING);
	return owner;
}

func (_ MediaOwner) IsLanguageModel() bool {
	return false;
}

func (_ MediaOwner) GetTable() string {
	return "media_owner";
}

func (_ MediaOwner) GetPrimaryKey() []string {
	return []string{"id"};
}

func NewMediaOwner(MediaId int64, OwnerId int64, Year int) MediaOwner {
	return MediaOwner{
		MediaId: MediaId,
		OwnerId: OwnerId,
		Year:    Year,
	};
}

func NewEmptyMediaOwner() MediaOwner {
	return NewMediaOwner(0, 0, 0);
}

func GetMediaOwnerFormValidator(ctx *fasthttp.RequestCtx, mediaOwner MediaOwner) Validator {
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
	Validator.AddField("owner_id", map[string]interface{}{
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

func (m MediaOwner) BuildStructure() {
	dbmap := db.DbMap;
	Conf := h.GetConfig();

	if (!Conf.Mode.Rebuild_structure) {
		return;
	}

	h.PrintlnIf(fmt.Sprintf("Drop %v table", m.GetTable()), Conf.Mode.Rebuild_structure);
	dbmap.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s;", m.GetTable()));
	h.PrintlnIf(fmt.Sprintf("Create %v table", m.GetTable()), Conf.Mode.Rebuild_structure);
	dbmap.CreateTablesIfNotExists();
	var indexes map[int]map[string]interface{} = make(map[int]map[string]interface{})

	indexes = map[int]map[string]interface{}{
		0: {
			"name":  "FK_MEDIA_OWNER_MEDIA_ID",
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
			"name":  "FK_MEDIA_OWNER_OWNER_ID",
			"type":  "fk",
			"field": []string{"owner_id"},
			"options": map[string]string{
				"field":           "owner_id",
				"reference_table": "owner",
				"reference_field": "id",
			},
			"unique": false,
		},
		2: {
			"name":   "IDX_MEDIA_OWNER_YEAR",
			"type":   "hash",
			"field":  []string{"year"},
			"unique": false,
		},
		3: {
			"name":   "UIDX_MEDIA_MEDIA_ID_OWNER_ID_YEAR",
			"type":   "hash",
			"field":  []string{"media_id", "year", "owner_id"},
			"unique": true,
		},
	};
	tablemap, err := dbmap.TableFor(reflect.TypeOf(MediaOwner{}), false);
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
			_, err := dbmap.Db.Query(indexQ);
			h.Error(err, "", h.ERROR_LVL_WARNING)
		} else {
			h.PrintlnIf(fmt.Sprintf("Create %s index", index["name"].(string)), Conf.Mode.Rebuild_structure);
			tablemap.AddIndex(index["name"].(string), index["type"].(string), index["field"].([]string)).SetUnique(index["unique"].(bool));
		}
	}

	dbmap.CreateIndex();
}

func (m MediaOwner) PrepeareData() {
	h.PrintlnIf("Start importing media owner data...", h.GetConfig().Mode.Debug);
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
		"Data/Halozati_tablak_vegtulajjal_2017.csv",
	};
	var missing map[int]interface{} = make(map[int]interface{});
	for _, fileName := range csvFiles {
		var year int = 0;
		h.PrintlnIf(fmt.Sprintf("Parsing file %s\r\n",fileName), h.GetConfig().Mode.Debug)
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
				yearExp := regexp2.MustCompile("^\\d{4}$");
				if(i == 1 && yearExp.MatchString(line[0])){
					year, err = strconv.Atoi(line[0]);
					continue;
				}
				yearExp = regexp2.MustCompile("^\\d{4}.*");
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
				//h.PrintlnIf(fmt.Sprintf("year is %v\r\n", year), h.GetConfig().Mode.Debug);
				continue;
			}

			var media Media;
			media,err = media.Get(int64(mediaId));
			if(media.Id == 0){
				missing[mediaId] = map[string]interface{}{
					"id":mediaId,
					"name":line[1],
				}
				continue;
			}

			var l = 4;
			for (l < len(line)) {
				if (strings.Trim(line[l], " ") == "") {
					break;
				}
				var mediaOw MediaOwner;
				ownerId, err := strconv.Atoi(line[l]);
				if(ownerId == 0){
					var owner Owner;
					owner,err = owner.GetByName(line[l]);
					h.Error(err,"",h.ERROR_LVL_ERROR);
					ownerId = int(owner.Id);
				}

				var owner Owner;
				owner,err = owner.Get(int64(ownerId));
				if(err == sql.ErrNoRows){
					owner.Id = int64(ownerId);
					owner.Name = line[l+1];
					h.PrintlnIf(fmt.Sprintf("Missiong owner, inserting %+v",owner),h.GetConfig().Mode.Debug);
					err = db.DbMap.Insert(&owner);
					h.Error(err,"",h.ERROR_LVL_ERROR);
				}

				db.DbMap.SelectOne(&mediaOw, fmt.Sprintf("SELECT * FROM %s WHERE media_id = ? AND owner_id = ? AND year = ?", mediaOw.GetTable()),mediaId,ownerId,year)

				if(mediaOw.Id==0) {
					mediaOw.MediaId = int64(mediaId);
					mediaOw.Year = year;
					mediaOw.OwnerId = int64(ownerId);
					err = db.DbMap.Insert(&mediaOw);
					h.Error(err, "", h.ERROR_LVL_ERROR);
				}
				l += 2;
			}
		}
	}
	/*for _,mm := range missing{
		fmt.Printf("id: %v, name: %s\r\n", mm.(map[string]interface{})["id"].(int),mm.(map[string]interface{})["name"].(string))
	}*/
	h.PrintlnIf(fmt.Sprintf("%v missing media found, skipped", len(missing)), h.GetConfig().Mode.Debug);
	h.PrintlnIf("Done importing media owner data...", h.GetConfig().Mode.Debug);
}
