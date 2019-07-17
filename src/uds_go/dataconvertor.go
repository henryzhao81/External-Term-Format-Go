package main

import (
	"ds_user_service"
	"github.com/juju/errors"
)

func TermToAdvertiserData (term Term) (r *ds_user_service.AdvertisersData, err error) {
	r = &ds_user_service.AdvertisersData{}
	if v, ok := term.(*Tuple); ok {
		rtb := v.get(0)
		if l, ok := rtb.(*ErlangList); ok {
			advertisers_rtb := make([]*ds_user_service.AdvRtbData, l.size())
			for i := 0; i < l.size(); i++ {
				each := l.get(i)
				if t, ok := each.(*Tuple); ok {
					var advertiser_id ds_user_service.AdvertiserId
					var rtb_data ds_user_service.RtbData
					if sub_t, ok := t.get(0).(*Tuple); ok {
						if sub_t.size() != 2 {
							err = errors.New("Bad format")
						} else {
							if t_adv_id, ok := sub_t.get(0).([]byte); ok {
								advertiser_id = ds_user_service.AdvertiserId(string(t_adv_id))
							} else {
								err = errors.New("Bad format")
							}
							if t_rtb_data, ok := sub_t.get(1).([]byte); ok {
								rtb_data = ds_user_service.RtbData(string(t_rtb_data))
							} else {
								err = errors.New("Bad format")
							}
						}
					} else {
						err = errors.New("Bad format")
					}
					var timestamp ds_user_service.Timestamp
					if t_timestamp, ok := t.get(1).(int); ok {
						timestamp = ds_user_service.Timestamp(int64(t_timestamp))
					}
					advertisers_rtb[i] = &ds_user_service.AdvRtbData{AdvertiserID:advertiser_id, RtbData:rtb_data, Timestamp:timestamp}
				} else {
					err = errors.New("Bad format")
				}
			}
			r.Rtb = advertisers_rtb
			r.PubSegments = make(map[string]ds_user_service.DmpSegmentMap)
		} else {
			err = errors.New("Bad format")
		}
	} else {
		err = errors.New("Bad format")
	}
	return
}

