package main

import (
	c "mos/controller"
	dbhelper "mos/db"
	h "mos/helper"
	m "mos/model"

	"fmt"
	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/valyala/fasthttp"
	"time"
)

func main() {
	var Log = h.SetLog()
	var Conf h.Conf = h.GetConfig()
	h.PrintlnIf(fmt.Sprintf("Config values are the following: %+v", Conf), Conf.Mode.Debug)
	h.InitLanguage()
	dbhelper.InitDb()
	mapTables(dbhelper.DbMap)

	models := buildStructure()

	h.SecureCookieSet()
	c.InitControllers()
	defer func() {
		go func() {
			if(h.GetConfig().Mode.Rebuild_data) {
				m.TruncateTables();
				m.MatchOperators()
				var start time.Time = time.Now().Round(time.Second)
				for _, curM := range models {
					curM.PrepeareData()
				}
				var end time.Time = time.Now().Round(time.Second)
				var dif= end.Unix() - start.Unix()
				var minutes= dif / 60
				var seconds= dif - (minutes * 60)
				h.PrintlnIf(fmt.Sprintf("%v:%v passed since data process started", minutes, seconds), h.GetConfig().Mode.Debug)
				h.PrintlnIf("Done importing and prepearing database.", h.GetConfig().Mode.Debug)
			}
		}()

		srv := &fasthttp.Server{
			Name:         "Branditorial Server",
			ReadTimeout:  time.Duration(h.GetConfig().Server.ReadTimeoutSeconds) * time.Second,
			WriteTimeout: time.Duration(h.GetConfig().Server.WriteTimeoutSeconds) * time.Second,
			Handler:      c.Route,
		}

		err := srv.ListenAndServe(fmt.Sprintf(":%s", Conf.ListenPort))
		h.Error(err, "", h.ERROR_LVL_ERROR)
		h.PrintlnIf("The server is listening...", h.GetConfig().Mode.Debug)
		Log.Close()
	}()
}

func buildStructure() []m.DbInterface {
	if h.GetConfig().Mode.Rebuild_structure {
		defer h.PrintlnIf("STRUCTURE BUILDING DONE", h.GetConfig().Mode.Debug)
	}

	var models []m.DbInterface = []m.DbInterface{
		m.Ban{},
		m.Request{},
		m.Status{},
		m.UserRole{},
		m.User{},
		m.Config{},
		m.Block{},
		// data
		m.MediaType{},
		m.Media{},
		m.Owner{},
		m.Operator{},
		m.Interest{},
		m.MediaOwner{},
		m.MediaOperator{},
		m.OperatorYearData{},
		m.OperatorInterest{},
	}

	h.PrintlnIf("Rebuild database structure because config rebuild flag is true", h.GetConfig().Mode.Rebuild_structure)


	if(h.GetConfig().Mode.Rebuild_structure) {
		_, err := dbhelper.DbMap.Exec("SET FOREIGN_KEY_CHECKS=0")
		h.Error(err, "", h.ERROR_LVL_ERROR)
		for _, mod := range models {
			mod.BuildStructure()
		}
		_, err = dbhelper.DbMap.Exec("SET FOREIGN_KEY_CHECKS=1")
		h.Error(err, "", h.ERROR_LVL_ERROR)
	}

	return models
}

func mapTables(dbmap *gorp.DbMap) {
	var ban m.Ban
	TableMap := dbmap.AddTableWithName(ban, ban.GetTable())
	TableMap.SetKeys(true, ban.GetPrimaryKey()...)

	var block m.Block
	TableMap = dbmap.AddTableWithName(block, block.GetTable())
	TableMap.SetKeys(true, block.GetPrimaryKey()...)

	var conf m.Config
	TableMap = dbmap.AddTableWithName(conf, conf.GetTable())
	TableMap.SetKeys(true, conf.GetPrimaryKey()...)

	var med m.Media
	TableMap = dbmap.AddTableWithName(med, med.GetTable())
	TableMap.SetKeys(false, med.GetPrimaryKey()...)

	var mo m.MediaOperator
	TableMap = dbmap.AddTableWithName(mo, mo.GetTable())
	TableMap.SetKeys(true, mo.GetPrimaryKey()...)

	var mow m.MediaOwner
	TableMap = dbmap.AddTableWithName(mow, mow.GetTable())
	TableMap.SetKeys(true, mow.GetPrimaryKey()...)

	var mt m.MediaType
	TableMap = dbmap.AddTableWithName(mt, mt.GetTable())
	TableMap.SetKeys(true, mt.GetPrimaryKey()...)

	var o m.Operator
	TableMap = dbmap.AddTableWithName(o, o.GetTable())
	TableMap.SetKeys(false, o.GetPrimaryKey()...)

	var i m.Interest
	TableMap = dbmap.AddTableWithName(i, i.GetTable())
	TableMap.SetKeys(true, i.GetPrimaryKey()...)

	var oyd m.OperatorYearData
	TableMap = dbmap.AddTableWithName(oyd, oyd.GetTable())
	TableMap.SetKeys(false, oyd.GetPrimaryKey()...)

	var oi m.OperatorInterest
	TableMap = dbmap.AddTableWithName(oi, oi.GetTable())
	TableMap.SetKeys(true, oi.GetPrimaryKey()...)

	var ow m.Owner
	TableMap = dbmap.AddTableWithName(ow, ow.GetTable())
	TableMap.SetKeys(false, ow.GetPrimaryKey()...)

	var req m.Request
	TableMap = dbmap.AddTableWithName(req, req.GetTable())
	TableMap.SetKeys(true, req.GetPrimaryKey()...)

	var s m.Status
	TableMap = dbmap.AddTableWithName(s, s.GetTable())
	TableMap.SetKeys(true, s.GetPrimaryKey()...)

	var u m.User
	TableMap = dbmap.AddTableWithName(u, u.GetTable())
	TableMap.SetKeys(true, u.GetPrimaryKey()...)

	var ur m.UserRole
	TableMap = dbmap.AddTableWithName(ur, ur.GetTable())
	TableMap.SetKeys(true, ur.GetPrimaryKey()...)
}
