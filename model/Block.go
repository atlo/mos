package model

import (
	"mos/db"
	"fmt"
	"errors"
	h "mos/helper"
	"mos/model/FElement"
	"github.com/valyala/fasthttp"
	"reflect"
)

type Block struct {
	Id         int64  `db:"id, primarykey, autoincrement"`
	Identifier string `db:"identifier, size:255"`
	Title string `db:"title, size:500"`
	Content    string `db:"content, size:1000"`
	Lc         string `db:"lc, size:2"`
}

func (b Block) GetAll() []Block {
	var blocks []Block
	_, err := db.DbMap.Select(&blocks, fmt.Sprintf("SELECT * FROM %s order by %v", b.GetTable(), b.GetPrimaryKey()[0]));
	h.Error(err, "", h.ERROR_LVL_ERROR);
	return blocks;
}

func (_ Block) Get(blockId int64) (Block, error) {
	var block Block;
	if (blockId == 0) {
		return block, errors.New(fmt.Sprintf("Could not retrieve block to ID %v", blockId));
	}

	err := db.DbMap.SelectOne(&block, fmt.Sprintf("SELECT * FROM %s WHERE %v = ?", block.GetTable(), block.GetPrimaryKey()[0]), blockId);
	h.Error(err, "", h.ERROR_LVL_ERROR);
	if (err != nil) {
		return block, err;
	}

	if (block.Id == 0) {
		return block, errors.New(fmt.Sprintf("Could not retrieve block to ID %v", blockId))
	}

	return block, nil;
}

func (_ Block) IsLanguageModel() bool {
	return true;
}

func (_ Block) GetTable() string {
	return "block";
}

func (_ Block) GetPrimaryKey() []string {
	return []string{"id"};
}

func GetBlockForm(data map[string]interface{}, action string) Form {
	var Elements []FormElement;
	var id = FElement.InputHidden{"id", "id", "", false, true, data["id"].(string)}
	Elements = append(Elements, id);
	var lc = FElement.InputHidden{"lc", "lc", "", false, true, data["lc"].(string)}
	Elements = append(Elements, lc);
	var identifier = FElement.InputText{"Identifier", "identifier", "", "", "fe.: iden-ti-fier", false, false, data["identifier"].(string), "Unique per language (this will be used to load the block)", "", "", "", ""}
	Elements = append(Elements, identifier);
	var title = FElement.InputText{"Title", "title", "", "", "", false, false, data["title"].(string), "", "", "", "", ""}
	Elements = append(Elements, title);
	var content = FElement.InputTextarea{"Content", "content", "content", "", "Content to display", false, false, data["content"].(string), "", 80, 5}
	Elements = append(Elements, content);
	var fullColMap = map[string]string{"lg": "12", "md": "12", "sm": "12", "xs": "12"};
	var Fieldsets []Fieldset;
	Fieldsets = append(Fieldsets, Fieldset{"left", Elements, fullColMap});
	button := FElement.InputButton{"Submit", "submit", "submit", "pull-right", false, "", true, false, false, nil}
	Fieldsets = append(Fieldsets, Fieldset{"bottom", []FormElement{button}, fullColMap});
	var form = Form{h.GetUrl(action, nil, true, "admin"), "POST", false, Fieldsets, false, nil, nil}

	return form;
}

func NewBlock(Id int64, Identifier string, Content string, Lang string) Block {
	return Block{
		Id:         Id,
		Identifier: Identifier,
		Content:    Content,
		Lc:         Lang,
	};
}

func NewEmptyBlock() Block {
	return NewBlock(0, "", "", h.DefLang)
}

func GetBlockFormValidator(ctx *fasthttp.RequestCtx, Block Block) Validator {
	var Validator Validator;
	Validator = Validator.New(ctx);
	Validator.AddField("id", map[string]interface{}{
		"roles": map[string]interface{}{
			"required": false,
		},
	});
	Validator.AddField("title", map[string]interface{}{
		"roles": map[string]interface{}{
			"required": true,
		},
	});
	Validator.AddField("identifier", map[string]interface{}{
		"roles": map[string]interface{}{
			"required": true,
			"format": map[string]interface{}{
				"type":   "regexp",
				"regexp": "^([a-zA-Z0-9\\-\\_]*)+$",
			},
		},
	});
	Validator.AddField("content", map[string]interface{}{
		"roles": map[string]interface{}{
			"required": true,
		},
	});
	return Validator;
}

