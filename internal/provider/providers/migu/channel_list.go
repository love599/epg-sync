package migu

import "github.com/epg-sync/epgsync/internal/model"

var channelList = []*model.ProviderChannel{
	{
		ID:      "608807420",
		Name:    "CCTV-1综合",
		Aliases: []string{"CCTV1", "CCTV-1"},
	},
	{
		ID:      "631780532",
		Name:    "CCTV-2财经",
		Aliases: []string{"CCTV2", "CCTV-2"},
	},
	{
		ID:      "624878271",
		Name:    "CCTV-3综艺",
		Aliases: []string{"CCTV3", "CCTV-3"},
	},
	{
		ID:      "631780421",
		Name:    "CCTV-4中文国际",
		Aliases: []string{"CCTV4亚洲", "CCTV-4中文国际", "CCTV-4中文国际(亚)"},
	},
	{
		ID:      "641886683",
		Name:    "CCTV-5体育",
		Aliases: []string{"CCTV5", "CCTV-5"},
	},
	{
		ID:      "641886773",
		Name:    "CCTV-5体育赛事",
		Aliases: []string{"CCTV5+", "CCTV-5+", "CCTV5plus"},
	},
	{
		ID:      "624878396",
		Name:    "CCTV-6电影",
		Aliases: []string{"CCTV6", "CCTV-6"},
	},
	{
		ID:      "673168121",
		Name:    "CCTV-7军事农业",
		Aliases: []string{"CCTV7", "CCTV-7"},
	},
	{
		ID:      "624878356",
		Name:    "CCTV-8电视剧",
		Aliases: []string{"CCTV8", "CCTV-8"},
	},
	{
		ID:      "673168140",
		Name:    "CCTV-9纪录",
		Aliases: []string{"CCTV9", "CCTV-9"},
	},
	{
		ID:      "624878405",
		Name:    "CCTV-10科教",
		Aliases: []string{"CCTV10", "CCTV-10"},
	},
	{
		ID:      "667987558",
		Name:    "CCTV-11戏曲",
		Aliases: []string{"CCTV11", "CCTV-11"},
	},
	{
		ID:      "673168185",
		Name:    "CCTV-12社会与法",
		Aliases: []string{"CCTV12", "CCTV-12"},
	},
	{
		ID:      "608807423",
		Name:    "CCTV-13新闻",
		Aliases: []string{"CCTV13", "CCTV-13"},
	},
	{
		ID:      "624878440",
		Name:    "CCTV-14少儿",
		Aliases: []string{"CCTV14", "CCTV-14"},
	},
	{
		ID:      "673168223",
		Name:    "CCTV-15音乐",
		Aliases: []string{"CCTV15", "CCTV-15"},
	},
	{
		ID:      "673168256",
		Name:    "CCTV-17农业农村",
		Aliases: []string{"CCTV17", "CCTV-17"},
	},
	{
		ID:   "608807419",
		Name: "CCTV4欧洲",
	},
	{
		ID:   "608807416",
		Name: "CCTV4美洲",
	},

	{
		ID:   "609017205",
		Name: "CGTN",
	},
	{
		ID:   "609006476",
		Name: "CGTN法语",
	},
	{
		ID:   "609006446",
		Name: "CGTN俄语",
	},
	{
		ID:      "609154345",
		Name:    "CGTN阿拉伯语",
		Aliases: []string{"CGTN阿语"},
	},
	{
		ID:      "609006450",
		Name:    "CGTN西班牙语",
		Aliases: []string{"CGTN西语"},
	},
	{
		ID:      "609006487",
		Name:    "CGTN外语纪录",
		Aliases: []string{"CGTN记录"},
	},
	{
		ID:   "651632648",
		Name: "东方卫视",
	},
	{
		ID:   "623899368",
		Name: "江苏卫视",
	},
	{
		ID:   "608831231",
		Name: "广东卫视",
	},
	{
		ID:   "783847495",
		Name: "江西卫视",
	},
	{
		ID:   "790187291",
		Name: "河南卫视",
	},
	{
		ID:   "738910838",
		Name: "陕西卫视",
	},
	{
		ID:   "608917627",
		Name: "大湾区卫视",
	},
	{
		ID:   "947472496",
		Name: "湖北卫视",
	},
	{
		ID:   "947472500",
		Name: "吉林卫视",
	},
	{
		ID:   "947472506",
		Name: "青海卫视",
	},
	{
		ID:   "849116810",
		Name: "东南卫视",
	},
	{
		ID:   "947472502",
		Name: "海南卫视",
	},
	{
		ID:   "849119120",
		Name: "海峡卫视",
	},
	{
		ID:   "956904896",
		Name: "中国农林卫视",
	},
	{
		ID:   "956923145",
		Name: "兵团卫视",
	},
	{
		ID:   "630291707",
		Name: "辽宁卫视",
	},
	{
		ID:   "738910535",
		Name: "宁夏卫视",
	},

	{
		ID:   "644368714",
		Name: "CHC动作电影",
	},
	{
		ID:   "644368373",
		Name: "CHC家庭影院",
	},
	{
		ID:   "952383261",
		Name: "CHC影迷电影",
	},
	{
		ID:   "651632657",
		Name: "上海新闻综合",
	},
	{
		ID:      "617290047",
		Name:    "上视东方影视",
		Aliases: []string{"东方影视"},
	},
	{
		ID:   "838109047",
		Name: "南京新闻综合频道",
	},
	{
		ID:   "838153729",
		Name: "南京教科频道",
	},
	{
		ID:   "838151753",
		Name: "南京十八频道",
	},
	{
		ID:   "626064707",
		Name: "体育休闲频道",
	},
	{
		ID:   "626064714",
		Name: "江苏城市频道",
	},
	{
		ID:   "626064674",
		Name: "江苏国际",
	},
	{
		ID:   "628008321",
		Name: "江苏教育",
	},
	{
		ID:   "626064697",
		Name: "江苏影视频道",
	},
	{
		ID:   "626065193",
		Name: "江苏综艺频道",
	},
	{
		ID:   "626064693",
		Name: "江苏新闻",
	},
	{
		ID:   "639731825",
		Name: "盐城新闻综合",
	},
	{
		ID:   "639731826",
		Name: "淮安综合",
	},
	{
		ID:   "639731818",
		Name: "泰州新闻综合",
	},
	{
		ID:   "639731715",
		Name: "连云港新闻综合",
	},
	{
		ID:   "639731832",
		Name: "宿迁新闻综合",
	},

	{
		ID:   "639731747",
		Name: "徐州新闻综合",
	},
	{
		ID:   "626064703",
		Name: "优漫卡通频道",
	},
	{
		ID:   "955227979",
		Name: "江阴新闻综合",
	},
	{
		ID:   "955227985",
		Name: "南通新闻综合",
	},
	{
		ID:   "955227996",
		Name: "宜兴新闻综合",
	},

	{
		ID:   "639737327",
		Name: "溧水新闻综合",
	},
	{
		ID:   "956909362",
		Name: "陕西银龄频道",
	},
	{
		ID:   "956909358",
		Name: "陕西都市青春频道",
	},

	{
		ID:   "956909356",
		Name: "陕西体育休闲频道",
	},
	{
		ID:   "956909303",
		Name: "陕西秦腔频道",
	},
	{
		ID:   "956909289",
		Name: "陕西新闻资讯频道",
	},
	{
		ID:   "956923159",
		Name: "财富天下",
	},
	{
		ID:   "959986621",
		Name: "中国天气",
	},
}
