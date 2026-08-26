package service

import "message-queue/internal/store"

// ValidateOptionalTopicPolicy 是可选主题策略的统一入口。
//
// 设计意图：
//   - 策略未启用（校验器不可用）：放行，返回 nil。
//   - 策略已启用：交由规则集裁决，仅放行已登记的主题；未登记主题返回 "route rejected"。
//
// 早期实现依赖 OptionalTopicPolicyValidatorUsable 单纯判 nil，
// 但 store 层在未启用时返回的是「包裹 nil 指针的接口」(typed-nil interface)，
// 导致 v != nil 恒为真、守卫被穿透，随后 v.Validate(key) 在 nil 接收者上解引用 map 而崩溃。
// 现由 store.OptionalTopicPolicyValidatorUsable 统一承担 typed-nil 防御，
// 此处仅做一次可用性判断即可，两条路径（空策略 / 有效策略）都收敛到正确语义。
func ValidateOptionalTopicPolicy(v store.OptionalTopicPolicyValidator, key string) (err error) {
	if !store.OptionalTopicPolicyValidatorUsable(v) {
		return nil
	}
	return v.Validate(key)
}