func (b Block) GetByIdentifier(identifier string, languageCode string) (Block, error) {
	var block Block
	var err error;

	err = db.DbMap.SelectOne(&block, fmt.Sprintf("SELECT * FROM %s WHERE %v= ? AND %v = ?", b.GetTable(), "lc", "identifier"), languageCode, identifier);

	return block, err;
}

func (b Block) BuildStructure() {
	dbmap := db.DbMap;
	Conf := h.GetConfig();

	if(!Conf.Mode.Rebuild_structure){
		return;
	}

	h.PrintlnIf(fmt.Sprintf("Drop %v table", b.GetTable()), Conf.Mode.Rebuild_structure);
	dbmap.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s;", b.GetTable()));

	h.PrintlnIf(fmt.Sprintf("Create %v table", b.GetTable()), Conf.Mode.Rebuild_structure);
	dbmap.CreateTablesIfNotExists();
	var indexes map[int]map[string]interface{} = make(map[int]map[string]interface{})

	indexes = map[int]map[string]interface{}{
		0: {
			"name":   "IDX_BLOCK_IDENTIFIER_LC",
			"type":   "hash",
			"field":  []string{"identifier", "lc"},
			"unique": true,
		},
	};
	tablemap, err := dbmap.TableFor(reflect.TypeOf(Block{}), false);
	h.Error(err, "", h.ERROR_LVL_ERROR);
	for _, index := range indexes {
		h.PrintlnIf(fmt.Sprintf("Create %s index", index["name"].(string)), Conf.Mode.Rebuild_structure);
		tablemap.AddIndex(index["name"].(string), index["type"].(string), index["field"].([]string)).SetUnique(index["unique"].(bool));
	}

	dbmap.CreateIndex();
	var blockCont map[string]map[string]string = map[string]map[string]string{
		"introduction": {
			"en": "Introduction English - Suspendisse a aliquet massa. Etiam sed ante in diam molestie sollicitudin. Vivamus vulputate lacus diam, nec auctor urna lacinia ac. Quisque elementum tempor scelerisque. Donec id dui lacus. Ut in mauris fermentum, varius purus ac, consequat felis. ",
			"de": "Introduction Deutsch - Suspendisse a aliquet massa. Etiam sed ante in diam molestie sollicitudin. Vivamus vulputate lacus diam, nec auctor urna lacinia ac. Quisque elementum tempor scelerisque. Donec id dui lacus. Ut in mauris fermentum, varius purus ac, consequat felis. ",
			"hu": "Introduction Hungarian - Suspendisse a aliquet massa. Etiam sed ante in diam molestie sollicitudin. Vivamus vulputate lacus diam, nec auctor urna lacinia ac. Quisque elementum tempor scelerisque. Donec id dui lacus. Ut in mauris fermentum, varius purus ac, consequat felis. ",
		},
		"about": {
			"en": "About English - Duis sit amet tincidunt purus. Cras finibus mi ut purus pretium porttitor. Nunc sed varius ipsum. Cras a pharetra nibh. Nam sagittis urna sapien, at feugiat lacus pretium eu. Nunc nec nulla ut elit euismod auctor. Duis id consectetur lacus. Aliquam erat volutpat. Aliquam erat volutpat. Pellentesque vulputate, tortor vel blandit fringilla, sem neque blandit lorem, sed vulputate nisl magna in tortor. Sed congue sapien quis est semper, a vehicula mi pulvinar. Curabitur pulvinar sapien id ligula sagittis suscipit. Maecenas dictum, neque ultrices aliquam accumsan, ante nunc vulputate magna, ac ultrices sapien nulla sit amet mi. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Sed sollicitudin lacus diam, sit amet rutrum libero facilisis sed. Curabitur feugiat eleifend nisl ut vulputate. Etiam eleifend nibh eget massa aliquam, at volutpat turpis porttitor. Aliquam diam sapien, vestibulum sed tincidunt quis, interdum vitae quam. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Ut id sodales nunc, vitae porttitor nibh. Vestibulum convallis, neque blandit fringilla aliquam, turpis ligula vehicula leo, eu imperdiet dui dolor sed mi. Maecenas vitae feugiat libero. Quisque ac rutrum augue. Phasellus metus nibh, hendrerit sit amet sapien at, ullamcorper sodales lacus. Vivamus nisi diam, hendrerit nec ligula sed, viverra sagittis massa. Curabitur at libero dolor. Nullam dolor ipsum, tincidunt et libero eu, tristique sodales orci. Praesent lorem ex, sodales non aliquam eu, accumsan at arcu. Proin laoreet eget ex nec efficitur. In sollicitudin tellus enim, et placerat dui scelerisque vitae. Etiam non turpis vel massa ullamcorper aliquam nec quis felis. Suspendisse mollis turpis felis, eget dapibus ante sodales suscipit. Etiam eget diam et mi efficitur lobortis. Proin eu felis ac urna tempus ornare. Donec convallis ligula vel libero porta imperdiet. Donec nibh nunc, pretium et ultrices at, porta vitae lorem. Nunc gravida a felis eget tempor.",
			"de": "About Deutsch - Duis sit amet tincidunt purus. Cras finibus mi ut purus pretium porttitor. Nunc sed varius ipsum. Cras a pharetra nibh. Nam sagittis urna sapien, at feugiat lacus pretium eu. Nunc nec nulla ut elit euismod auctor. Duis id consectetur lacus. Aliquam erat volutpat. Aliquam erat volutpat. Pellentesque vulputate, tortor vel blandit fringilla, sem neque blandit lorem, sed vulputate nisl magna in tortor. Sed congue sapien quis est semper, a vehicula mi pulvinar. Curabitur pulvinar sapien id ligula sagittis suscipit. Maecenas dictum, neque ultrices aliquam accumsan, ante nunc vulputate magna, ac ultrices sapien nulla sit amet mi. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Sed sollicitudin lacus diam, sit amet rutrum libero facilisis sed. Curabitur feugiat eleifend nisl ut vulputate. Etiam eleifend nibh eget massa aliquam, at volutpat turpis porttitor. Aliquam diam sapien, vestibulum sed tincidunt quis, interdum vitae quam. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Ut id sodales nunc, vitae porttitor nibh. Vestibulum convallis, neque blandit fringilla aliquam, turpis ligula vehicula leo, eu imperdiet dui dolor sed mi. Maecenas vitae feugiat libero. Quisque ac rutrum augue. Phasellus metus nibh, hendrerit sit amet sapien at, ullamcorper sodales lacus. Vivamus nisi diam, hendrerit nec ligula sed, viverra sagittis massa. Curabitur at libero dolor. Nullam dolor ipsum, tincidunt et libero eu, tristique sodales orci. Praesent lorem ex, sodales non aliquam eu, accumsan at arcu. Proin laoreet eget ex nec efficitur. In sollicitudin tellus enim, et placerat dui scelerisque vitae. Etiam non turpis vel massa ullamcorper aliquam nec quis felis. Suspendisse mollis turpis felis, eget dapibus ante sodales suscipit. Etiam eget diam et mi efficitur lobortis. Proin eu felis ac urna tempus ornare. Donec convallis ligula vel libero porta imperdiet. Donec nibh nunc, pretium et ultrices at, porta vitae lorem. Nunc gravida a felis eget tempor.",
			"hu": "About Hungarian - Duis sit amet tincidunt purus. Cras finibus mi ut purus pretium porttitor. Nunc sed varius ipsum. Cras a pharetra nibh. Nam sagittis urna sapien, at feugiat lacus pretium eu. Nunc nec nulla ut elit euismod auctor. Duis id consectetur lacus. Aliquam erat volutpat. Aliquam erat volutpat. Pellentesque vulputate, tortor vel blandit fringilla, sem neque blandit lorem, sed vulputate nisl magna in tortor. Sed congue sapien quis est semper, a vehicula mi pulvinar. Curabitur pulvinar sapien id ligula sagittis suscipit. Maecenas dictum, neque ultrices aliquam accumsan, ante nunc vulputate magna, ac ultrices sapien nulla sit amet mi. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Sed sollicitudin lacus diam, sit amet rutrum libero facilisis sed. Curabitur feugiat eleifend nisl ut vulputate. Etiam eleifend nibh eget massa aliquam, at volutpat turpis porttitor. Aliquam diam sapien, vestibulum sed tincidunt quis, interdum vitae quam. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Ut id sodales nunc, vitae porttitor nibh. Vestibulum convallis, neque blandit fringilla aliquam, turpis ligula vehicula leo, eu imperdiet dui dolor sed mi. Maecenas vitae feugiat libero. Quisque ac rutrum augue. Phasellus metus nibh, hendrerit sit amet sapien at, ullamcorper sodales lacus. Vivamus nisi diam, hendrerit nec ligula sed, viverra sagittis massa. Curabitur at libero dolor. Nullam dolor ipsum, tincidunt et libero eu, tristique sodales orci. Praesent lorem ex, sodales non aliquam eu, accumsan at arcu. Proin laoreet eget ex nec efficitur. In sollicitudin tellus enim, et placerat dui scelerisque vitae. Etiam non turpis vel massa ullamcorper aliquam nec quis felis. Suspendisse mollis turpis felis, eget dapibus ante sodales suscipit. Etiam eget diam et mi efficitur lobortis. Proin eu felis ac urna tempus ornare. Donec convallis ligula vel libero porta imperdiet. Donec nibh nunc, pretium et ultrices at, porta vitae lorem. Nunc gravida a felis eget tempor.",
		},
		"impressum": {
			"en": "Impressum English - Duis sit amet tincidunt purus. Cras finibus mi ut purus pretium porttitor. Nunc sed varius ipsum. Cras a pharetra nibh. Nam sagittis urna sapien, at feugiat lacus pretium eu. Nunc nec nulla ut elit euismod auctor. Duis id consectetur lacus. Aliquam erat volutpat. Aliquam erat volutpat. Pellentesque vulputate, tortor vel blandit fringilla, sem neque blandit lorem, sed vulputate nisl magna in tortor. Sed congue sapien quis est semper, a vehicula mi pulvinar. Curabitur pulvinar sapien id ligula sagittis suscipit. Maecenas dictum, neque ultrices aliquam accumsan, ante nunc vulputate magna, ac ultrices sapien nulla sit amet mi. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Sed sollicitudin lacus diam, sit amet rutrum libero facilisis sed. Curabitur feugiat eleifend nisl ut vulputate. Etiam eleifend nibh eget massa aliquam, at volutpat turpis porttitor. Aliquam diam sapien, vestibulum sed tincidunt quis, interdum vitae quam. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Ut id sodales nunc, vitae porttitor nibh. Vestibulum convallis, neque blandit fringilla aliquam, turpis ligula vehicula leo, eu imperdiet dui dolor sed mi. Maecenas vitae feugiat libero. Quisque ac rutrum augue. Phasellus metus nibh, hendrerit sit amet sapien at, ullamcorper sodales lacus. Vivamus nisi diam, hendrerit nec ligula sed, viverra sagittis massa. Curabitur at libero dolor. Nullam dolor ipsum, tincidunt et libero eu, tristique sodales orci. Praesent lorem ex, sodales non aliquam eu, accumsan at arcu. Proin laoreet eget ex nec efficitur. In sollicitudin tellus enim, et placerat dui scelerisque vitae. Etiam non turpis vel massa ullamcorper aliquam nec quis felis. Suspendisse mollis turpis felis, eget dapibus ante sodales suscipit. Etiam eget diam et mi efficitur lobortis. Proin eu felis ac urna tempus ornare. Donec convallis ligula vel libero porta imperdiet. Donec nibh nunc, pretium et ultrices at, porta vitae lorem. Nunc gravida a felis eget tempor.",
			"de": "Impressum Deutsch - Duis sit amet tincidunt purus. Cras finibus mi ut purus pretium porttitor. Nunc sed varius ipsum. Cras a pharetra nibh. Nam sagittis urna sapien, at feugiat lacus pretium eu. Nunc nec nulla ut elit euismod auctor. Duis id consectetur lacus. Aliquam erat volutpat. Aliquam erat volutpat. Pellentesque vulputate, tortor vel blandit fringilla, sem neque blandit lorem, sed vulputate nisl magna in tortor. Sed congue sapien quis est semper, a vehicula mi pulvinar. Curabitur pulvinar sapien id ligula sagittis suscipit. Maecenas dictum, neque ultrices aliquam accumsan, ante nunc vulputate magna, ac ultrices sapien nulla sit amet mi. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Sed sollicitudin lacus diam, sit amet rutrum libero facilisis sed. Curabitur feugiat eleifend nisl ut vulputate. Etiam eleifend nibh eget massa aliquam, at volutpat turpis porttitor. Aliquam diam sapien, vestibulum sed tincidunt quis, interdum vitae quam. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Ut id sodales nunc, vitae porttitor nibh. Vestibulum convallis, neque blandit fringilla aliquam, turpis ligula vehicula leo, eu imperdiet dui dolor sed mi. Maecenas vitae feugiat libero. Quisque ac rutrum augue. Phasellus metus nibh, hendrerit sit amet sapien at, ullamcorper sodales lacus. Vivamus nisi diam, hendrerit nec ligula sed, viverra sagittis massa. Curabitur at libero dolor. Nullam dolor ipsum, tincidunt et libero eu, tristique sodales orci. Praesent lorem ex, sodales non aliquam eu, accumsan at arcu. Proin laoreet eget ex nec efficitur. In sollicitudin tellus enim, et placerat dui scelerisque vitae. Etiam non turpis vel massa ullamcorper aliquam nec quis felis. Suspendisse mollis turpis felis, eget dapibus ante sodales suscipit. Etiam eget diam et mi efficitur lobortis. Proin eu felis ac urna tempus ornare. Donec convallis ligula vel libero porta imperdiet. Donec nibh nunc, pretium et ultrices at, porta vitae lorem. Nunc gravida a felis eget tempor.",
			"hu": "Impressum Hungarian - Duis sit amet tincidunt purus. Cras finibus mi ut purus pretium porttitor. Nunc sed varius ipsum. Cras a pharetra nibh. Nam sagittis urna sapien, at feugiat lacus pretium eu. Nunc nec nulla ut elit euismod auctor. Duis id consectetur lacus. Aliquam erat volutpat. Aliquam erat volutpat. Pellentesque vulputate, tortor vel blandit fringilla, sem neque blandit lorem, sed vulputate nisl magna in tortor. Sed congue sapien quis est semper, a vehicula mi pulvinar. Curabitur pulvinar sapien id ligula sagittis suscipit. Maecenas dictum, neque ultrices aliquam accumsan, ante nunc vulputate magna, ac ultrices sapien nulla sit amet mi. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Sed sollicitudin lacus diam, sit amet rutrum libero facilisis sed. Curabitur feugiat eleifend nisl ut vulputate. Etiam eleifend nibh eget massa aliquam, at volutpat turpis porttitor. Aliquam diam sapien, vestibulum sed tincidunt quis, interdum vitae quam. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Ut id sodales nunc, vitae porttitor nibh. Vestibulum convallis, neque blandit fringilla aliquam, turpis ligula vehicula leo, eu imperdiet dui dolor sed mi. Maecenas vitae feugiat libero. Quisque ac rutrum augue. Phasellus metus nibh, hendrerit sit amet sapien at, ullamcorper sodales lacus. Vivamus nisi diam, hendrerit nec ligula sed, viverra sagittis massa. Curabitur at libero dolor. Nullam dolor ipsum, tincidunt et libero eu, tristique sodales orci. Praesent lorem ex, sodales non aliquam eu, accumsan at arcu. Proin laoreet eget ex nec efficitur. In sollicitudin tellus enim, et placerat dui scelerisque vitae. Etiam non turpis vel massa ullamcorper aliquam nec quis felis. Suspendisse mollis turpis felis, eget dapibus ante sodales suscipit. Etiam eget diam et mi efficitur lobortis. Proin eu felis ac urna tempus ornare. Donec convallis ligula vel libero porta imperdiet. Donec nibh nunc, pretium et ultrices at, porta vitae lorem. Nunc gravida a felis eget tempor.",
		},
		"privacy": {
			"en": "Impressum English - Duis sit amet tincidunt purus. Cras finibus mi ut purus pretium porttitor. Nunc sed varius ipsum. Cras a pharetra nibh. Nam sagittis urna sapien, at feugiat lacus pretium eu. Nunc nec nulla ut elit euismod auctor. Duis id consectetur lacus. Aliquam erat volutpat. Aliquam erat volutpat. Pellentesque vulputate, tortor vel blandit fringilla, sem neque blandit lorem, sed vulputate nisl magna in tortor. Sed congue sapien quis est semper, a vehicula mi pulvinar. Curabitur pulvinar sapien id ligula sagittis suscipit. Maecenas dictum, neque ultrices aliquam accumsan, ante nunc vulputate magna, ac ultrices sapien nulla sit amet mi. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Sed sollicitudin lacus diam, sit amet rutrum libero facilisis sed. Curabitur feugiat eleifend nisl ut vulputate. Etiam eleifend nibh eget massa aliquam, at volutpat turpis porttitor. Aliquam diam sapien, vestibulum sed tincidunt quis, interdum vitae quam. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Ut id sodales nunc, vitae porttitor nibh. Vestibulum convallis, neque blandit fringilla aliquam, turpis ligula vehicula leo, eu imperdiet dui dolor sed mi. Maecenas vitae feugiat libero. Quisque ac rutrum augue. Phasellus metus nibh, hendrerit sit amet sapien at, ullamcorper sodales lacus. Vivamus nisi diam, hendrerit nec ligula sed, viverra sagittis massa. Curabitur at libero dolor. Nullam dolor ipsum, tincidunt et libero eu, tristique sodales orci. Praesent lorem ex, sodales non aliquam eu, accumsan at arcu. Proin laoreet eget ex nec efficitur. In sollicitudin tellus enim, et placerat dui scelerisque vitae. Etiam non turpis vel massa ullamcorper aliquam nec quis felis. Suspendisse mollis turpis felis, eget dapibus ante sodales suscipit. Etiam eget diam et mi efficitur lobortis. Proin eu felis ac urna tempus ornare. Donec convallis ligula vel libero porta imperdiet. Donec nibh nunc, pretium et ultrices at, porta vitae lorem. Nunc gravida a felis eget tempor.",
			"de": "Impressum Deutsch - Duis sit amet tincidunt purus. Cras finibus mi ut purus pretium porttitor. Nunc sed varius ipsum. Cras a pharetra nibh. Nam sagittis urna sapien, at feugiat lacus pretium eu. Nunc nec nulla ut elit euismod auctor. Duis id consectetur lacus. Aliquam erat volutpat. Aliquam erat volutpat. Pellentesque vulputate, tortor vel blandit fringilla, sem neque blandit lorem, sed vulputate nisl magna in tortor. Sed congue sapien quis est semper, a vehicula mi pulvinar. Curabitur pulvinar sapien id ligula sagittis suscipit. Maecenas dictum, neque ultrices aliquam accumsan, ante nunc vulputate magna, ac ultrices sapien nulla sit amet mi. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Sed sollicitudin lacus diam, sit amet rutrum libero facilisis sed. Curabitur feugiat eleifend nisl ut vulputate. Etiam eleifend nibh eget massa aliquam, at volutpat turpis porttitor. Aliquam diam sapien, vestibulum sed tincidunt quis, interdum vitae quam. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Ut id sodales nunc, vitae porttitor nibh. Vestibulum convallis, neque blandit fringilla aliquam, turpis ligula vehicula leo, eu imperdiet dui dolor sed mi. Maecenas vitae feugiat libero. Quisque ac rutrum augue. Phasellus metus nibh, hendrerit sit amet sapien at, ullamcorper sodales lacus. Vivamus nisi diam, hendrerit nec ligula sed, viverra sagittis massa. Curabitur at libero dolor. Nullam dolor ipsum, tincidunt et libero eu, tristique sodales orci. Praesent lorem ex, sodales non aliquam eu, accumsan at arcu. Proin laoreet eget ex nec efficitur. In sollicitudin tellus enim, et placerat dui scelerisque vitae. Etiam non turpis vel massa ullamcorper aliquam nec quis felis. Suspendisse mollis turpis felis, eget dapibus ante sodales suscipit. Etiam eget diam et mi efficitur lobortis. Proin eu felis ac urna tempus ornare. Donec convallis ligula vel libero porta imperdiet. Donec nibh nunc, pretium et ultrices at, porta vitae lorem. Nunc gravida a felis eget tempor.",
			"hu": "Impressum Hungarian - Duis sit amet tincidunt purus. Cras finibus mi ut purus pretium porttitor. Nunc sed varius ipsum. Cras a pharetra nibh. Nam sagittis urna sapien, at feugiat lacus pretium eu. Nunc nec nulla ut elit euismod auctor. Duis id consectetur lacus. Aliquam erat volutpat. Aliquam erat volutpat. Pellentesque vulputate, tortor vel blandit fringilla, sem neque blandit lorem, sed vulputate nisl magna in tortor. Sed congue sapien quis est semper, a vehicula mi pulvinar. Curabitur pulvinar sapien id ligula sagittis suscipit. Maecenas dictum, neque ultrices aliquam accumsan, ante nunc vulputate magna, ac ultrices sapien nulla sit amet mi. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Sed sollicitudin lacus diam, sit amet rutrum libero facilisis sed. Curabitur feugiat eleifend nisl ut vulputate. Etiam eleifend nibh eget massa aliquam, at volutpat turpis porttitor. Aliquam diam sapien, vestibulum sed tincidunt quis, interdum vitae quam. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Ut id sodales nunc, vitae porttitor nibh. Vestibulum convallis, neque blandit fringilla aliquam, turpis ligula vehicula leo, eu imperdiet dui dolor sed mi. Maecenas vitae feugiat libero. Quisque ac rutrum augue. Phasellus metus nibh, hendrerit sit amet sapien at, ullamcorper sodales lacus. Vivamus nisi diam, hendrerit nec ligula sed, viverra sagittis massa. Curabitur at libero dolor. Nullam dolor ipsum, tincidunt et libero eu, tristique sodales orci. Praesent lorem ex, sodales non aliquam eu, accumsan at arcu. Proin laoreet eget ex nec efficitur. In sollicitudin tellus enim, et placerat dui scelerisque vitae. Etiam non turpis vel massa ullamcorper aliquam nec quis felis. Suspendisse mollis turpis felis, eget dapibus ante sodales suscipit. Etiam eget diam et mi efficitur lobortis. Proin eu felis ac urna tempus ornare. Donec convallis ligula vel libero porta imperdiet. Donec nibh nunc, pretium et ultrices at, porta vitae lorem. Nunc gravida a felis eget tempor.",
		},
		"methodology": {
			"en": "Impressum English - Duis sit amet tincidunt purus. Cras finibus mi ut purus pretium porttitor. Nunc sed varius ipsum. Cras a pharetra nibh. Nam sagittis urna sapien, at feugiat lacus pretium eu. Nunc nec nulla ut elit euismod auctor. Duis id consectetur lacus. Aliquam erat volutpat. Aliquam erat volutpat. Pellentesque vulputate, tortor vel blandit fringilla, sem neque blandit lorem, sed vulputate nisl magna in tortor. Sed congue sapien quis est semper, a vehicula mi pulvinar. Curabitur pulvinar sapien id ligula sagittis suscipit. Maecenas dictum, neque ultrices aliquam accumsan, ante nunc vulputate magna, ac ultrices sapien nulla sit amet mi. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Sed sollicitudin lacus diam, sit amet rutrum libero facilisis sed. Curabitur feugiat eleifend nisl ut vulputate. Etiam eleifend nibh eget massa aliquam, at volutpat turpis porttitor. Aliquam diam sapien, vestibulum sed tincidunt quis, interdum vitae quam. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Ut id sodales nunc, vitae porttitor nibh. Vestibulum convallis, neque blandit fringilla aliquam, turpis ligula vehicula leo, eu imperdiet dui dolor sed mi. Maecenas vitae feugiat libero. Quisque ac rutrum augue. Phasellus metus nibh, hendrerit sit amet sapien at, ullamcorper sodales lacus. Vivamus nisi diam, hendrerit nec ligula sed, viverra sagittis massa. Curabitur at libero dolor. Nullam dolor ipsum, tincidunt et libero eu, tristique sodales orci. Praesent lorem ex, sodales non aliquam eu, accumsan at arcu. Proin laoreet eget ex nec efficitur. In sollicitudin tellus enim, et placerat dui scelerisque vitae. Etiam non turpis vel massa ullamcorper aliquam nec quis felis. Suspendisse mollis turpis felis, eget dapibus ante sodales suscipit. Etiam eget diam et mi efficitur lobortis. Proin eu felis ac urna tempus ornare. Donec convallis ligula vel libero porta imperdiet. Donec nibh nunc, pretium et ultrices at, porta vitae lorem. Nunc gravida a felis eget tempor.",
			"de": "Impressum Deutsch - Duis sit amet tincidunt purus. Cras finibus mi ut purus pretium porttitor. Nunc sed varius ipsum. Cras a pharetra nibh. Nam sagittis urna sapien, at feugiat lacus pretium eu. Nunc nec nulla ut elit euismod auctor. Duis id consectetur lacus. Aliquam erat volutpat. Aliquam erat volutpat. Pellentesque vulputate, tortor vel blandit fringilla, sem neque blandit lorem, sed vulputate nisl magna in tortor. Sed congue sapien quis est semper, a vehicula mi pulvinar. Curabitur pulvinar sapien id ligula sagittis suscipit. Maecenas dictum, neque ultrices aliquam accumsan, ante nunc vulputate magna, ac ultrices sapien nulla sit amet mi. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Sed sollicitudin lacus diam, sit amet rutrum libero facilisis sed. Curabitur feugiat eleifend nisl ut vulputate. Etiam eleifend nibh eget massa aliquam, at volutpat turpis porttitor. Aliquam diam sapien, vestibulum sed tincidunt quis, interdum vitae quam. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Ut id sodales nunc, vitae porttitor nibh. Vestibulum convallis, neque blandit fringilla aliquam, turpis ligula vehicula leo, eu imperdiet dui dolor sed mi. Maecenas vitae feugiat libero. Quisque ac rutrum augue. Phasellus metus nibh, hendrerit sit amet sapien at, ullamcorper sodales lacus. Vivamus nisi diam, hendrerit nec ligula sed, viverra sagittis massa. Curabitur at libero dolor. Nullam dolor ipsum, tincidunt et libero eu, tristique sodales orci. Praesent lorem ex, sodales non aliquam eu, accumsan at arcu. Proin laoreet eget ex nec efficitur. In sollicitudin tellus enim, et placerat dui scelerisque vitae. Etiam non turpis vel massa ullamcorper aliquam nec quis felis. Suspendisse mollis turpis felis, eget dapibus ante sodales suscipit. Etiam eget diam et mi efficitur lobortis. Proin eu felis ac urna tempus ornare. Donec convallis ligula vel libero porta imperdiet. Donec nibh nunc, pretium et ultrices at, porta vitae lorem. Nunc gravida a felis eget tempor.",
			"hu": "Impressum Hungarian - Duis sit amet tincidunt purus. Cras finibus mi ut purus pretium porttitor. Nunc sed varius ipsum. Cras a pharetra nibh. Nam sagittis urna sapien, at feugiat lacus pretium eu. Nunc nec nulla ut elit euismod auctor. Duis id consectetur lacus. Aliquam erat volutpat. Aliquam erat volutpat. Pellentesque vulputate, tortor vel blandit fringilla, sem neque blandit lorem, sed vulputate nisl magna in tortor. Sed congue sapien quis est semper, a vehicula mi pulvinar. Curabitur pulvinar sapien id ligula sagittis suscipit. Maecenas dictum, neque ultrices aliquam accumsan, ante nunc vulputate magna, ac ultrices sapien nulla sit amet mi. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Sed sollicitudin lacus diam, sit amet rutrum libero facilisis sed. Curabitur feugiat eleifend nisl ut vulputate. Etiam eleifend nibh eget massa aliquam, at volutpat turpis porttitor. Aliquam diam sapien, vestibulum sed tincidunt quis, interdum vitae quam. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Ut id sodales nunc, vitae porttitor nibh. Vestibulum convallis, neque blandit fringilla aliquam, turpis ligula vehicula leo, eu imperdiet dui dolor sed mi. Maecenas vitae feugiat libero. Quisque ac rutrum augue. Phasellus metus nibh, hendrerit sit amet sapien at, ullamcorper sodales lacus. Vivamus nisi diam, hendrerit nec ligula sed, viverra sagittis massa. Curabitur at libero dolor. Nullam dolor ipsum, tincidunt et libero eu, tristique sodales orci. Praesent lorem ex, sodales non aliquam eu, accumsan at arcu. Proin laoreet eget ex nec efficitur. In sollicitudin tellus enim, et placerat dui scelerisque vitae. Etiam non turpis vel massa ullamcorper aliquam nec quis felis. Suspendisse mollis turpis felis, eget dapibus ante sodales suscipit. Etiam eget diam et mi efficitur lobortis. Proin eu felis ac urna tempus ornare. Donec convallis ligula vel libero porta imperdiet. Donec nibh nunc, pretium et ultrices at, porta vitae lorem. Nunc gravida a felis eget tempor.",
		},
	};
	for k, lMap := range blockCont {
		for lk, c := range lMap {
			block, err := b.GetByIdentifier(k, lk);
			h.Error(err, "", h.ERROR_LVL_ERROR);
			if (block.Id == 0) {
				block.Identifier = k;
				block.Lc = lk;
				block.Title = k;
				block.Content = c;
				err := dbmap.Insert(&block);
				h.Error(err, "", h.ERROR_LVL_ERROR)
			}
		}
	}
}

func (b Block) PrepeareData(){}
