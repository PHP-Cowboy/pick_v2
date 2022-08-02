package helper

func GetDeliveryMethod(method int) (ret string) {

	switch method {
	case 1:
		ret = "公司配送"
		break
	case 2:
		ret = "用户自提"
		break
	case 3:
		ret = "三方物流"
		break
	case 4:
		ret = "快递配送"
		break
	case 5:
		ret = "首批物料|设备单"
		break
	}
	return ret
}
