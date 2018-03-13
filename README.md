# wss

## 游戏地图

00|0|1|2|3|4|5|6|7|8|9|10|
---|---|---|---|---|---|---|---|---|---|---|---
0|起点|沈阳|天津|核能发电|-|-|入狱|银川|兰州|大陆运输
1|钓鱼岛|-|-|北京|-|-|长沙|-|-|拉萨
2|上海|-|-|新闻|大连|贵阳|建筑公司|-|-|运气
3|税务|-|-|-|-|-|-|-|-|杭州
4|澳门|-|-|-|-|-|-|-|-|电视台
5|运气|-|-|公园|成都|重庆|新闻|-|-|南京
6|香港|-|-|深圳|-|-|航空运输|-|-|苏州
7|大洋运输|三亚|广州|污水处理|-|-|监狱|台北|厦门|证券


00|0|1|2|3|4|5|6|7|8|9|10|
---|---|---|---|---|---|---|---|---|---|---|---
0|起点(0,0)|沈阳(0,1)|天津(0,2)|投资(0,3)|-|-|入狱(0,6)|银川(0,7)|兰州(0,8)|投资(0,9)
1|钓鱼岛(1,0)|-|-|北京(1,3)|-|-|长沙(1,6)|-|-|拉萨(1,9)
2|上海(2,0)|-|-|新闻(2,3)|大连(2,4)|贵阳(2,5)|投资(2,6)|-|-|运气(2,9)
3|税务(3,0)|-|-|-|-|-|-|-|-|杭州(3,9)
4|澳门(4,0)|-|-|-|-|-|-|-|-|投资(4,9)
5|运气(5,0)|-|-|公园(5,3)|成都(5,4)|重庆(5,5)|新闻(5,6)|-|-|南京(5,9)
6|香港(6,0)|-|-|深圳(6,3)|-|-|投资(6,6)|-|-|苏州(6,9)
7|投资(7,0)|三亚(7,1)|广州(7,2)|投资(7,3)|-|-|监狱(7,6)|台北(7,7)|厦门(7,8)|证券(7,9)



## 游戏规则
- 初始化游戏，四人房间
  - 初始化给15000
  - 过了起点给3000
- A投掷骰子
  - A移动到指定位置
  - A购买地产
  - A升级地产
  - A缴纳租金
  - A赎回地产
  - A抵押地产
  - A
- B投掷骰子
  - B移动到指定位置
  - B购买地产
  - B升级地产
  - B缴纳租金
  - B赎回地产
  - B抵押地产
  - B
- 运气卡
  - 随机抽取卡片
  - 银行发钱
  - ......
- 新闻卡
  - 随机抽取卡片
  - 公园捡钱
  - ......
- 游戏结束
  - 某人钱足够多
  - 其它人都破产

- 某人破产
  - 清空其名下地产
  - 把所有的钱给银行，或者缴纳租金
  - 
- 额外运气卡新闻卡规则
  - 入狱
  - 抓入监狱
  - 倒退规则
  - 
  
## 规则说明

- 这是一款简单轻松的桌面游戏，你和你的小伙伴将成为游戏中的地产大亨，购买中国的各大名城
从其他玩家手中收取租金，拥有最多财富的玩家将成为真正的超级地产富翁
- 地产地契30张
- 运气卡牌12张
- 新闻卡12张
- 每个玩家从起点，准备出发
- 每个玩家获取初始资金15000元，将多余的钱放置在一边作为银行
- 将12张 运气卡和12张新闻卡 洗匀放置在地图相应位置作为等待抽取的牌堆
- 将地契按牌色分开排列在桌子的一边，等待玩家购买
- 从钱包最大的玩家开始，依次进行游戏行动回合
- 轮到玩家移动时，先投掷骰子一次，然后按照骰子的点数像箭头方向移动相应的步数，
最后执行人物停下时到达的格子的效果
- 玩家到达无人拥有的地产或投资项目时，可以支付金钱向银行购买该地产或投资项目，
购买之后，将该出的“地契”放在自己面前，表示该地产归你所有
- 当你再次到达自己的地产时，可以支付金钱将地产升级，升级时，将一个房屋标记物放置在
对应的位置上，来表示该地产的等级
- 其它玩家到达你的地产，时他必须向你支付该处地产的地租，
支付地租的金额由该地产的等级决定
- 玩家到达其余地点时，按照地点格子内的文字描述来执行效果
- 入狱：立即进入监狱
- 税务中心：缴纳每块地产300元税金，此外若你拥有2个或以上投资
项目，额外缴纳1000元
- 证券中心：获得你拥有投资项目数量*500元的奖励
- 运气：抽取1张运气卡
- 新闻：抽取1张新闻卡
- 起点：经过时领取奖励3000元


