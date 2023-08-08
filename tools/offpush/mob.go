package offpush

import (
	mob_push_sdk "github.com/MobClub/mobpush-websdkv3-go"
)

type MobPush struct {
	client *mob_push_sdk.PushClient
}

func NewMobPush(appKey, appSecret string) *MobPush {
	hc := mob_push_sdk.NewPushClient(appKey, appSecret)

	return &MobPush{
		client: hc,
	}
}

func (m *MobPush) PushByRIDs(title, content string, rids []string) error {
	_, err := m.client.PushByRids("", title, content, rids)
	return err
}
