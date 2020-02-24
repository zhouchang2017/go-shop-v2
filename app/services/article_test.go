package services

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"go-shop-v2/pkg/db/mongodb"
	"go-shop-v2/pkg/qiniu"
	"testing"
)

func TestArticleService_Create(t *testing.T) {
	mongodb.TestConnect()
	defer mongodb.Close()

	service := MakeArticleService()

	article, err := service.Create(context.Background(), ArticleOption{
		Title:      "不要“恐鄂”，要妥善安置在外的湖北人",
		ShortTitle: "为包括武汉在内的湖北游客提供指定酒店，既能体现对同胞的温情，也能有力地防止疫情蔓延。",
		Photos:     []qiniu.Image{"https://inews.gtimg.com/newsapp_bt/0/11258845869/1000"},
		Content: `1月27日，一张“@来西安的武汉人——请在西安入住指定酒店”的二维码图片在网络流传。

在一些地方“谈鄂色变”、标签化、妖魔化武汉人，武汉返乡人员被避之不及的语境下，西安等地为赴当地的武汉人提供指定酒店的做法值得肯定。面对有可能存在的疫情传播，这些地方不是一味地“堵”、“赶”，而是积极为来自武汉的民众提供准确的住宿支持，这提振了滞留外地武汉人对抗疫情的信心。

当然，最重要的是，这也是当前公共卫生防疫至关重要的一环。

就在1月26日晚间，武汉市长披露，受春节和疫情的影响，目前有500多万人离开武汉。按照常识以及大数据分析，春节正常返乡的大学生与外来务工人员，应该占了不小的比例。但除此之外，显然还有大批人离开武汉，是为了旅游、探亲或参加其他活动。

虽然这部分人到底是多少，尚无明晰的统计数据，但以500多万如此庞大的基数而言，哪怕仅占1%的概率，也是个不小的数字。在当前疫情防控的关键时刻，与其让他们以个体原子化的形式分散在各地，不如通过指定入住酒店等形式集中管理。

一方面，这也是保障这部分人基本生活需求的必要措施。从社交平台来看，在武汉交通关闭后，一夜之间，不少来自武汉的人从“旅客变成了不受欢迎的武汉人”。他们既回不了家，也被一些酒店拒之门外，甚至有人“在机场睡了一宿，不知道明天要如何”。

就在27日凌晨，武汉文旅局发布请求信，发出了希望各地“对所有在外旅行的武汉市民给予必要的帮助”的呼声，这也从侧面印证了一些武汉游客确实深陷窘境的现实。


另一方面，为武汉游客指定入住酒店，也是一种精准管控追踪的措施，从根本上切断了他们与外界密切接触的条件。此举也便于以酒店为网格，加强对他们的健康监测，一旦发现情况也能第一时间启动诊断、隔离治疗等程序，并进一步做好病人的隔离控制和转送至定点医院的各项准备。

一言以蔽之，为包括武汉在内的湖北游客提供指定酒店，既能体现对同胞的温情，也能有力地防止疫情蔓延。

实际上，除了西安包括云南省各市州、广东省湛江市、广西桂林、上海、海南海口、厦门鼓浪屿等地的有关部门也早已开始行动，为滞留当地的武汉游客提供对接宾馆，甚至是免费的住宿服务。

与此同时，也有一些人开始“自救”，朋友圈就流传着不少“**省武汉同胞回家群”的微信群，试图通过民间力量，来寻找愿意收留武汉人的旅馆。

大疫面前，民间自助自然值得肯定，但来自各级各地政府层面的统筹安排，更是不可或缺。因此，其他尚未行动起来的城市，不妨也借鉴这些“走在前列”的省市的做法，通过为滞留在外的湖北游客提供指定酒店等形式，给同胞多些温暖，也给防疫多些助益。

编辑：陈静 校对：刘军`,
		ProductId: "5e268351001f5053d8b2b1e0",
	})

	if err != nil {
		t.Fatal(err)
	}

	spew.Dump(article.GetID())
}