- 停留一轮。当你在监狱时，不能收取租金或获得任何金钱奖励
- 公园捡到300元
- 当玩家需要支付金钱但现金不足时，玩家必须将手中的“地契”按照其抵押价格向银行抵押。
将地契暂时背面向上放置，同时从银行获得金钱。被抵押的地产不能收取地租，玩家可以在自己
的行动回合内，再次支付金钱给银行以赎回被抵押地契，赎回金额为抵押价格+地产售价的10%
- 当玩家无法支付金钱，并且没有任何卡牌可以抵押时，该玩家破产。破产的玩家作为失败方离
开游戏，他的所有卡牌都归还给银行所有。
- 当任意玩家的现金达到50000元时，获得游戏的胜利
- 当其余玩家宣告破产时，剩余的玩家获得游戏的胜利
- 当新闻卡或者运气卡抽完时，将所有卡牌重新再次组成牌堆
- 只有在向前经过，或刚好到达起点时，才可以领取奖励，后退经过起点不能领取奖励（新闻卡效果）
- 运气卡对玩家移动效果没有方向，故不能领取起点的奖励
- 当新闻卡发放奖励时，有多个合适的人选，则都能拿到奖励
- 已经被购买的电视台和建筑公司没有任何作用
- 建筑公司只对自己的空地有效
- 当你因投掷骰子向前移动而到达监狱时，也不能收取租金或奖励，但你无须停留一回合
- 玩家应随时将手中的小面额金钱与银行交换成大面额金钱
- 你可以保留被抵押的地产，在缴纳税金时，仍会将其计入地产数量
- 投资项目不能被抵押
- 玩家可以通过自行协议，来处理任何说明书内未涉及的情况

## 运气卡

- 遗失钱包，你失去300元，位于你后方的第一位玩家获得300元
- 黑历史被查，立即移动到监狱，并停留一回合
- 社会主义春风吹过，你可以立即免费升级一块抵偿
- 双十一期间，疯狂消费，支付300元
- 前往九寨沟旅游，支付500元
- 潜入银行 系统内部，从每位玩家手中收取300元
- 在香港乘坐豪华游轮，你可以选择支付1000元，并立即移动
到起点领取奖励
- 立即移动到你的左边手玩家的位置，并按该结果结算
- 立即移动到你右手边玩家的位置，并按该结果结算
- 发票刮中奖，获得400元
- 额外获得遗产，获得600元
- 购买最新款私人坐骑，支付400元，并立即额外进行一回合的行动

## 新闻卡

- 投资项目分红，距离证券中心最近的玩家获得500元
- 社会发放福利，每位玩家获得1000元
- 经营不善，拥有核能发电站的玩家失去300元
- 经营不善，拥有污水处理厂的玩家失去300元
- 经营不善，每位拥有运输业的玩家失去300元（大陆运输，大洋运输，空中运输）
- 政府公开补助土地少者500元
- 无名慈善家资助，每位玩家可以立即免费赎回一块抵偿
- 全体玩家参加狂欢节，在你下次行动结束前，所有玩家移动移动时都变为后退
- 经营不善，拥有建筑公司将自己的一个地产下降一级
- 百年一遇特大暴雨，所有玩家原地停留一回合
- 发生灵异事件，在你下次行动结束前，所有玩家都无须支付任何费用
- 所有玩家缴纳个人所得税，每块地产300元