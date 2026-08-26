// Package config 负责从环境变量加载服务配置。
package config

// SignaturePolicy 描述签名算法白名单策略。
type SignaturePolicy struct{ Allowed map[string]bool }

// LoadSignaturePolicy 根据算法名称白名单构建策略。
// 当 names 为空（白名单未配置）时返回 nil。
//
// 注意：由于返回值 *SignaturePolicy 在赋值给接口时，即便是 nil 指针也会
// 被装箱为非 nil 接口，调用方不能仅靠 v.policy != nil 判定“未配置”，
// 必须在调用 Allow 前显式检查 nil 指针，否则会解引用空指针 panic。
func LoadSignaturePolicy(names []string) *SignaturePolicy {
	if len(names) == 0 {
		return nil
	}
	p := &SignaturePolicy{Allowed: make(map[string]bool)}
	for _, name := range names {
		p.Allowed[name] = true
	}
	return p
}

// Allow 判断算法是否在白名单内。nil 接收者视为不允许任何算法。
func (p *SignaturePolicy) Allow(name string) bool {
	if p == nil {
		return false
	}
	return p.Allowed[name]
}
