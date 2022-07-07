package global

var search SearcherClient

type SearcherClient struct {
	Url            string
	AppId          string
	AppKey         string
	AppType        string
	QueryCondition map[string]string
	Columns        []string
	Wheres         map[string]string
	Orders         map[string]string
	TimeOut        int
}

func init() {
	search.Url = "http://info.grid.leju.com/search/default/index"

}

func getType(poolType string) (err error) {
	switch poolType {
	case "etc_goods":
		search.AppId = ""
		search.AppKey = ""
		search.AppType = "etc_goods"
	case "house":
		search.AppId = ""
		search.AppKey = ""
		search.AppType = "house"
	}
	return
}
