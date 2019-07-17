package main

import (
	"ds_user_service"
	"fmt"
	"bytes"
)

type DsuThriftHandler struct {
}

func NewDsuThriftHandler() *DsuThriftHandler {
	return &DsuThriftHandler{}
}

func (p *DsuThriftHandler) GetPublisherSegments(storage_id ds_user_service.StorageId, context ds_user_service.Context) (r ds_user_service.PublisherSegments, err error)  {
	fmt.Println("Request GetPublisherSegments")
	return nil, nil
}

func (p *DsuThriftHandler) GetConversions(storage_id ds_user_service.StorageId, conversion_id ds_user_service.ConversionId, context ds_user_service.Context) (r ds_user_service.Conversions, err error) {
	fmt.Println("Request GetConversions")
	return nil, nil
}

func (p *DsuThriftHandler) GetAdvertisersData(market_id ds_user_service.MarketId, context ds_user_service.Context) (r *ds_user_service.AdvertisersData, err error) {
	fmt.Println("Request GetAdvertisersData", market_id)
	riakObj := &RiakObject{}
	payload := &Payload{dbRequest:riakObj, table:"a", key:string(market_id)}
	rsp := dRiak.Execute(&Job{jobName:fmt.Sprintf("job-[%s]", "GetAdvertisersData"), requestType:"get", Payload:payload,}, 2000)
	switch rsp.(type) {
	case Term:
		var buffer *bytes.Buffer
		buffer = &bytes.Buffer{}
		ToString(Term(rsp), buffer)
		fmt.Println("riak result:", string(buffer.Bytes()))
		// Filling the object
                r, err = TermToAdvertiserData(Term(rsp))
		fmt.Println("rtb len", len(r.Rtb))
		return
	}
	return nil, nil
}

func (p *DsuThriftHandler) GetCookieState(id ds_user_service.CookieStateId, platform_hash ds_user_service.PlatformHash, context ds_user_service.Context) (r ds_user_service.CookieState, err error) {
	fmt.Println("Request GetCookieState")
	return "", nil
}

func (p *DsuThriftHandler) SetDeviceData(id ds_user_service.DeviceId, data ds_user_service.DeviceData, context ds_user_service.Context) (err error) {
	fmt.Println("Request SetDeviceData")
	return nil
}

func (p *DsuThriftHandler) DelDeviceData(id ds_user_service.DeviceId, context ds_user_service.Context) (err error) {
	fmt.Println("Request DelDeviceData")
	return nil
}

func (p *DsuThriftHandler) DelDeviceDataByKey(id ds_user_service.DeviceId, keys ds_user_service.DeviceDataKeys, context ds_user_service.Context) (err error) {
	fmt.Println("Request DelDeviceDataByKey")
	return nil
}

func (p *DsuThriftHandler) GetDeviceData(id ds_user_service.DeviceId, context ds_user_service.Context) (r ds_user_service.DeviceData, err error) {
	fmt.Println("Request GetDeviceData")
	return nil, nil
}

func (p *DsuThriftHandler) IsOptedOut(id ds_user_service.DeviceId, context ds_user_service.Context) (r bool, err error) {
	fmt.Println("Request IsOptedOut")
	return false, nil
}

func (p *DsuThriftHandler) GetCrossDeviceData(storage_id ds_user_service.StorageId, context ds_user_service.Context) (r *ds_user_service.CrossDeviceData, err error) {
	fmt.Println("Request GetCrossDeviceData")
	return nil, nil
}