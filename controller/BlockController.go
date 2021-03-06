package controller

import (
	"github.com/valyala/fasthttp"
	"html/template"
	h "mos/helper"
	"mos/model/list"
	m "mos/model"
	"mos/db"
	"strconv"
	"fmt"
	"mos/model/view/admin"
	"mos/model/view"
)

type BlockController struct {
	AuthAction map[string][]string;
}

func (b *BlockController) Init() {
	b.AuthAction = make(map[string][]string);
	b.AuthAction["edit"] = []string{"block/edit"};
	b.AuthAction["delete"] = []string{"block/delete"};
	b.AuthAction["save"] = []string{"block/edit", "block/new"};
	b.AuthAction["new"] = []string{"block/new"};
	b.AuthAction["list"] = []string{"block/list"};
}

func (b *BlockController) ListAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(b.AuthAction["list"],session)) {
		var bl list.BlockList;
		bl.Init(ctx, session.GetActiveLang());
		pageInstance.Title = "List Blocks"

		AdminContent := admin.Content{};
		AdminContent.Title = "Blocks"
		AdminContent.SubTitle = "List Blocks"

		AdminContent.Content = template.HTML(bl.Render(bl.GetToPage()))
		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent,pageInstance.Scope), "", nil, false, 0)
	} else {
		Redirect(ctx, "user/welcome", fasthttp.StatusForbidden, true,pageInstance);
		return;
	}
}

func (b *BlockController) EditAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(b.AuthAction["edit"],session)) {
		//azért nem kell vizsgálni az errort, mert a request reguláris kifejezése csak akkor hozza ide, ha a végén \d van :)
		var id, _ = strconv.Atoi(h.GetParamFromCtxPath(ctx, 3, ""));
		var blockId = int64(id);
		var ModelBlock m.Block;
		Block, err := ModelBlock.Get(blockId);
		if (err != nil) {
			session.AddError(err.Error());
			h.Error(err, "", h.ERROR_LVL_WARNING);
			Redirect(ctx, "block/index", fasthttp.StatusOK, true,pageInstance);
			return;
		}

		var data map[string]interface{};
		if (!ctx.IsPost()) {
			data = map[string]interface{}{
				"id":         strconv.Itoa(int(Block.Id)),
				"identifier": Block.Identifier,
				"title": Block.Title,
				"content":    Block.Content,
				"lc" : Block.Lc,
			};
		} else {
			data = map[string]interface{}{
				"id":         h.GetFormData(ctx, "id", false).(string),
				"identifier": h.GetFormData(ctx, "identifier", false).(string),
				"title": h.GetFormData(ctx, "title", false).(string),
				"content":    h.GetFormData(ctx, "content", false).(string),
				"lc": h.GetFormData(ctx, "lc", false).(string),
			};
		}

		var form = m.GetBlockForm(data, fmt.Sprintf("block/edit/%v", data["id"].(string)));
		if (ctx.IsPost()) {
			succ, formErrors := b.saveBlock(ctx, session, Block);
			form.SetErrors(formErrors);
			if (succ) {
				session.AddSuccess("Block save was successful.");
				Redirect(ctx, fmt.Sprintf("block/edit/%v", data["id"].(string)), fasthttp.StatusOK, true,pageInstance);
				return;
			}
		}

		pageInstance.Title = "Block - Edit"

		AdminContent := admin.Content{};
		AdminContent.Title = "Block"
		AdminContent.SubTitle = fmt.Sprintf("Edit block %v", Block.Identifier);
		AdminContent.Content = template.HTML(form.Render())
		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent,pageInstance.Scope), "", nil, false, 0)
	} else {
		Redirect(ctx, "block", fasthttp.StatusForbidden, true,pageInstance)
		return;
	}
}

func (b *BlockController) NewAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(b.AuthAction["new"],session)) {
		var Block = m.NewEmptyBlock();
		var data map[string]interface{} = map[string]interface{}{};
		var dataKeys []string = []string{"id", "identifier","title", "content"};
		for _,k := range dataKeys {
			var val string = "";
			if(ctx.IsPost()) {
				val = h.GetFormData(ctx, k, false).(string)
			}
			data[k] = val;
		}
		data["lc"] = session.GetActiveLang();

		var form = m.GetBlockForm(data, "block/new");
		if (ctx.IsPost()) {
			succ, formErrors := b.saveBlock(ctx, session, Block);
			form.SetErrors(formErrors);
			if (succ) {
				session.AddSuccess("Block save was successful.");
				Redirect(ctx, "block", fasthttp.StatusOK, true,pageInstance);
				return;
			}
		}

		pageInstance.Title = "Block - New"

		AdminContent := admin.Content{};
		AdminContent.Title = "Block"
		AdminContent.SubTitle = "New";
		AdminContent.Content = template.HTML(form.Render())
		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent,pageInstance.Scope), "", nil, false, 0)
	} else {
		Redirect(ctx, "", fasthttp.StatusForbidden, true,pageInstance)
		return;
	}
}

func (b *BlockController) saveBlock(ctx *fasthttp.RequestCtx, session *h.Session, Block m.Block) (bool, map[string]error) {
	if (ctx.IsPost() && ((Ah.HasRights(b.AuthAction["edit"],session) && Block.Id != 0) || (Ah.HasRights(b.AuthAction["new"],session) && Block.Id == 0))) {
		var err error;
		var succ bool;
		var Validator = m.GetBlockFormValidator(ctx, Block);
		succ, errors := Validator.Validate();
		if (!succ) {
			return false, errors;
		}

		Block.Identifier = h.GetFormData(ctx, "identifier", false).(string);
		Block.Title = h.GetFormData(ctx, "title", false).(string);
		Block.Content = h.GetFormData(ctx, "content", false).(string);
		Block.Lc = h.GetFormData(ctx, "lc", false).(string);

		if (Block.Id > 0) {
			_, err = db.DbMap.Update(&Block);
		} else {
			err = db.DbMap.Insert(&Block);
		}
		h.Error(err, "", h.ERROR_LVL_ERROR)
		succ = err == nil;
		return succ, nil;
	} else {
		return false, nil;
	}
}

func (b *BlockController) DeleteAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	var status int;
	if (Ah.HasRights(b.AuthAction["delete"],session)) {
		//azért nem kell vizsgálni az errort, mert a request reguláris kifejezése csak akkor hozza ide, ha a végén \d van :)
		var id, _ = strconv.Atoi(h.GetParamFromCtxPath(ctx, 3, ""));
		var blockId = int64(id);
		var ModelBlock m.Block;
		Block, err := ModelBlock.Get(blockId);
		if (err != nil) {
			session.AddError(err.Error());
			h.Error(err, "", h.ERROR_LVL_WARNING);
			Redirect(ctx, "block/index", fasthttp.StatusOK, true,pageInstance);
			return;
		}

		blockIdentifier := Block.Identifier;
		count,err := db.DbMap.Delete(&Block);
		h.Error(err,"",h.ERROR_LVL_WARNING);
		if(err != nil || count == 0){
			session.AddError("An error occured, could not delete the block.");
			status = fasthttp.StatusBadRequest;
		} else {
			session.AddSuccess(fmt.Sprintf("Block %v has been deleted",blockIdentifier));
			status = fasthttp.StatusOK;
		}
	} else {
		status = fasthttp.StatusOK;
	}
	Redirect(ctx, "block/index", status, true,pageInstance)
}
