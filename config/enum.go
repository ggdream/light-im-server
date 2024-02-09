package config

// 数据表的记录状态
const (
	RecordStatusNormal      int8 = 1 // 正常
	RecordStatusBanned      int8 = 2 // 禁用
	RecordStatusReview      int8 = 3 // 待审核
	RecordStatusInactivated int8 = 4 // 待激活
	RecordStatusDelete      int8 = 5 // 已删除
)

const (
	ConvTypePrivate int8 = 1 // 单聊
	ConvTypeGroup   int8 = 2 // 群聊
)

const (
	GroupMemberRoleOwner    int8 = 1 // 拥有者
	GroupMemberRoleAdmin    int8 = 2 // 管理员
	GroupMemberRoleOrdinary int8 = 3 // 普通成员
	GroupMemberRoleRobot    int8 = 4 // 机器人
)

//// 列表排序方式
//const (
//	OrderByCreateAtDesc  int8 = 1 // 创建时间倒序
//	OrderByLikeCountDesc int8 = 2 // 点赞倒序
//	OrderByHotDesc       int8 = 3 // 热度倒序
//	OrderByTopDesc       int8 = 4 // 顶流倒序
//)

// 数据同步动作类型
const (
	SyncActionCreate int8 = 1
	SyncActionUpdate int8 = 2
	SyncActionDelete int8 = 3
)
